package config

import (
	"github.com/sahalazain/go-common/config"
)

var DefaultConfig = map[string]interface{}{
	"PORT":                            "8008",
	"KAFKA_BROKER":                    "localhost:29092",
	"HYSTRIX_TIMEOUT":                 1000,
	"HYSTRIX_MAX_CONCURRENT_REQUESTS": 100,
	"HYSTRIX_ERROR_PERCENT_THRESHOLD": 25,
	"SENTRY_DSN":                      "",
	"SENTRY_TRACES_SAMPLE_RATE":       0.1,
	"SENTRY_ENVIRONMENT":              "development",
	"SENTRY_DEBUG":                    false,
	"TOPIC":                           "dev_notification",
	"KAFKA_CONSUMER_GROUP":            "dev_notification_group",
	"MYSQL_NOTIFICATION_USER":         "root",
	"MYSQL_NOTIFICATION_PASS":         "",
	"MYSQL_NOTIFICATION_HOST":         "localhost",
	"MYSQL_NOTIFICATION_PORT":         "3306",
	"MYSQL_NOTIFICATION_DBNAME":       "local_notification",
	"CLEVERTAP_URL":                   "",
	"CLEVERTAP_ACCOUNT_ID":            "",
	"CLEVERTAP_PASSCODE":              "",
	"SENDGRID_ENDPOINT":               "",
	"SENDGRID_HOST":                   "",
	"SENDGRID_API_KEY":                "",
	"MAIL_SENDER_NAME":                "",
	"MAIL_SENDER_ADDRESS":             "",
	"ALIBABA_SMS_ACCESS_KEY_ID":	   "",
	"ALIBABA_SMS_ACCESS_KEY_SECRET":   "",

}

var Config config.Getter
var Url string

func Load() error {
	cfgClient, err := config.Load(DefaultConfig, Url)
	if err != nil {
		return err
	}

	Config = cfgClient

	return nil
}
