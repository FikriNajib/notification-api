package consumer

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"go.elastic.co/apm"
	"notification-api/app/consumer/handler"
	kafka2 "notification-api/app/infrastructure/kafka"
	"notification-api/app/infrastructure/mysql"
	"notification-api/config"
	s "notification-api/repositories/email"
	notification2 "notification-api/repositories/notification"
	cl "notification-api/repositories/push_notification"
	sms2 "notification-api/repositories/sms"
	"notification-api/service/notification"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

type ConsumerImpl struct {
	router *message.Router
}
type Consumer interface {
	InitRouter() *ConsumerImpl
	Listen() *ConsumerImpl
	Serve()
}

func NewConsumer() Consumer {
	return &ConsumerImpl{}
}

func (c *ConsumerImpl) InitRouter() *ConsumerImpl {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)

	// Router level middleware are executed for every message sent to the router
	router.AddMiddleware(
		middleware.CorrelationID,

		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,

		middleware.Recoverer,
	)
	c.router = router
	return c
}

func (c *ConsumerImpl) Listen() *ConsumerImpl {
	tx := apm.DefaultTracer.StartTransaction("Consumer Notification", "CreateJob")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	span, ctx := apm.StartSpan(ctx, "Listen", "CreateJob")
	defer span.End()

	subscriber := kafka2.NewSubscriber()
	database := mysql.NewDatabase().ConnectNotificationDB().Execute()

	repo := notification2.NewNotificationRepository(database)
	sendgrid := s.NewSendgridRepository()
	clevertap := cl.NewClevertapRepository()
	sms := sms2.NewSMSRepository()
	svcConsumerNotification := notification.NewConsumerNotificationService(subscriber, repo, sendgrid, clevertap, sms)
	notificationHandler := handler.NewConsumerNotificationHandler(svcConsumerNotification, ctx, tx)

	c.router.AddNoPublisherHandler(
		"NotificationHandle",
		config.Config.GetString("TOPIC"),
		subscriber,
		notificationHandler.NotificationReportHandle,
	)

	if err := c.router.Run(ctx); err != nil {
		panic(err)
	}
	return c
}

func (c *ConsumerImpl) Serve() {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	for {
	}
}
