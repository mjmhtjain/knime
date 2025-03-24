package outbox

import (
	"errors"

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
	messageService service.IMessageService
}

func New(
	outboxDBConfig *config.OutboxDBConfig,
	natsConfig *config.NatsConfig,
) *Outbox {
	if outboxIns == nil {
		outboxIns = &Outbox{
			messageService: service.NewMessageService(outboxDBConfig),
		}

		// start the outbox service in a new goroutine
		outboxService := service.NewOutboxService(outboxDBConfig)
		go outboxService.ConsumeOutboxMessages()
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
