package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/afex/hystrix-go/hystrix"
	"go.elastic.co/apm/v2"
	"log"
	"notification-api/app/infrastructure/kafka"
	"notification-api/config"
	"notification-api/domain"
	"notification-api/domain/entities"
	"notification-api/repositories/notification"
	"strconv"
	"time"
)

type Notification interface {
	NotificationReport(ctx context.Context, request *domain.NotificationData) (string, error)
	GetStatus(ctx context.Context, reqID string) (string, error)
}

type notificationService struct {
	kp         *kafka.Publisher
	repository notification.Repository
}

func NewNotificationService(kp *kafka.Publisher, r notification.Repository) Notification {
	return &notificationService{
		kp:         kp,
		repository: r,
	}
}
func (s *notificationService) GetStatus(ctx context.Context, reqID string) (string, error) {
	span, ctx := apm.StartSpan(ctx, "service/notification/service_notification.go", "GetStatus")
	defer span.End()
	status, err := s.repository.GetJobStatus(ctx, reqID)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		log.Println(err)
		return "", err
	}

	return status, nil
}

func (s *notificationService) NotificationReport(ctx context.Context, request *domain.NotificationData) (string, error) {
	span, ctx := apm.StartSpan(ctx, "service/notification/service_notification.go", "NotificationReport")
	defer span.End()
	request.Timestamp = time.Now().Unix()
	tsNano := strconv.Itoa(int(time.Now().UnixNano()))
	request.RequestID = request.Caller + "-" + tsNano

	payload, err := json.Marshal(request)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		fmt.Println("Error marshal:", err)
		return "", err
	}
	dataJob := entities.Job{
		RequestID:    request.RequestID,
		Caller:       request.Caller,
		Type:         request.Type,
		Name:         request.EvtName,
		Payload:      string(payload),
		Status:       "pending",
		ErrorMessage: "",
		CreatedAt:    time.Now().Add(7 * time.Hour),
	}
	reqID, err := s.repository.InsertJob(ctx, dataJob)
	if err != nil {
		apm.CaptureError(ctx, err)
		log.Println("Error Insert Job :", err)
		return "", err
	}
	hystrix.ConfigureCommand("publish", hystrix.CommandConfig{
		Timeout:               config.Config.GetInt("HYSTRIX_TIMEOUT"),
		MaxConcurrentRequests: config.Config.GetInt("HYSTRIX_MAX_CONCURRENT_REQUESTS"),
		ErrorPercentThreshold: config.Config.GetInt("HYSTRIX_ERROR_PERCENT_THRESHOLD"),
	})

	hystrix.Go("publish", func() error {
		return s.publish(ctx, config.Config.GetString("TOPIC"), request)
	}, func(err error) error {
		fmt.Println("PublishNotificationDetail =>", err.Error())

		log.Println("Circuit breaker opened for publish command", err.Error())

		return nil
	})

	return reqID, nil
}

func (s *notificationService) publish(ctx context.Context, topic string, request *domain.NotificationData) error {
	span, ctx := apm.StartSpan(ctx, "service/notification/service_notification.go", "publish")
	defer span.End()

	jsonString, err := json.Marshal(request)
	if err != nil {
		fmt.Println("publish =>", err.Error())
		apm.CaptureError(ctx, err).Send()
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), jsonString)
	if err = s.kp.Publish(topic, msg); err != nil {
		fmt.Println("publish =>", err.Error())
		apm.CaptureError(ctx, err).Send()
		return err
	}

	return nil
}
