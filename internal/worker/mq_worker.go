package worker

import (
	"context"
	"encoding/json"
	"fmt"
	// "log"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/provider/smtp"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/emersion/go-imap"
	// "github.com/emersion/go-imap/client"
	"go.uber.org/zap"
)

type MQWorker struct {
	cfg *config.Config
	mq     mq.MessageQueueProvider
	smtp   smtp.SMTPProvider
	// imap   *client.Client
	s3 *awsS3.Client
	logger *zap.Logger
}

func NewMQWorker(
	cfg *config.Config,
	mq mq.MessageQueueProvider,
	smtp smtp.SMTPProvider,
	s3 *awsS3.Client,
	// imap *client.Client,
	logger *zap.Logger,
) *MQWorker {
	return &MQWorker{
		cfg,
		mq,
		smtp,
		s3,
		// imap,
		logger,
	}
}

func (w *MQWorker) StartSendAuthEmail() {
	if err := w.mq.ConsumeMessage(common.QueueNameAuthEmail, common.ExchangeEmail, common.RoutingKeyAuthEmail, func(body []byte) error {
		var emailMsg types.AuthEmailMessage
		if err := json.Unmarshal(body, &emailMsg); err != nil {
			return err
		}

		if err := w.smtp.AuthEmail(emailMsg.To, emailMsg.Subject, emailMsg.Otp); err != nil {
			return err
		}

		w.logger.Info(fmt.Sprintf("Email sent successfully to: %s", emailMsg.To))
		return nil
	}); err != nil {
		w.logger.Error("start consumer send auth email failed", zap.Error(err))
	}
}

func (w *MQWorker) StartDeleteFile() {
	if err := w.mq.ConsumeMessage(common.QueueNameDeleteFile, common.ExchangeFile, common.RoutingKeyDeleteFile, func(body []byte) error {
		key := string(body)

		ctx := context.Background()

		if _, err := w.s3.HeadObject(ctx, &awsS3.HeadObjectInput{
			Bucket: aws.String(w.cfg.S3.Bucket),
			Key: aws.String(key),
		}); err != nil {
			w.logger.Error("file check failed", zap.Error(err))
		}

		if _, err := w.s3.DeleteObject(ctx, &awsS3.DeleteObjectInput{
			Bucket: aws.String(w.cfg.S3.Bucket),
			Key: aws.String(key),
		}); err != nil {
			w.logger.Error("file delete failed", zap.Error(err))
		}

		return nil
	}); err != nil {
		w.logger.Error("start consumer delete file failed", zap.Error(err))
	}
}

// func (w *EmailWorker) StartListenEmail() {
// 	mbox, err := w.imap.Select("INBOX", false)
// 	if err != nil {
// 		w.logger.Error("select inbox failed", zap.Error(err))
// 	}
// 	lastSeenSeqNum := mbox.Messages

// 	updates := make(chan client.Update, 1)
// 	w.imap.Updates = updates

// 	stop := make(chan struct{})
// 	done := make(chan error, 1)

// 	go func() {
// 		log.Println("Bắt đầu IDLE...")
// 		done <- w.imap.Idle(stop, nil)
// 	}()

// 	for {
// 		select {
// 		case update := <-updates:
// 			// Chúng ta nhận được một cập nhật
// 			if mboxUpdate, ok := update.(*client.MailboxUpdate); ok {
// 				log.Println("Có cập nhật Mailbox:", mboxUpdate.Mailbox.Messages)
// 				if mboxUpdate.Mailbox.Messages > lastSeenSeqNum {
// 					log.Println("Phát hiện email mới!")

// 					// --- 6a. Dừng IDLE ---
// 					log.Println("Đang dừng IDLE để Fetch...")
// 					// Đóng channel 'stop' để ra hiệu c.Idle() kết thúc
// 					close(stop)

// 					// Đợi goroutine 'Idle' báo là đã dừng
// 					if err := <-done; err != nil {
// 						log.Println("Lỗi khi dừng IDLE:", err)
// 						// Cần xử lý kết nối lại ở đây
// 						return
// 					}
// 					log.Println("IDLE đã dừng.")

// 					// --- 6b. Fetch thư mới ---
// 					seqSet := new(imap.SeqSet)
// 					seqSet.AddRange(lastSeenSeqNum+1, mboxUpdate.Mailbox.Messages)

// 					messages := make(chan *imap.Message, 10)
// 					go func() {
// 						if err := w.imap.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
// 							log.Println("Lỗi khi Fetch:", err)
// 						}
// 					}()

// 					for msg := range messages {
// 						log.Printf("* Email Mới: %d - Chủ đề: %s", msg.SeqNum, msg.Envelope.Subject)
// 					}

// 					lastSeenSeqNum = mboxUpdate.Mailbox.Messages

// 					// --- 6c. Khởi động lại IDLE ---
// 					// Tạo lại các channel
// 					stop = make(chan struct{})
// 					done = make(chan error, 1)
// 					go func() {
// 						log.Println("Bắt đầu IDLE trở lại...")
// 						done <- w.imap.Idle(stop, nil)
// 					}()
// 				}
// 			}

// 		case err := <-done:
// 			// Goroutine 'Idle' đã kết thúc (do lỗi hoặc bị dừng)
// 			log.Println("IDLE đã thoát vòng lặp.")
// 			if err != nil {
// 				log.Println("IDLE bị dừng do lỗi:", err)
// 				log.Println("Cần thực hiện kết nối lại (reconnect)...")
// 				// Đây là nơi bạn đặt logic kết nối lại
// 				return // Thoát chương trình cho đơn giản
// 			}
// 			// Nếu err == nil, có nghĩa là chúng ta đã chủ động 'close(stop)'
// 			// Vòng lặp sẽ tiếp tục và chờ 'update' (đã fetch xong)
// 		}
// 	}
// }
