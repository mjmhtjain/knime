package main

import "github.com/mjmhtjain/knime/src/outbox"

func main() {
	outboxClient := outbox.New()

	msg := outbox.NewMessage("test", "test")
	outboxClient.PostMessage(msg)
}
