package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Eureka   EurekaConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

type RedisConfig struct {
	Host string
	Port string
}

type EurekaConfig struct {
	URL string
}

type AuthConfig struct {
	URL string
}

func LoadConfig() Config {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Значения по умолчанию
	viper.SetDefault("server.port", "8083")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
