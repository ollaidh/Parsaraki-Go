package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	msgProcessor := NewKafkaClient("test-topic")

	telegramClient := NewTelegramClient(config, &msgProcessor)
	telegramClient.pingBot()

	// EP for Telegram Bot Webhook
	http.HandleFunc("/bot-message", telegramClient.ProcessBotMessage)

	// launch server at 8443 port
	go func() {
		if err := http.ListenAndServe(":"+config.Gateway.Port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// set webhook for bot
	telegramClient.setWebhook()

	// wait for Ctrl+C sugnal tp stop the app
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nExiting...")

}
