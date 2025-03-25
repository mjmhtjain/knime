package repository

import (
	"encoding/json"

	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/internal/client"
	"github.com/mjmhtjain/knime/src/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type INatsRepository interface {
	PublishMessage(message *model.OutboxMessageEntity) error
}

type NatsRepository struct {
	natsClient *nats.Conn
}

func NewNatsRepository(natsConfig *config.NatsConfig) INatsRepository {
	natsClient, err := client.NewNatsClient(natsConfig)
	if err != nil {
		logrus.Fatalf("Failed to create NATS client: %v", err)
	}

	return &NatsRepository{
		natsClient: natsClient,
	}
}

func (r *NatsRepository) PublishMessage(message *model.OutboxMessageEntity) error {
	jsonBody, err := json.Marshal(message.Body)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"Repository": "NatsRepository",
				"Method":     "PublishMessage",
				"Error":      err,
			}).Errorf("Failed to marshal message body: %v", err)
		return err
	}

	err = r.natsClient.Publish(message.Subject, jsonBody)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"Repository": "NatsRepository",
				"Method":     "PublishMessage",
				"Error":      err,
			}).Errorf("Failed to publish message to nats: %v", err)
		return err
	}

	return nil
}
