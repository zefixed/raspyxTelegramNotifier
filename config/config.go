package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App   App
		Log   Log
		PG    PG
		Bot   Bot
		Kafka Kafka
	}
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
		Type  string `env:"LOG_TYPE,required"`
	}

	PG struct {
		PGURL    string `env:"PG_URL,required"`
		Timeout  string `env:"PG_TIMEOUT,required"`
		Attempts int    `env:"PG_ATTEMPTS,required"`
	}

	Bot struct {
		Token string `env:"BOT_TOKEN,required"`
	}

	Kafka struct {
		Port      string `env:"KAFKA_PORT,required"`
		URL       string `env:"KAFKA_URL,required"`
		TopicName string `env:"KAFKA_TOPIC_NAME,required"`
		Group     string `env:"KAFKA_GROUP,required"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
