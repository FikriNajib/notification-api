package push_notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.elastic.co/apm/v2"
	"io/ioutil"
	"net/http"
	"notification-api/config"
	"notification-api/domain"
)

type repository struct {
}

func NewClevertapRepository() PushNotifRepository {
	return &repository{}
}
func (r *repository) PostPushNotif(ctx context.Context, request domain.PushNotificationRequest) (domain.ResponseNotif, error) {
	span, ctx := apm.StartSpan(ctx, "repositories/push_notification/push_notification.go", "PostPushNotif")
	defer span.End()
	url := config.Config.GetString("CLEVERTAP_URL")
	method := "POST"

	jsonPayload, err := json.Marshal(request)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		fmt.Println(err)
		return domain.ResponseNotif{}, err
	}
	fmt.Println("payload: ", string(jsonPayload))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonPayload))

	if err != nil {
		apm.CaptureError(ctx, err).Send()
		fmt.Println(err)
		return domain.ResponseNotif{}, err
	}
	req.Header.Add("X-CleverTap-Account-Id", config.Config.GetString("CLEVERTAP_ACCOUNT_ID"))
	req.Header.Add("X-CleverTap-Passcode", config.Config.GetString("CLEVERTAP_PASSCODE"))

	res, err := client.Do(req)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		fmt.Println(err)
		return domain.ResponseNotif{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		fmt.Println(err)
		return domain.ResponseNotif{}, err
	}
	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err)
		apm.CaptureError(ctx, err).Send()
		return domain.ResponseNotif{}, err
	}
	resp := domain.ResponseNotif{}
	resp.Status = data["status"].(string)
	resp.ErrorMessage = string(body)

	return resp, nil
}
