package main

import (
	"log"
	"net/http"
)

var CONFIG Config

func init() {
	CONFIG = loadConfig()
}

func main() {
	pingBot()

	http.HandleFunc("/bot-message", ProcessBotMessage)

	go func() {
		if err := http.ListenAndServe(":8443", nil); err != nil {
			log.Fatal(err)
		}
	}()

	setWebhook()
	select {}

}
