package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mjmhtjain/knime/src/outbox"
)

func main() {
	outboxClient := outbox.New()

	numGoroutines := 3
	waitGroup := sync.WaitGroup{}

	for i := 0; i < numGoroutines; i++ {
		waitGroup.Add(1)
		go generateMessages(outboxClient, &waitGroup)
	}

	waitGroup.Wait()
}

func generateMessages(outboxClient *outbox.Outbox, waitGroup *sync.WaitGroup) {
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
