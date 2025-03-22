package outbox

var outboxIns *Outbox

type Outbox struct {
}

func New() *Outbox {
	if outboxIns == nil {
		outboxIns = &Outbox{}
	}

	// TODO: initialize NATS for sending the messages
	// TODO: initialize DB for persisting the messages
	return outboxIns
}

func (o *Outbox) PostMessage(message *Message) error {
	return nil
}
