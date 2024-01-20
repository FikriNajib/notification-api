package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"notification-api/app/consumer"
	"notification-api/app/infrastructure/kafka"
	"notification-api/config"
)

var kp *kafka.Publisher

var (
	consumerCmd = &cobra.Command{
		Use:              "consumer",
		Short:            "A Cronjob for consume data job to push notification",
		Long:             "A Cronjob for consume data job to push notification",
		PersistentPreRun: consumerPreRun,
		RunE:             runConsumer,
	}
)

func ConsumerCmd() *cobra.Command {
	return consumerCmd
}

func consumerPreRun(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		config.Url = args[0]
	}
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
}

func runConsumer(cmd *cobra.Command, args []string) error {
	kp = kafka.NewPublisher()
	go consumer.NewConsumer().InitRouter().Listen().Serve()
	return nil
}
