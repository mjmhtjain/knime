package outbox

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/internal/obj"
	"github.com/mjmhtjain/knime/src/internal/service"
	"github.com/sirupsen/logrus"
)

var outboxIns *Outbox

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

type Outbox struct {
	pushMessageInterval int
	outboxDBConfig      *config.OutboxDBConfig
	natsConfig          *config.NatsConfig
	messageService      service.IMessageService
	outboxService       service.IOutboxService
}

func New(
	outboxDBConfig *config.OutboxDBConfig,
	natsConfig *config.NatsConfig,
) *Outbox {
	if outboxIns == nil {
		pushMessageInterval := os.Getenv("OUTBOX_PUSH_MESSAGE_INTERVAL")
		if pushMessageInterval == "" {
			pushMessageInterval = "5" // default to 5 seconds
		}

		interval, err := strconv.Atoi(pushMessageInterval)
		if err != nil {
			logrus.Fatalf("Error converting OUTBOX_PUSH_MESSAGE_INTERVAL to int: %v", err)
		}

		outboxIns = &Outbox{
			pushMessageInterval: interval,
			outboxDBConfig:      outboxDBConfig,
			natsConfig:          natsConfig,
			messageService:      service.NewMessageService(outboxDBConfig, natsConfig),
			outboxService:       service.NewOutboxService(outboxDBConfig, natsConfig),
		}
	}

	return outboxIns
}

func (o *Outbox) PostMessage(message *Message) error {
	if message == nil {
		err := errors.New("message is nil")
		return err
	}

	msg := obj.NewMessage(message.Subject, message.Body)
	return o.messageService.SaveMessage(msg)
}

func (o *Outbox) LaunchOutboxService(ctx context.Context) {

	ticker := time.NewTicker(time.Duration(o.pushMessageInterval) * time.Second) // Push messages every pushMessageInterval seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := o.outboxService.ConsumeOutboxMessages()
			if err != nil {
				logrus.Error("Error consuming outbox messages", err)
			}
		}
	}
}
