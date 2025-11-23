package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type EurekaInstance struct {
	InstanceID string `json:"instanceId"`
	HostName   string `json:"hostName"`
	App        string `json:"app"`
	IPAddr     string `json:"ipAddr"`
	Status     string `json:"status"`
	Port       struct {
		Port    int  `json:"$"`
		Enabled bool `json:"@enabled"`
	} `json:"port"`
	SecurePort struct {
		Port    int  `json:"$"`
		Enabled bool `json:"@enabled"`
	} `json:"securePort"`
	DataCenterInfo struct {
		Class string `json:"@class"`
		Name  string `json:"name"`
	} `json:"dataCenterInfo"`
	LeaseInfo struct {
		RenewalIntervalInSecs int `json:"renewalIntervalInSecs"`
		DurationInSecs        int `json:"durationInSecs"`
	} `json:"leaseInfo"`
}

type EurekaRegistration struct {
	Instance EurekaInstance `json:"instance"`
}

func RegisterWithEureka(eurekaURL string, config *Config) error {
	instance := EurekaInstance{
		InstanceID: fmt.Sprintf("%s:%s:%d",
			config.Eureka.Instance.Hostname,
			config.Eureka.Instance.App,
			config.Eureka.Instance.Port),
		HostName: config.Eureka.Instance.Hostname,
		App:      config.Eureka.Instance.App,
		IPAddr:   getLocalIP(),
		Status:   "UP",
		Port: struct {
			Port    int  `json:"$"`
			Enabled bool `json:"@enabled"`
		}{
			Port:    config.Eureka.Instance.Port,
			Enabled: true,
		},
		SecurePort: struct {
			Port    int  `json:"$"`
			Enabled bool `json:"@enabled"`
		}{
			Port:    443,
			Enabled: false,
		},
		DataCenterInfo: struct {
			Class string `json:"@class"`
			Name  string `json:"name"`
		}{
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Name:  "MyOwn",
		},
		LeaseInfo: struct {
			RenewalIntervalInSecs int `json:"renewalIntervalInSecs"`
			DurationInSecs        int `json:"durationInSecs"`
		}{
			RenewalIntervalInSecs: 30,
			DurationInSecs:        90,
		},
	}

	registration := EurekaRegistration{Instance: instance}
	payload, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("failed to marshal registration: %w", err)
	}

	url := fmt.Sprintf("%s/apps/%s", eurekaURL, config.Eureka.Instance.App)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send registration request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("eureka registration failed with status: %d", resp.StatusCode)
	}

	log.Printf("Successfully registered with Eureka: %s", url)

	// Запускаем heartbeat в фоне
	go sendHeartbeat(eurekaURL, config)

	return nil
}

func sendHeartbeat(eurekaURL string, config *Config) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	instanceID := fmt.Sprintf("%s:%s:%d",
		config.Eureka.Instance.Hostname,
		config.Eureka.Instance.App,
		config.Eureka.Instance.Port)

	for range ticker.C {
		url := fmt.Sprintf("%s/apps/%s/%s",
			eurekaURL, config.Eureka.Instance.App, instanceID)

		req, err := http.NewRequest("PUT", url, nil)
		if err != nil {
			log.Printf("Failed to create heartbeat request: %v", err)
			continue
		}

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to send heartbeat: %v", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Heartbeat failed with status: %d", resp.StatusCode)
		}
	}
}

func getLocalIP() string {
	// В реальном приложении здесь должна быть логика получения реального IP
	return "127.0.0.1"
}
