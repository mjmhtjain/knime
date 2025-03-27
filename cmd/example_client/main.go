package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

	// Initialize the configs
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

type IMessageService interface {
	LaunchPostMessageJob(ctx context.Context)
}

type MessageService struct {
	outboxClient        *outbox.Outbox
	postMessageInterval int
}

func NewMessageService(dbConfig *config.OutboxDBConfig, natsConfig *config.NatsConfig) IMessageService {
	postMessageInterval := getEnv("POST_MESSAGE_INTERVAL", "1")
	interval, err := strconv.Atoi(postMessageInterval)
	if err != nil {
		logrus.Fatalf("Error converting POST_MESSAGE_INTERVAL to int: %v", err)
	}

	return &MessageService{
		postMessageInterval: interval,
		outboxClient:        outbox.New(dbConfig, natsConfig),
	}
}

func (s *MessageService) LaunchPostMessageJob(ctx context.Context) {
	count := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping message posting...")
			return
		case <-time.Tick(time.Duration(s.postMessageInterval) * time.Second): // Post a message every interval seconds
			count++
			msg := outbox.NewMessage(fmt.Sprintf("subject-%d", count), fmt.Sprintf("body-%d", count))
			err := s.outboxClient.PostMessage(msg)
			if err != nil {
				logrus.Error("Error posting message:", "error", err)
			}
		}
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
