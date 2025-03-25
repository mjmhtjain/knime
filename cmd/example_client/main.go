package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/outbox"
)

func main() {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown signals
	go func() {
		<-sigChan
		fmt.Println("Received shutdown signal, initiating graceful shutdown...")
		cancel()
	}()

	dbConfig := config.NewOutboxDBConfig(
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "postgres"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
	)

	natsConfig := config.NewNatsConfig(
		getEnv("NATS_URL", "nats://localhost:4222"),
	)

	outboxClient := outbox.New(dbConfig, natsConfig)
	go outboxClient.LaunchOutboxService(ctx)

	numGoroutines := 1
	waitGroup := sync.WaitGroup{}

	for i := 0; i < numGoroutines; i++ {
		waitGroup.Add(1)
		go postMessagesJob(ctx, outboxClient, &waitGroup)
	}

	waitGroup.Wait()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func postMessagesJob(ctx context.Context, outboxClient *outbox.Outbox, waitGroup *sync.WaitGroup) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	defer waitGroup.Done()

	count := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping message posting...")
			return
		case <-ticker.C:
			count++
			msg := outbox.NewMessage(fmt.Sprintf("subject-%d", count), fmt.Sprintf("body-%d", count))
			outboxClient.PostMessage(msg)
		}
	}
}
