package msgproducer

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(topic string) KafkaProducer {
	return KafkaProducer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (kp *KafkaProducer) Close() error {
	if err := kp.writer.Close(); err != nil {
		return err
	}
	return nil
}

func (kc *KafkaProducer) WriteMessage(message string) {
	err := kc.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(message),
		},
	)
	if err != nil {
		fmt.Printf("Kafka send error: %s", err)
	}
}
