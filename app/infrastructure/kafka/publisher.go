package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"notification-api/config"
	"strings"
)

var KafkaPublisher *kafka.Publisher

type Publisher struct {
	*kafka.Publisher
}

func NewPublisher() *Publisher {
	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   strings.Split(config.Config.GetString("KAFKA_BROKER"), ","),
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		fmt.Println(err)
	}
	KafkaPublisher = publisher
	return &Publisher{publisher}
}
