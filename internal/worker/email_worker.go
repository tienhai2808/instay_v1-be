package worker

import (
	"encoding/json"
	"fmt"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/provider/smtp"
	"github.com/InstaySystem/is-be/internal/types"
	"go.uber.org/zap"
)

type EmailWorker struct {
	mq     mq.MessageQueueProvider
	smtp   smtp.SMTPProvider
	logger *zap.Logger
}

func NewEmailWorker(
	mq mq.MessageQueueProvider,
	smtp smtp.SMTPProvider,
	logger *zap.Logger,
) *EmailWorker {
	return &EmailWorker{
		mq,
		smtp,
		logger,
	}
}

func (w *EmailWorker) StartSendAuthEmail() {
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
