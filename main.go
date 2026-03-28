package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"parsaraki-go/config"
	telegramemulator "parsaraki-go/external/telegram_emulator"
	telegramapi "parsaraki-go/internal/app/api/telegram"
	msgconsumer "parsaraki-go/internal/app/consumers"
	msgproducer "parsaraki-go/internal/infrastructure/kafka"
	inmemoryrepo "parsaraki-go/internal/repository/memory"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// telegramClient := telegram.NewTelegramClient(&config)
	msgProducer := msgproducer.NewKafkaProducer("all-messages")

	telegramHandler := telegramapi.NewTelegramHandler(&msgProducer, config)

	// endpoint for Telegram Bot Webhook
	http.HandleFunc("/bot-message", telegramHandler.ProcessBotMessage)

	// launch server at 8443 port
	go func() { // HOW to gracefully shutdown this??
		if err := http.ListenAndServe(":"+config.Gateway.Port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		tgMsgEmulator := telegramemulator.NewTelegramMessagesEmulator()
		tgMsgEmulator.SendMessages(ctx, "http://127.0.0.1:8443/bot-message")
	}()

	// Create and run consumer

	repo := inmemoryrepo.NewMemoryDB()

	consumer := msgconsumer.NewKafkaConsumer("all-messages", &repo)
	consumer.Consume(ctx)
	// ADD server shutdown

	fmt.Println("Exiting...")

}
