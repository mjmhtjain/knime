package outbox

import (
	"errors"

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

func New() *Outbox {
	if outboxIns == nil {
		outboxIns = &Outbox{
			messageService: service.NewMessageService(),
		}
	}

	// TODO: initialize NATS for sending the messages
	// TODO: initialize DB for persisting the messages
	// TODO: initialize the internal
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
