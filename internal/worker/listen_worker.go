package worker

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"github.com/PuerkitoBio/goquery"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"go.uber.org/zap"
)

const (
	maxReconnectAttempts = 10
	baseReconnectDelay   = 2 * time.Second
	maxReconnectDelay    = 5 * time.Minute
	idleTimeout          = 25 * time.Minute
	healthCheckInterval  = 5 * time.Minute
	targetSubject        = "CONGRATULATIONS! You've received a new booking"
)

type ListenWorker struct {
	cfg         *config.Config
	bookingRepo repository.BookingRepository
	sfGen       snowflake.Generator
	logger      *zap.Logger
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewListenWorker(
	cfg *config.Config,
	bookingRepo repository.BookingRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) *ListenWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &ListenWorker{
		cfg,
		bookingRepo,
		sfGen,
		logger,
		ctx,
		cancel,
	}
}

func (w *ListenWorker) Start() {
	w.logger.Info("Starting email listener")
	go w.runEmailListener()
}

func (w *ListenWorker) Stop() {
	w.logger.Info("Stopping listen worker")
	w.cancel()
}

func (w *ListenWorker) runEmailListener() {
	attempt := 0
	for {
		select {
		case <-w.ctx.Done():
			w.logger.Info("Email listener context cancelled, exiting")
			return
		default:
		}

		if err := w.listenEmails(); err != nil {
			attempt++
			delay := calculateBackoff(attempt)
			w.logger.Error("Email listener error, reconnecting",
				zap.Error(err),
				zap.Int("attempt", attempt),
				zap.Duration("retry_in", delay))

			select {
			case <-time.After(delay):
				continue
			case <-w.ctx.Done():
				return
			}
		}

		attempt = 0
	}
}

func (w *ListenWorker) listenEmails() error {
	imapAddr := fmt.Sprintf("%s:%d", w.cfg.IMAP.Host, w.cfg.IMAP.Port)
	c, err := client.DialTLS(imapAddr, nil)
	if err != nil {
		return err
	}
	defer c.Logout()

	if err := c.Login(w.cfg.IMAP.User, w.cfg.IMAP.Password); err != nil {
		return err
	}

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return err
	}

	lastSeenSeqNum := mbox.Messages
	w.logger.Info("INBOX selected",
		zap.Uint32("total_messages", mbox.Messages),
		zap.Uint32("last_seen", lastSeenSeqNum))

	updates := make(chan client.Update, 10)
	c.Updates = updates

	stop := make(chan struct{})
	done := make(chan error, 1)
	healthTicker := time.NewTicker(healthCheckInterval)
	defer healthTicker.Stop()

	go func() {
		done <- c.Idle(stop, &client.IdleOptions{
			LogoutTimeout: idleTimeout,
			PollInterval:  0,
		})
	}()

	for {
		select {
		case <-w.ctx.Done():
			close(stop)
			<-done
			return nil

		case <-healthTicker.C:

		case update := <-updates:
			if mboxUpdate, ok := update.(*client.MailboxUpdate); ok {
				if mboxUpdate.Mailbox.Messages > lastSeenSeqNum {
					close(stop)
					if err := <-done; err != nil {
						return err
					}

					if err := w.fetchNewEmails(c, lastSeenSeqNum+1, mboxUpdate.Mailbox.Messages); err != nil {
						w.logger.Error("fetch new emails failed", zap.Error(err))
					}

					lastSeenSeqNum = mboxUpdate.Mailbox.Messages

					stop = make(chan struct{})
					done = make(chan error, 1)
					go func() {
						done <- c.Idle(stop, &client.IdleOptions{
							LogoutTimeout: idleTimeout,
							PollInterval:  0,
						})
					}()
				}
			}

		case err := <-done:
			if err != nil {
				return err
			}
		}
	}
}

func (w *ListenWorker) fetchNewEmails(c *client.Client, from, to uint32) error {
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchUid, section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	for msg := range messages {
		if msg == nil {
			continue
		}

		if msg.Envelope != nil && msg.Envelope.Subject == targetSubject {
			w.processEmail(msg, section)
		}
	}

	if err := <-done; err != nil {
		return err
	}

	return nil
}

func (w *ListenWorker) processEmail(msg *imap.Message, section *imap.BodySectionName) {
	r := msg.GetBody(section)
	if r == nil {
		return
	}
	mr, err := mail.CreateReader(r)
	if err != nil {
		w.logger.Error("create mail reader failed", zap.Error(err))
		return
	}

	var htmlBody string
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			w.logger.Error("read part failed", zap.Error(err))
			break
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()
			body, err := io.ReadAll(p.Body)
			if err != nil {
				w.logger.Error("reading body failed", zap.Error(err))
				continue
			}

			if contentType == "text/html" {
				htmlBody = string(body)
			}
		}
	}

	if htmlBody == "" {
		return
	}

	bookingData, err := w.parseBookingFromHTML(htmlBody)
	if err != nil {
		w.logger.Error("parse booking from HTML failed", zap.Error(err))
		return
	}

	bookingData.ID, err = w.sfGen.NextID()
	if err != nil {
		w.logger.Error("generate booking ID failed", zap.Error(err))
		return
	}

	if err := w.bookingRepo.CreateBooking(w.ctx, bookingData); err != nil {
		w.logger.Error("create booking failed", zap.Error(err))
	} else {
		w.logger.Info("booking created successfully",
			zap.String("booking_number", bookingData.BookingNumber),
			zap.Int64("id", bookingData.ID))
	}
}

