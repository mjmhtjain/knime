package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/internal/model"
	"github.com/mjmhtjain/knime/src/internal/obj"
	"github.com/mjmhtjain/knime/src/internal/repository"
	"github.com/mjmhtjain/knime/src/internal/util"
	"github.com/sirupsen/logrus"
)

type IMessageService interface {
	SaveMessage(msg *obj.Message) error
}

type MessageService struct {
	repo repository.IOutboxMessageRepository
}

func NewMessageService(outboxDBConfig *config.OutboxDBConfig) *MessageService {
	repo := repository.NewOutboxMessageRepository(outboxDBConfig)
	return &MessageService{repo: repo}
}

// SaveMessage saves the message to the database
func (s *MessageService) SaveMessage(msg *obj.Message) error {

	if msg.Subject == "" {
		err := errors.New("message subject is empty")
		logrus.WithFields(logrus.Fields{
			"Service": "MessageService",
			"Method":  "SaveMessage",
		}).Error(err)

		return err
	}

	if msg.Body == nil {
		err := errors.New("message body is nil")
		logrus.WithFields(logrus.Fields{
			"Service": "MessageService",
			"Method":  "SaveMessage",
		}).Error(err)

		return err
	}

	// Assuming message.Body is interface{}
	jsonBytes, err := json.Marshal(msg.Body)
	if err != nil {
		return err
	}

	messageEntity := &model.OutboxMessageEntity{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Subject:   msg.Subject,
		Body:      json.RawMessage(jsonBytes),
		Status:    util.MessageStatusPending,
	}

	err = s.repo.Create(messageEntity)
	if err != nil {
		return err
	}

	return nil
}
