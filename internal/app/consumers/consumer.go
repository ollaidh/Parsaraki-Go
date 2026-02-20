package msgconsumer

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
		return kafka.Message{}, err
	}
	return msg, nil
}

func (kc *KafkaConsumer) RunConsumer() {
	for {
		msg, err := kc.readMessage()
		if err != nil {
			fmt.Printf("Error reading message from Kafka: %s", err)
			continue
		}
		fmt.Printf("\nGot message from Kafka: %s, topic: %s", msg.Value, msg.Topic)
	}
}
