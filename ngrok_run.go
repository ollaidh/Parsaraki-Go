package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

type Tunnel struct {
	PublicURL string `json:"public_url"`
}

type TunnelsResponse struct {
	Tunnels []Tunnel `json:"tunnels"`
}

func startNgrok(port string) (*exec.Cmd, error) {
	cmd := exec.Command("ngrok", "http", port)
	return cmd, cmd.Start()
}

func getPublicURL() (string, error) {
	resp, err := http.Get("http://127.0.0.1:4040/api/tunnels")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tr TunnelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	if len(tr.Tunnels) == 0 {
		return "", fmt.Errorf("no tunnels")
	}
	return tr.Tunnels[0].PublicURL, nil
}

func waitForURL(timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if url, err := getPublicURL(); err == nil {
			return url, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return "", fmt.Errorf("timeout waiting for ngrok")
}
