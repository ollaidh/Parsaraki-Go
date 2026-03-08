package msgconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"parsaraki-go/internal/repository"

	"parsaraki-go/internal/infrastructure/telegram"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	repo   repository.Repository
}

func NewKafkaConsumer(topic string, repo repository.Repository) KafkaConsumer {
	return KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			GroupID: "my-group",
		},
		),
		repo: repo,
	}
}

func (kc *KafkaConsumer) Close() error {
	if err := kc.reader.Close(); err != nil {
		return err
	}
	return nil
}

func (kc *KafkaConsumer) readMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := kc.reader.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}
	return msg, nil
}

func (kc *KafkaConsumer) RunConsumer(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			break
		}

		// ??? в readMessage нужно передавать другой контекст с таймаутом? Где его создавать?
		msg, err := kc.readMessage(ctx)
		if err != nil { // not handling context cancellation here, just logging error and continue
			fmt.Printf("Error reading message from Kafka: %s", err)
			continue
		}

		fmt.Printf("\nGot message from Kafka: %s, topic: %s", msg.Value, msg.Topic)

		var botMsg telegram.BotMessage
		err = json.Unmarshal(msg.Value, &botMsg)
		kc.repo.SaveBotRequest(botMsg)
	}
	log.Println("consumer stopped")
}
