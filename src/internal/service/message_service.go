package service

import (
	"errors"

	"github.com/mjmhtjain/knime/src/internal/obj"
	"github.com/sirupsen/logrus"
)

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

	if msg.Subject == "" {
		err := errors.New("message subject is empty")
		return err
	}

	if msg.Body == nil {
		err := errors.New("message body is nil")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"message.subject": msg.Subject,
	}).Info("Saving message")

	return nil
}
