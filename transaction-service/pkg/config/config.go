package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `yaml:"port" json:"port"`
	} `yaml:"server" json:"server"`

	Database struct {
		PostgresURL string `yaml:"postgres-url" json:"postgres-url"`
		RedisURL    string `yaml:"redis-url" json:"redis-url"`
	} `yaml:"database" json:"database"`

	Eureka struct {
		URL      string `yaml:"url" json:"url"`
		Instance struct {
			Hostname string `yaml:"hostname" json:"hostname"`
			App      string `yaml:"app" json:"app"`
			Port     int    `yaml:"port" json:"port"`
		} `yaml:"instance" json:"instance"`
	} `yaml:"eureka" json:"eureka"`

	Services struct {
		Auth    string `yaml:"auth" json:"auth"`
		Deposit string `yaml:"deposit" json:"deposit"`
	} `yaml:"services" json:"services"`

	Security struct {
		OAuth2 struct {
			ClientID     string `yaml:"client-id" json:"client-id"`
			ClientSecret string `yaml:"client-secret" json:"client-secret"`
			TokenURL     string `yaml:"token-url" json:"token-url"`
		} `yaml:"oauth2" json:"oauth2"`
	} `yaml:"security" json:"security"`

	Logging struct {
		Level string `yaml:"level" json:"level"`
	} `yaml:"logging" json:"logging"`

	Cache struct {
		TTL string `yaml:"ttl" json:"ttl"`
	} `yaml:"cache" json:"cache"`
}

var AppConfig *Config

func InitConfig() error {
	// Сначала пробуем загрузить из Spring Cloud Config
	configServerURL := os.Getenv("CONFIG_SERVER_URL")
	if configServerURL == "" {
		configServerURL = "http://localhost:8888"
	}

	appName := "transaction-service"
	profile := os.Getenv("SPRING_PROFILES_ACTIVE")
	if profile == "" {
		profile = "default"
	}

	config, err := LoadConfigFromConfigServer(configServerURL, appName, profile)
	if err != nil {
		log.Printf("Failed to load from config server: %v. Using local config.", err)
		return loadLocalConfig()
	}

	AppConfig = config
	log.Printf("Successfully loaded config from config server for profile: %s", profile)
	return nil
}

func LoadConfigFromConfigServer(configServerURL, appName, profile string) (*Config, error) {
	url := fmt.Sprintf("%s/%s/%s", configServerURL, appName, profile)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to config server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("config server returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var config Config
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &config, nil
}

func loadLocalConfig() error {
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./pkg/config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read local config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal local config: %w", err)
	}

	AppConfig = &config
	log.Println("Successfully loaded local config")
	return nil
}

func GetConfig() *Config {
	if AppConfig == nil {
		if err := InitConfig(); err != nil {
			log.Fatalf("Failed to initialize config: %v", err)
		}
	}
	return AppConfig
}