func (w *ListenWorker) parseBookingFromHTML(htmlContent string) (*model.Booking, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	booking := &model.Booking{}

	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Booking number:") {
			booking.BookingNumber = strings.TrimSpace(s.Find("span").Last().Text())
		}
	})

	findValueByLabel := func(label string) *goquery.Selection {
		var result *goquery.Selection
		doc.Find("div[style*='display: flex']").Each(func(i int, s *goquery.Selection) {
			labelDiv := s.Find("div").First()
			if strings.Contains(strings.TrimSpace(labelDiv.Text()), label) {
				result = s.Find("div").Eq(1)
			}
		})
		return result
	}

	if guestDiv := findValueByLabel("Guest:"); guestDiv != nil {
		guestDiv.Find("div").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text == "" {
				return
			}
			if strings.Contains(text, "@") {
				booking.GuestEmail = text
			} else if isPhoneNumber(text) {
				booking.GuestPhone = text
			} else {
				if booking.GuestFullName == "" {
					booking.GuestFullName = text
				}
			}
		})
	}

	if val := findValueByLabel("Check-in:"); val != nil {
		booking.CheckIn = parseDateString(val.Text())
	}
	if val := findValueByLabel("Check-out:"); val != nil {
		booking.CheckOut = parseDateString(val.Text())
	}
	if val := findValueByLabel("Booked on:"); val != nil {
		booking.BookedOn = parseDateString(val.Text())
	}

	if val := findValueByLabel("Rooms booked:"); val != nil {
		rawRoom := strings.TrimSpace(val.Text())
		parts := strings.SplitN(rawRoom, " ", 2)
		if len(parts) > 0 {
			if num, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
				booking.RoomNumber = uint32(num)
			}
		}
		if len(parts) > 1 {
			booking.RoomType = strings.TrimSpace(parts[1])
		} else {
			booking.RoomType = rawRoom
		}
	}

	if val := findValueByLabel("Booking source:"); val != nil {
		booking.Source = strings.TrimSpace(val.Text())
	}

	if val := findValueByLabel("Total net price:"); val != nil {
		booking.TotalNetPrice = parsePrice(val.Text())
	}
	if val := findValueByLabel("Total sell price:"); val != nil {
		booking.TotalSellPrice = parsePrice(val.Text())
	}

	if val := findValueByLabel("Number of guests:"); val != nil {
		booking.GuestNumber = strings.TrimSpace(val.Text())
	}

	if val := findValueByLabel("Promo name:"); val != nil {
		booking.PromotionName = strings.TrimSpace(val.Text())
	}
	if val := findValueByLabel("Booking conditions:"); val != nil {
		booking.BookingConditions = strings.TrimSpace(val.Text())
	}

	if booking.BookingNumber == "" {
		return nil, fmt.Errorf("could not parse booking number")
	}

	return booking, nil
}

func parseDateString(raw string) time.Time {
	raw = strings.TrimSpace(raw)

	if strings.Contains(raw, "from") {
		layoutCheckIn := "Monday, January 2, 2006 from 15:04"
		if t, err := time.Parse(layoutCheckIn, raw); err == nil {
			return t
		}
	}

	if strings.Contains(raw, "until") {
		layoutCheckOut := "Monday, January 2, 2006 until 15:04"
		if t, err := time.Parse(layoutCheckOut, raw); err == nil {
			return t
		}
	}

	layoutDateOnly := "Monday, January 2, 2006"
	t, err := time.Parse(layoutDateOnly, raw)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parsePrice(raw string) float64 {
	reg, _ := regexp.Compile("[^0-9,.]+")
	processed := reg.ReplaceAllString(raw, "")

	processed = strings.ReplaceAll(processed, ".", "")
	processed = strings.ReplaceAll(processed, ",", ".")

	price, err := strconv.ParseFloat(processed, 64)
	if err != nil {
		return 0
	}
	return price
}

func isPhoneNumber(s string) bool {
	hasDigit := false
	for _, r := range s {
		if r >= '0' && r <= '9' {
			hasDigit = true
			break
		}
	}
	return hasDigit && (strings.Contains(s, "+") || len(s) > 6)
}

func calculateBackoff(attempt int) time.Duration {
	if attempt > maxReconnectAttempts {
		attempt = maxReconnectAttempts
	}
	delay := min(baseReconnectDelay*time.Duration(1<<uint(attempt)), maxReconnectDelay)
	return delay
}
