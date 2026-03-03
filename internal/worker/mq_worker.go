package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/InstaySystem/is_v1-be/internal/common"
	"github.com/InstaySystem/is_v1-be/internal/config"
	"github.com/InstaySystem/is_v1-be/internal/hub"
	"github.com/InstaySystem/is_v1-be/internal/provider/mq"
	"github.com/InstaySystem/is_v1-be/internal/provider/smtp"
	"github.com/InstaySystem/is_v1-be/internal/types"
	"go.uber.org/zap"
)

type MQWorker struct {
	cfg    *config.Config
	mq     mq.MessageQueueProvider
	smtp   smtp.SMTPProvider
	gcs    *storage.Client
	logger *zap.Logger
	sseHub *hub.SSEHub
}

func NewMQWorker(
	cfg *config.Config,
	mq mq.MessageQueueProvider,
	smtp smtp.SMTPProvider,
	gcs *storage.Client,
	logger *zap.Logger,
	sseHub *hub.SSEHub,
) *MQWorker {
	return &MQWorker{
		cfg,
		mq,
		smtp,
		gcs,
		logger,
		sseHub,
	}
}

func (w *MQWorker) Start() {
	go w.startSendAuthEmail()
	go w.startDeleteFile()
	go w.startSendServiceNotification()
	go w.startSendRequestNotification()
}

func (w *MQWorker) startSendAuthEmail() {
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

func (w *MQWorker) startDeleteFile() {
	if err := w.mq.ConsumeMessage(common.QueueNameDeleteFile, common.ExchangeFile, common.RoutingKeyDeleteFile, func(body []byte) error {
		key := string(body)

		ctx := context.Background()

		obj := w.gcs.Bucket(w.cfg.GCS.Bucket).Object(key)

		if _, err := obj.Attrs(ctx); err != nil {
			if err == storage.ErrObjectNotExist {
				w.logger.Warn("file not found", zap.String("key", key))
				return nil
			}
			w.logger.Error("file check failed", zap.Error(err))
			return err
		}

		if err := obj.Delete(ctx); err != nil {
			w.logger.Error("file delete failed", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		w.logger.Error("start consumer delete file failed", zap.Error(err))
	}
}

func (w *MQWorker) startSendServiceNotification() {
	if err := w.mq.ConsumeMessage(common.QueueNameServiceNotification, common.ExchangeNotification, common.RoutingKeyServiceNotification, func(body []byte) error {
		var serviceNotificationMsg types.NotificationMessage
		if err := json.Unmarshal(body, &serviceNotificationMsg); err != nil {
			return err
		}

		data := map[string]any{
			"content":      serviceNotificationMsg.Content,
			"content_id":   serviceNotificationMsg.ContentID,
			"content_type": serviceNotificationMsg.Type,
			"receiver":     serviceNotificationMsg.Receiver,
		}

		event := types.SSEEventData{
			Event:        "order_service",
			Type:         serviceNotificationMsg.Receiver,
			DepartmentID: serviceNotificationMsg.DepartmentID,
			Data:         data,
		}

		for _, clientID := range serviceNotificationMsg.ReceiverIDs {
			w.sseHub.SendToClient(clientID, event)
		}

		w.logger.Info("Service notification sent successfully")
		return nil
	}); err != nil {
		w.logger.Error("start consumer send service notification failed", zap.Error(err))
	}
}

func (w *MQWorker) startSendRequestNotification() {
	if err := w.mq.ConsumeMessage(common.QueueNameRequestNotification, common.ExchangeNotification, common.RoutingKeyRequestNotification, func(body []byte) error {
		var requestNotificationMsg types.NotificationMessage
		if err := json.Unmarshal(body, &requestNotificationMsg); err != nil {
			return err
		}

		data := map[string]any{
			"content":      requestNotificationMsg.Content,
			"content_id":   requestNotificationMsg.ContentID,
			"content_type": requestNotificationMsg.Type,
			"receiver":     requestNotificationMsg.Receiver,
		}

		event := types.SSEEventData{
			Event:        "request",
			Type:         requestNotificationMsg.Receiver,
			DepartmentID: requestNotificationMsg.DepartmentID,
			Data:         data,
		}

		for _, clientID := range requestNotificationMsg.ReceiverIDs {
			w.sseHub.SendToClient(clientID, event)
		}

		w.logger.Info("Request notification sent successfully")
		return nil
	}); err != nil {
		w.logger.Error("start consumer send request notification failed", zap.Error(err))
	}
}
