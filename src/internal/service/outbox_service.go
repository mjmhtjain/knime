package service

import (
	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/internal/repository"
	"github.com/sirupsen/logrus"
)

type IOutboxService interface {
	ConsumeOutboxMessages() error
}

type OutboxService struct {
	outboxRepository repository.IOutboxMessageRepository
}

func NewOutboxService(outboxDBConfig *config.OutboxDBConfig, natsConfig *config.NatsConfig) *OutboxService {
	repo := repository.NewOutboxMessageRepository(outboxDBConfig, natsConfig)

	return &OutboxService{
		outboxRepository: repo,
	}
}

func (s *OutboxService) ConsumeOutboxMessages() error {
	// read the latest message from the outbox table
	messages, err := s.outboxRepository.ReadLatestPendingMessages()
	if err != nil {
		return err
	}

	// log the messages
	for _, message := range messages {
		logrus.
			WithFields(logrus.Fields{
				"Service": "OutboxService",
				"Method":  "ConsumeOutboxMessages",
				"Message": message,
			}).Info("Consuming outbox message")
	}

	return nil
}
