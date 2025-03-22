package obj

type Message struct {
	Subject string
	Body    interface{}
}

func NewMessage(subject string, body interface{}) *Message {
	return &Message{
		Subject: subject,
		Body:    body,
	}
}
