package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"parsaraki-go/config"
	telegramapi "parsaraki-go/internal/app/api/telegram"
	msgconsumer "parsaraki-go/internal/app/consumers"
	msgproducer "parsaraki-go/internal/infrastructure/kafka"
	"syscall"
	"time"
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
	go func() {
		if err := http.ListenAndServe(":"+config.Gateway.Port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Create and run consumer

	// repo := repository.NewMemoryDB()

	consumer := msgconsumer.NewKafkaConsumer("all-messages")
	defer consumer.Close()

	go func() {
		consumer.RunConsumer(ctx)
	}()

	// wait for Ctrl+C signal tp stop the app
	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// <-sigCh

	<-ctx.Done() // waiting shutdown

	// ADD wait group wg, pass to consumer

	// Add channel
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		log.Println("shutdown timeout")
	}

	log.Println("consumer stopped")

	// ADD server shutdown
	// ADD Consumer shutdown

	fmt.Println("\nExiting...")

}
