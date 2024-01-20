package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"go.elastic.co/apm/v2"
	"log"
	"notification-api/app/infrastructure/kafka"
	"notification-api/config"
	"notification-api/domain"
	"notification-api/domain/entities"
	sg "notification-api/repositories/email"
	"notification-api/repositories/notification"
	ctap "notification-api/repositories/push_notification"
	sms2 "notification-api/repositories/sms"
	"time"
)

type ConsumerNotification interface {
	ProcessMessage(ctx context.Context, data *domain.NotificationData) error
}

type consumerNotificationService struct {
	kp         *kafka.Subscriber
	repository notification.Repository
	email      sg.EmailRepository
	pushNotif  ctap.PushNotifRepository
	sms        sms2.SmsRepository
}

func NewConsumerNotificationService(kp *kafka.Subscriber, repository notification.Repository, email sg.EmailRepository, pushNotif ctap.PushNotifRepository, sms sms2.SmsRepository) ConsumerNotification {
	return &consumerNotificationService{
		kp:         kp,
		repository: repository,
		email:      email,
		pushNotif:  pushNotif,
		sms:        sms,
	}
}

func (c consumerNotificationService) ProcessMessage(ctx context.Context, data *domain.NotificationData) error {
	span, ctx := apm.StartSpan(ctx, "service/notification/service_consumer_notification", "ProcessMessage")
	defer span.End()
	var (
		jobData entities.Job
	)
	switch data.Type {
	case "event":
		d := domain.MetadataPushNotif{
			Identity: data.Identity,
			Ts:       data.Timestamp,
			Type:     data.Type,
			EvtName:  data.EvtName,
			EvtData:  data.EvtData,
		}

		payload := domain.PushNotificationRequest{Data: []domain.MetadataPushNotif{d}}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			log.Println(err)
		}
		fmt.Println("===================>  Payload: ", string(jsonData))
		resultChan := make(chan *domain.ResponseNotif)

		// Configuring the Hystrix circuit
		hystrix.ConfigureCommand("post", hystrix.CommandConfig{
			Timeout:               config.Config.GetInt("HYSTRIX_TIMEOUT"),
			MaxConcurrentRequests: config.Config.GetInt("HYSTRIX_MAX_CONCURRENT_REQUESTS"),
			ErrorPercentThreshold: config.Config.GetInt("HYSTRIX_ERROR_PERCENT_THRESHOLD"),
		})

		// Using hystrix.Do to execute the HTTP request with circuit breaker protection
		errChan := hystrix.Go("post", func() error {
			resp, err := c.pushNotif.PostPushNotif(ctx, payload)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				return err
			}
			resultChan <- &resp
			return nil
		}, nil)
		select {
		case errors := <-errChan:
			log.Println("Failed Post", errors.Error())
			jobData.Status = "failed"
			jobData.ErrorMessage = errors.Error()
			jobData.FinishedAt = time.Now().Add(7 * time.Hour)
			err := c.repository.UpdateJob(ctx, data.RequestID, jobData)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				log.Println(err)
			}
			return errors
		case result := <-resultChan:
			if result.Status != "success" {
				log.Println("Failed Post:", result)
				jobData.Status = "failed"
				jobData.ErrorMessage = result.ErrorMessage
				jobData.FinishedAt = time.Now().Add(7 * time.Hour)
				err = c.repository.UpdateJob(ctx, data.RequestID, jobData)
				if err != nil {
					apm.CaptureError(ctx, err).Send()
					log.Println(err)
					return err
				}
				return nil
			}
			log.Println("Success:", result)
			jobData.Status = "success"
			jobData.FinishedAt = time.Now().Add(7 * time.Hour)
			err = c.repository.UpdateJob(ctx, data.RequestID, jobData)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				log.Println(err)
				return err
			}
			return nil
		}

	case "email":
		d := domain.MetadataEmail{
			Identity:      data.Identity,
			Type:          data.Type,
			TemplateName:  data.EvtName,
			TemplateParam: data.EvtData,
		}

		jsonData, err := json.Marshal(d)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			log.Println(err)
		}
		fmt.Println("===================>  Payload: ", string(jsonData))
		resultChan := make(chan *domain.ResponseNotif)

		// Configuring the Hystrix circuit
		hystrix.ConfigureCommand("post", hystrix.CommandConfig{
			Timeout:               config.Config.GetInt("HYSTRIX_TIMEOUT"),
			MaxConcurrentRequests: config.Config.GetInt("HYSTRIX_MAX_CONCURRENT_REQUESTS"),
			ErrorPercentThreshold: config.Config.GetInt("HYSTRIX_ERROR_PERCENT_THRESHOLD"),
		})

		// Using hystrix.Do to execute the HTTP request with circuit breaker protection
		errChan := hystrix.Go("post", func() error {
			resp, err := c.email.PostEmail(ctx, d)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				return err
			}
			resultChan <- &resp
			return nil
		}, nil)
		select {
		case errors := <-errChan:
			log.Println("Failed Post", errors.Error())
			jobData.Status = "failed"
			jobData.ErrorMessage = errors.Error()
			jobData.FinishedAt = time.Now().Add(7 * time.Hour)
			err := c.repository.UpdateJob(ctx, data.RequestID, jobData)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				log.Println(err)
			}
			return errors
		case result := <-resultChan:
			if result.Status != "success" {
				log.Println("Failed Post:", result)
				jobData.Status = "failed"
				jobData.ErrorMessage = result.ErrorMessage
				jobData.FinishedAt = time.Now().Add(7 * time.Hour)
				err = c.repository.UpdateJob(ctx, data.RequestID, jobData)
				if err != nil {
					apm.CaptureError(ctx, err).Send()
					log.Println(err)
					return err
				}
				return nil
			}
			log.Println("Success:", result)
			jobData.Status = "success"
			jobData.FinishedAt = time.Now().Add(7 * time.Hour)
			err = c.repository.UpdateJob(ctx, data.RequestID, jobData)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				log.Println(err)
				return err
			}
			return nil
		}
	case "sms":
		resultChan := make(chan *domain.ResponseNotif)
		otp, _ := data.EvtData["otp"].(string)

		// Configuring the Hystrix circuit
		hystrix.ConfigureCommand("post", hystrix.CommandConfig{
			Timeout:               config.Config.GetInt("HYSTRIX_TIMEOUT"),
			MaxConcurrentRequests: config.Config.GetInt("HYSTRIX_MAX_CONCURRENT_REQUESTS"),
			ErrorPercentThreshold: config.Config.GetInt("HYSTRIX_ERROR_PERCENT_THRESHOLD"),
		})

		// Using hystrix.Do to execute the HTTP request with circuit breaker protection
		errChan := hystrix.Go("post", func() error {
			resp, err := c.sms.SendSMS(ctx, data.Identity, otp, data.EvtName)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				return err
			}
			resultChan <- &resp
			return nil
		}, nil)
		select {
		case errors := <-errChan:
			log.Println("Failed Post", errors.Error())
			jobData.Status = "failed"
			jobData.ErrorMessage = errors.Error()
			jobData.FinishedAt = time.Now().Add(7 * time.Hour)
			err := c.repository.UpdateJob(ctx, data.RequestID, jobData)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				log.Println(err)
			}
			return errors
		case result := <-resultChan:
			if result.Status != "success" {
				log.Println("Failed Post:", result)
				jobData.Status = "failed"
				jobData.ErrorMessage = result.ErrorMessage
				jobData.FinishedAt = time.Now().Add(7 * time.Hour)
				err := c.repository.UpdateJob(ctx, data.RequestID, jobData)
				if err != nil {
					apm.CaptureError(ctx, err).Send()
					log.Println(err)
					return err
				}
				return nil
			}
			log.Println("Success:", result)
			jobData.Status = "success"
			jobData.FinishedAt = time.Now().Add(7 * time.Hour)
			err := c.repository.UpdateJob(ctx, data.RequestID, jobData)
			if err != nil {
				apm.CaptureError(ctx, err).Send()
				log.Println(err)
				return err
			}
			return nil
		}
	}
	return nil
}
