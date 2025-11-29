package main

import (
	"fmt"
	"health_checker/config"
	"health_checker/internal/app"
	"os"
)

func main() {
	// Загрузка конфига
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config %v\n", err)
		os.Exit(1)
	}

	// Запуск сервера
	if err := app.StartServer(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server %v\n", err)
		os.Exit(1)
	}
}
