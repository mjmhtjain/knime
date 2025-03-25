package outbox

import (
	"context"
	"errors"
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
	outboxDBConfig *config.OutboxDBConfig
	natsConfig     *config.NatsConfig
	messageService service.IMessageService
	outboxService  service.IOutboxService
}

func New(
	outboxDBConfig *config.OutboxDBConfig,
	natsConfig *config.NatsConfig,
) *Outbox {
	if outboxIns == nil {
		outboxIns = &Outbox{
			outboxDBConfig: outboxDBConfig,
			natsConfig:     natsConfig,
			messageService: service.NewMessageService(outboxDBConfig, natsConfig),
			outboxService:  service.NewOutboxService(outboxDBConfig, natsConfig),
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
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go o.outboxService.ConsumeOutboxMessages()
		}
	}
}
