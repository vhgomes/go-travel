package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vhgomes/go-travel/internal/services/flight"
	"github.com/vhgomes/go-travel/internal/services/hotel"
	"github.com/vhgomes/go-travel/internal/services/order"
	sqssvc "github.com/vhgomes/go-travel/internal/services/sqs"
	"github.com/vhgomes/go-travel/pkg/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("pinging postgres: %w", err)
	}
	log.Println("connected to postgres")

	awsCfg := config.LoadFromEnv()
	sqsClient, err := config.NewSQSClient(ctx, awsCfg)
	if err != nil {
		return fmt.Errorf("creating sqs client: %w", err)
	}

	queueManager := sqssvc.NewQueueManager(sqsClient)

	sagaQueueName := getEnv("SAGA_QUEUE_NAME", "saga-start.fifo")
	sagaQueueURL, err := queueManager.GetQueueURL(ctx, sagaQueueName)
	if err != nil {
		return fmt.Errorf("resolving saga queue url (%s): %w", sagaQueueName, err)
	}

	producer := sqssvc.NewProducer(sqsClient, sagaQueueURL)

	orderRepo := order.NewOrderRepository(pool)

	flightRepo := flight.NewInMemoryFlightRepository()
	flightService := flight.NewFlightService(flightRepo)

	hotelRepo := hotel.NewInMemoryHotelRepository()
	hotelService := hotel.NewHotelService(hotelRepo)

	orderService := order.NewOrderService(orderRepo, flightService, hotelService, producer)
	orderHandler := order.NewOrderHandler(orderService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /healthz", healthCheck(pool))

	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("TravelGo API listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutting down http server: %w", err)
	}

	log.Println("server stopped gracefully")
	return nil
}

func healthCheck(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
