// pinger/main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ContainerStatus struct {
	IPAddress      string    `json:"ip_address"`
	PingDuration   int64     `json:"ping_duration"`
	LastSuccessful time.Time `json:"last_successful"`
}

func pingIP(ip string) (int64, error) {
	cmd := exec.Command("ping", "-c", "1", ip)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	outputStr := out.String()
	idx := strings.Index(outputStr, "time=")
	if idx == -1 {
		return 0, fmt.Errorf("ping output parsing error")
	}
	subStr := outputStr[idx+5:]
	fields := strings.Fields(subStr)
	if len(fields) < 1 {
		return 0, fmt.Errorf("ping output parsing error")
	}
	timeStr := strings.TrimSuffix(fields[0], "ms")
	durationFloat, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		return 0, err
	}
	return int64(durationFloat), nil
}

func sendStatus(status ContainerStatus, backendURL string) error {
	jsonData, err := json.Marshal(status)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/containers", backendURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post status, code: %d", resp.StatusCode)
	}
	return nil
}

func main() {
	targetIPs := os.Getenv("TARGET_IPS")
	if targetIPs == "" {
		log.Fatal("TARGET_IPS environment variable not set")
	}
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://backend:8080" 
	}

	ips := strings.Split(targetIPs, ",")
	intervalStr := os.Getenv("PING_INTERVAL")
	interval := 30 * time.Second
	if intervalStr != "" {
		if val, err := strconv.Atoi(intervalStr); err == nil {
			interval = time.Duration(val) * time.Second
		}
	}
	log.Printf("Starting pinger, will ping %v every %v", ips, interval)
	for {
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			start := time.Now()
			duration, err := pingIP(ip)
			if err != nil {
				log.Printf("Ping failed for %s: %v", ip, err)
				continue
			}
			log.Printf("Ping %s: %d ms", ip, duration)
			status := ContainerStatus{
				IPAddress:      ip,
				PingDuration:   duration,
				LastSuccessful: time.Now(),
			}
			if err = sendStatus(status, backendURL); err != nil {
				log.Printf("Failed to send status for %s: %v", ip, err)
			}
		}
		time.Sleep(interval)
	}
}

