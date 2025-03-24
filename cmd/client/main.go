package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/outbox"
)

func main() {
	dbConfig := config.NewOutboxDBConfig(
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "postgres"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
	)
	outboxClient := outbox.New(dbConfig)

	numGoroutines := 1
	waitGroup := sync.WaitGroup{}

	for i := 0; i < numGoroutines; i++ {
		waitGroup.Add(1)
		go postMessagesJob(outboxClient, &waitGroup)
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

func postMessagesJob(outboxClient *outbox.Outbox, waitGroup *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer waitGroup.Done()

	count := 0
	for range ticker.C {
		count++
		msg := outbox.NewMessage(fmt.Sprintf("subject-%d", count), fmt.Sprintf("body-%d", count))
		outboxClient.PostMessage(msg)

		if count >= 10 {
			fmt.Println("10 seconds elapsed. Exiting.")
			return
		}
	}
}
