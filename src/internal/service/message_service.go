package service

import "github.com/mjmhtjain/knime/src/internal/obj"

type IMessageService interface {
	SaveMessage(msg *obj.Message) error
}

type MessageService struct {
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

// SaveMessage saves the message to the database
func (s *MessageService) SaveMessage(msg *obj.Message) error {
	return nil
}
