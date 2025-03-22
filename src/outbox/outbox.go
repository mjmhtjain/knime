package outbox

import (
	"github.com/mjmhtjain/knime/src/internal/obj"
	"github.com/mjmhtjain/knime/src/internal/service"
)

var outboxIns *Outbox

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
	msg := obj.NewMessage(message.Subject, message.Body)
	return o.messageService.SaveMessage(msg)
}
