package notification

import (
	"context"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.elastic.co/apm/v2"
	"notification-api/config"
	"notification-api/domain"
)

type repository struct {
	personalization *mail.Personalization
	email           *mail.Email
	request         domain.MetadataEmail
}

func NewSendgridRepository() EmailRepository {
	return &repository{}
}

func (r *repository) SetPersonalization() *repository {
	p := mail.NewPersonalization()
	r.personalization = p

	return r
}

func (r *repository) SetEmailTo(identity string) *repository {
	to := mail.NewEmail("Client", identity)
	r.personalization.AddTos(to)

	return r
}

func (r *repository) SetTemplateData(req domain.MetadataEmail, senderName, fromAdmin string) *mail.SGMailV3 {
	from := mail.NewEmail(senderName, fromAdmin)
	m := mail.NewV3Mail()
	m.SetFrom(from)
	m.AddPersonalizations(r.personalization)
	for key, _ := range req.TemplateParam {
		result, _ := req.TemplateParam[key].(string)
		r.personalization.SetDynamicTemplateData(key, result)
		m.SetTemplateID(req.TemplateName)
	}
	return m
}

func (r *repository) PostEmail(ctx context.Context, req domain.MetadataEmail) (domain.ResponseNotif, error) {
	span, ctx := apm.StartSpan(ctx, "repositories/email/sendgrid.go", "PostEmail")
	defer span.End()
	apiKey := config.Config.GetString("SENDGRID_API_KEY")

	m := r.SetPersonalization().
		SetEmailTo(req.Identity).
		SetTemplateData(
			req,
			config.Config.GetString("MAIL_SENDER_NAME"),
			config.Config.GetString("MAIL_SENDER_ADDRESS"))

	body := mail.GetRequestBody(m)
	request := sendgrid.GetRequest(apiKey, config.Config.GetString("SENDGRID_ENDPOINT"), config.Config.GetString("SENDGRID_HOST"))
	request.Method = "POST"
	request.Body = body

	_, err := sendgrid.API(request)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		return domain.ResponseNotif{Status: "failed", ErrorMessage: err.Error()}, err
	}

	result := domain.ResponseNotif{
		Status:       "success",
		ErrorMessage: "",
	}

	return result, nil
}
