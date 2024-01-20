package api

import (
	"fmt"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.elastic.co/apm/module/apmfiber/v2"
	"log"
	handlerV1 "notification-api/app/api/handler/behavior/v1"
	"notification-api/app/infrastructure/kafka"
	"notification-api/app/infrastructure/mysql"
	"notification-api/config"
	notification2 "notification-api/repositories/notification"
	"notification-api/service/notification"
	"os"
	"os/signal"
)

var kp *kafka.Publisher

func configure() *fiber.App {
	kp = kafka.NewPublisher()

	app := fiber.New()
	app.Use(apmfiber.Middleware())
	app.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(pprof.New())

	return app
}

func Serve() {
	app := configure()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	app = loadRoute(app)

	if err := app.Listen(":" + config.Config.GetString("PORT")); err != nil {
		log.Println("Failed to start server", err.Error())
	}

	fmt.Println("Running cleanup tasks...")
}

func loadRoute(app *fiber.App) *fiber.App {
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Hello, World!") })
	app.Get("/ping", func(c *fiber.Ctx) error { return c.SendString("PONG!") })
	database := mysql.NewDatabase().ConnectNotificationDB().Execute()

	repo := notification2.NewNotificationRepository(database)
	svcNotification := notification.NewNotificationService(kp, repo)

	v1 := handlerV1.NewHandler(svcNotification)

	app.Post("/api/v1/job/push-notif", v1.CreateJob)
	app.Get("/api/v1/job/push-notif/:requestID", v1.CheckJobStatus)

	return app
}
