package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sezzle/backend/internal/adapters/in/httpapi"
	"sezzle/backend/internal/adapters/out/strategyregistry"
	"sezzle/backend/internal/application/calculator"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	strategyRegistry := strategyregistry.NewInMemoryCalculationStrategyRegistry()
	calculatorService := calculator.NewService(strategyRegistry)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           httpapi.NewCalculatorHandler(calculatorService),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("calculator API started", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case signal := <-stop:
		slog.Info("shutdown signal received", "signal", signal.String())
		shutdown(server)
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server stopped unexpectedly", "error", err)
		}
	}
}

func shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		return
	}
	slog.Info("calculator API stopped")
}
