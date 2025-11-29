package app

import (
	"context"
	"health_checker/config"
	"health_checker/internal/delivery"
	"health_checker/internal/repository"
	"health_checker/internal/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func StartServer(cfg *config.Config) error {
	logger := logrus.New()

	storage := repository.NewRepository()

	// Services
	healthCheckerSvc := services.NewHealthCheckerService(
		logger,
		storage,
		cfg,
	)

	pdfGeneratorSvc := services.NewPdfGeneratorService(
		logger,
	)

	httpServer := delivery.NewHTTPHandler(logger, healthCheckerSvc, pdfGeneratorSvc)

	// Настройка HTTP сервера
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      httpServer.NewRouters().Handler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Ожидание сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server stopped")
	return nil
}
