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
	"github.com/sirupsen/logrus"
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

	// Launch the outbox service to consume messages from the outbox
	outboxClient := outbox.New(dbConfig, natsConfig)
	go outboxClient.LaunchOutboxService(ctx)

	// Launch the message posting service to post messages to the outbox
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	messageService := NewMessageService(dbConfig, natsConfig)
	go func() {
		defer waitGroup.Done()
		messageService.LaunchPostMessageJob(ctx)
	}()

	waitGroup.Wait()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type IMessageService interface {
	LaunchPostMessageJob(ctx context.Context)
}

type MessageService struct {
	outboxClient *outbox.Outbox
}

func NewMessageService(dbConfig *config.OutboxDBConfig, natsConfig *config.NatsConfig) IMessageService {
	outboxClient := outbox.New(dbConfig, natsConfig)

	return &MessageService{
		outboxClient: outboxClient,
	}
}

func (s *MessageService) LaunchPostMessageJob(ctx context.Context) {
	count := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping message posting...")
			return
		case <-time.Tick(1 * time.Second):
			count++
			msg := outbox.NewMessage(fmt.Sprintf("subject-%d", count), fmt.Sprintf("body-%d", count))
			err := s.outboxClient.PostMessage(msg)
			if err != nil {
				logrus.Error("Error posting message:", "error", err)
			}
		}
	}
}
