package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		slog.Error("Error loading .env file.")
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         os.Getenv("HOST") + ":" + os.Getenv("PORT"),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// Spawn a goroutine to make the function non blocking
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
		}
	}()

	slog.Info("Started server on", "port", os.Getenv("PORT"))

	// Block until SIGINT or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info(fmt.Sprintf("Received signal: %s", sig))

	// Give active connections 60 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel() // Waits for the srv.Shutdown function to finish before executing

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Forced shutdown", "err", err)
	}

	slog.Info("Served Shutdown.")
}
