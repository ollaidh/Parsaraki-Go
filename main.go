package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"parsaraki-go/config"
	telegramapi "parsaraki-go/internal/app/api/telegram"
	"parsaraki-go/internal/infrastructure/kafka"
	"syscall"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// telegramClient := telegram.NewTelegramClient(&config)
	msgProducer := kafka.NewKafkaProducer("all-messages")

	telegramHandler := telegramapi.NewTelegramHandler(&msgProducer, config)

	// endpoint for Telegram Bot Webhook
	http.HandleFunc("/bot-message", telegramHandler.ProcessBotMessage)

	// launch server at 8443 port
	go func() {
		if err := http.ListenAndServe(":"+config.Gateway.Port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// wait for Ctrl+C signal tp stop the app
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nExiting...")

}
