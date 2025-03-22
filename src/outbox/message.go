package outbox

type Message struct {
	Subject string
	Body    string
}

func NewMessage(subject string, body string) *Message {
	return &Message{
		Subject: subject,
		Body:    body,
	}
}
