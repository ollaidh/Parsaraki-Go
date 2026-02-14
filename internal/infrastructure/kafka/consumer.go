package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(topic string) KafkaConsumer {
	return KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			GroupID: "my-group",
		}),
	}
}

func (kc *KafkaConsumer) Close() error {
	if err := kc.reader.Close(); err != nil {
		return err
	}
	return nil
}

func (kc *KafkaConsumer) readMessage() (kafka.Message, error) {
	msg, err := kc.reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Printf("Kafka read error: %s", err)
		return kafka.Message{}, err
	}

	fmt.Printf("Kafka read message: %s, topic: %s \n", msg.Value, msg.Topic)

	return msg, nil
}
