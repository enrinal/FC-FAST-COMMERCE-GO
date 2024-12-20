package config

import "github.com/segmentio/kafka-go"

var KafkaBrokerURLs = []string{"localhost:9092", "localhost:9093", "localhost:9094"}

func NewKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(KafkaBrokerURLs...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{}, // Balancer for selecting partition
		AllowAutoTopicCreation: true,
	}
}
