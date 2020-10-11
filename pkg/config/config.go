package config

import (
	"encoding/json"
	"os"

	env "github.com/Netflix/go-env"
	"github.com/sirupsen/logrus"
)

// APIConfig from marudor api
type APIConfig struct {
	APIEndpoint string `env:"API_ENDPOINT" json:"endpoint"`
}

// TelegramConfig for the bot
type TelegramConfig struct {
	Key string `env:"TELEGRAM_KEY" json:"key"`
}

// DatabaseConfig connection data
type DatabaseConfig struct {
	Dialect string `env:"DB_DIALECT" json:"dialect"`
	Path    string `env:"DB_PATH" json:"path"`
}

// WebhookConfig connection data
type WebhookConfig struct {
	Endpoint string `env:"WEBHOOK_ENDPOINT" json:"endpoint"`
	Port     string `env:"WEBHOOK_PORT" json:"port"`
}

// Config contains the complete service configuration
type Config struct {
	APIConfig     APIConfig      `json:"api"`
	Database      DatabaseConfig `json:"database"`
	Telegram      TelegramConfig `json:"telegram"`
	LogLevel      string         `env:"LOG_LEVEL" json:"loglevel"`
	EnableMetrics bool           `json:"enable_metrics"`
	Webhook       WebhookConfig  `json:"webhook"`
}

// NewTestConfig return a config object with test settings
func NewTestConfig() *Config {
	return &Config{
		APIConfig: APIConfig{
			APIEndpoint: "https://marudor.de/api",
		},
		LogLevel: "trace",
		Webhook: WebhookConfig{
			Endpoint: "http://localhost/hook",
		},
	}
}

// ReadConfig reads a json file and overwrite with ENV vars
func ReadConfig(file string, log *logrus.Entry) *Config {
	var config Config

	if fileExists(file) {
		fileContent, _ := os.Open(file)

		if err := json.NewDecoder(fileContent).Decode(&config); err != nil {
			log.Fatal(err)
		}
	}

	// Override ENVs
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Telegram.Key == "" {
		panic("Need TELEGRAM_KEY")
	}

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	if config.APIConfig.APIEndpoint == "" {
		config.APIConfig.APIEndpoint = "https://marudor.de/api"
	}

	if config.Webhook.Endpoint == "" {
		panic("Need WEBHOOK_ENDPOINT")
	}

	if config.Webhook.Port == "" {
		config.Webhook.Port = ":8080"
	}

	return &config
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
