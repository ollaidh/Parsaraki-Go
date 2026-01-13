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

var CONFIG Config

func init() {
	CONFIG = loadConfig()
}

func main() {
	telegramClient := TelegramClient{}
	telegramClient.pingBot()

	// EP for Telegram Bot Webhook
	http.HandleFunc("/bot-message", ProcessBotMessage)

	// launch server at 8443 port
	go func() {
		if err := http.ListenAndServe(":"+CONFIG.Gateway.Port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// use ngrok to get https address for dev
	if CONFIG.Mode == "DEV" {
		cmdNgrok, err := startNgrok(CONFIG.Gateway.Port)
		if err != nil {
			panic(err)
		}
		defer cmdNgrok.Process.Kill()

		urlForWebhook, err := waitForURL(5 * time.Second)
		if err != nil {
			panic(err)
		}

		fmt.Println("Ngrok public URL:", urlForWebhook)

		CONFIG.Webhooks.GatewayWebhooksUrl = urlForWebhook
	}

	// set webhook for bot
	telegramClient.setWebhook()

	// wait for Ctrl+C sugnal tp stop the app
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nExiting...")

}
