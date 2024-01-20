package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.elastic.co/apm"
	"log"
	"notification-api/domain"
	"notification-api/service/notification"
	"time"
)

type ConsumerNotificationHandler struct {
	service notification.ConsumerNotification
	ctx     context.Context
	tx      *apm.Transaction
}

func NewConsumerNotificationHandler(svc notification.ConsumerNotification, ctx context.Context, tx *apm.Transaction) *ConsumerNotificationHandler {
	return &ConsumerNotificationHandler{service: svc, ctx: ctx, tx: tx}
}

func (s ConsumerNotificationHandler) NotificationReportHandle(msg *message.Message) (err error) {
	fmt.Printf(
		"\n> Notification Consume successfully: \n> at: %s\n> \n\n", time.Now(),
	)

	log.Println("received message", msg.UUID)
	var data domain.NotificationData

	if err := json.Unmarshal(msg.Payload, &data); err != nil {
		apm.CaptureError(s.ctx, err).Send()
		fmt.Println("Consumer Handle | error parse to video with error: " + err.Error())
	}
	if err := s.service.ProcessMessage(s.ctx, &data); err != nil {
		log.Println("Consumer Handle | Failed Consume/Insert Database", err)
	}
	log.Println("Consumer Handle | Success insert Database")
	s.tx.End()
	msg.Ack()
	fmt.Printf(
		"\n> Received message NotificationHandler: %s\n> at: %s\n> payload: %s \n> metadata: %v\n\n",
		msg.UUID,
		time.Now(),
		string(msg.Payload),
		msg.Metadata,
	)
	return nil
}

func (s ConsumerNotificationHandler) PrintMessages(msg *message.Message) error {
	fmt.Printf(
		"\n> Received message: %s\n> at: %s\n> payload: %s \n> metadata: %v\n\n",
		msg.UUID,
		time.Now(),
		string(msg.Payload),
		msg.Metadata,
	)
	return nil
}
