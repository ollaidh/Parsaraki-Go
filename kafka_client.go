package main

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaClient(topic string) KafkaClient {
	return KafkaClient{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			GroupID: "my-group",
		}),
	}
}

func (kc *KafkaClient) Close() error {
	if err := kc.writer.Close(); err != nil {
		return err
	}
	if err := kc.reader.Close(); err != nil {
		return err
	}
	return nil
}

func (kc *KafkaClient) writeMessage(message string) {
	err := kc.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(message),
		},
	)
	if err != nil {
		fmt.Printf("Kafka send error: %s", err)
	}
}

func (kc *KafkaClient) readMessage() (kafka.Message, error) {
	msg, err := kc.reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Printf("Kafka read error: %s", err)
		return kafka.Message{}, err
	}

	fmt.Printf("Kafka read message: %s, topic: %s \n", msg.Value, msg.Topic)

	return msg, nil
}

func (kc *KafkaClient) ProcessMsg(content string) string {
	kc.writeMessage(content)
	msgFromKafka, _ := kc.readMessage()
	msgResponseToBot := fmt.Sprintf("Received from Kafka (topic: %s), message: %s", msgFromKafka.Topic, string(msgFromKafka.Value))

	return msgResponseToBot
}
