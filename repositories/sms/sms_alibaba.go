package sms

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"go.elastic.co/apm"
	"notification-api/config"
	"notification-api/domain"
)

type repository struct {
}

func NewSMSRepository() SmsRepository {
	return &repository{}
}

func (r *repository) SendSMS(ctx context.Context, recipient, otp, otpType string) (domain.ResponseNotif, error) {
	span, ctx := apm.StartSpan(ctx, "repositories/sms/sms_alibaba.go", "SendSMS")
	defer span.End()

	var msgJSON string
	var result domain.ResponseNotif

	configSMS := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(config.Config.GetString("ALIBABA_SMS_ACCESS_KEY_ID"), config.Config.GetString("ALIBABA_SMS_ACCESS_KEY_SECRET"))
	client, err := sdk.NewClientWithOptions("ap-southeast-1", configSMS, credential)
	if err != nil {
		result.Status = "failed"
		result.ErrorMessage = err.Error()
		apm.CaptureError(ctx, err).Send()
		return result, err
	}

	switch otpType {
	case "registration":
		msgJSON = fmt.Sprintf("Hi, Insert OTP Code to verify your account at RCTI+. Your OTP Code is %s", otp)
	case "forget-password":
		msgJSON = fmt.Sprintf("Hi, Insert OTP Code to create new password at RCTI+. Your OTP Code is %s", otp)
	case "change-password":
		msgJSON = fmt.Sprintf("Hi, Insert OTP Code to create new password at RCTI+. Your OTP Code is %s", otp)
	case "edit-profile":
		msgJSON = fmt.Sprintf("Hi, Insert OTP Code to edit profile data at RCTI+. Your OTP Code is %s", otp)
	case "delete-profile":
		msgJSON = fmt.Sprintf("Hi, Insert OTP Code to delete profile at RCTI+. Your OTP Code is %s", otp)
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = ""
	request.Version = "2018-05-01"
	request.ApiName = ""
	request.QueryParams["To"] = recipient
	request.QueryParams["From"] = ""
	request.QueryParams["Type"] = "OTP"
	request.QueryParams["Message"] = msgJSON
	request.RegionId = "ap-southeast-1"
	request.AcceptFormat = "JSON"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		result.Status = "failed"
		result.ErrorMessage = err.Error()
		apm.CaptureError(ctx, err).Send()
		return result, err
	}
	result.Status = "success"
	result.ErrorMessage = ""
	fmt.Print(response.GetHttpContentString())
	return result, nil
}
