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
	pingBot()

	// EP for Telegram Bot Webhook
	http.HandleFunc("/bot-message", ProcessBotMessage)

	go func() {
		if err := http.ListenAndServe(":8443", nil); err != nil {
			log.Fatal(err)
		}
	}()

	cmdNgrok, err := startNgrok("8443")
	if err != nil {
		panic(err)
	}
	defer cmdNgrok.Process.Kill()

	urlForWebhook, err := waitForURL(5 * time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Println("Ngrok public URL:", urlForWebhook)

	CONFIG.TelegramBot.WebhooksUrl = urlForWebhook

	setWebhook()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nExiting...")

}
