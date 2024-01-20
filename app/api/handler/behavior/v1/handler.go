package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
	"notification-api/domain"
	"notification-api/service/notification"
)

type handler struct {
	service notification.Notification
}

func NewHandler(svc notification.Notification) *handler {
	return &handler{
		service: svc,
	}
}

func (h *handler) CreateJob(c *fiber.Ctx) error {
	tx := apm.DefaultTracer.StartTransaction("API POST /api/v1/job/push-notif", "CreateJob")

	ctx := apm.ContextWithTransaction(c.Context(), tx)
	span, ctx := apm.StartSpan(ctx, "handler", "CreateJob")
	defer span.End()

	var p domain.NotificationData

	if err := c.BodyParser(&p); err != nil {
		apm.CaptureError(ctx, err).Send()
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad request",
		})
	}

	resp, err := h.service.NotificationReport(ctx, &p)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		fmt.Println(err.Error())
		return c.Status(500).JSON(map[string]interface{}{
			"data": nil,
			"status": map[string]interface{}{
				"code":           1,
				"message_server": "Something Went Wrong",
				"message_client": "Something Went Wrong",
			},
		})
	}
	tx.End()
	return c.Status(200).JSON(map[string]interface{}{
		"data": resp,
		"status": map[string]interface{}{
			"code":           0,
			"message_server": "success",
			"message_client": "success",
		},
	})
}

func (h *handler) CheckJobStatus(c *fiber.Ctx) error {
	span, ctx := apm.StartSpan(c.Context(), "API Get /api/v1/job/push-notif/:requestID", "CheckJobStatus")
	defer span.End()
	requestID := c.Params("requestID")

	resp, err := h.service.GetStatus(ctx, requestID)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		fmt.Println(err.Error())
		return c.Status(500).JSON(map[string]interface{}{
			"data": nil,
			"status": map[string]interface{}{
				"code":           1,
				"message_server": "Something Went Wrong",
				"message_client": "Something Went Wrong",
			},
		})
	}
	return c.Status(200).JSON(map[string]interface{}{
		"data": resp,
		"status": map[string]interface{}{
			"code":           0,
			"message_server": "success",
			"message_client": "success",
		},
	})
}
