package main

import (
	"log"
	"raspyxTelegramNotifier/config"
	"raspyxTelegramNotifier/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
