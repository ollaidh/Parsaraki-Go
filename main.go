package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	telegramClient := NewTelegramClient(config)
	telegramClient.pingBot()

	// EP for Telegram Bot Webhook
	http.HandleFunc("/bot-message", telegramClient.ProcessBotMessage)

	// launch server at 8443 port
	go func() {
		if err := http.ListenAndServe(":"+config.Gateway.Port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// use ngrok to get https address for dev
	if config.Mode == "DEV" {
		cmdNgrok, err := startNgrok(config.Gateway.Port)
		if err != nil {
			panic(err)
		}
		defer cmdNgrok.Process.Kill()

		urlForWebhook, err := waitForURL(5 * time.Second)
		if err != nil {
			panic(err)
		}

		fmt.Println("Ngrok public URL:", urlForWebhook)

		telegramClient.config.Webhooks.GatewayWebhooksUrl = urlForWebhook
	}

	// set webhook for bot
	telegramClient.setWebhook()

	// wait for Ctrl+C sugnal tp stop the app
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nExiting...")

}
