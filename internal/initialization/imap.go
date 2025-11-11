package initialization

import (
	"fmt"
	"log"

	"github.com/InstaySystem/is-be/internal/config"
	idle "github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
)

func InitIMAP(cfg *config.Config) (*client.Client, error) {
	imapAddr := fmt.Sprintf("%s:%d", cfg.IMAP.Host, cfg.IMAP.Port)
	c, err := client.DialTLS(imapAddr, nil)
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	if err = c.Login(cfg.IMAP.User, cfg.IMAP.Password); err != nil {
		return nil, err
	}

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, err
	}
	log.Printf("Mailbox status: %d messages\n", mbox.Messages)

	updates := make(chan client.Update)
	c.Updates = updates

	idleClient := idle.NewClient(c)
	stop := make(chan struct{})

	go func() {
		for update := range updates {
			// Đây là nơi nhận realtime event
			switch u := update.(type) {
			case *client.MailboxUpdate:
				fmt.Println("Mailbox updated:", u.Mailbox)
			case *client.MessageUpdate:
				fmt.Println("New message arrived!")
				// Tại đây bạn fetch nội dung mail
			}
		}
	}()

	log.Println("Listening for new messages...")

	for {
		// Kích hoạt chế độ IDLE
		if err := idleClient.Idle(stop); err != nil {
			log.Println("Idle error:", err)
		}
	}
}
