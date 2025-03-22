package main

import (
	"fmt"
	"time"

	"github.com/mjmhtjain/knime/src/outbox"
)

func main() {
	outboxClient := outbox.New()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

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
