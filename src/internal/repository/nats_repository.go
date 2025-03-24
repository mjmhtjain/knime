package repository

import (
	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/internal/client"
	"github.com/mjmhtjain/knime/src/internal/obj"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type INatsRepository interface {
	PublishMessage(message *obj.Message) error
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

func (r *NatsRepository) PublishMessage(message *obj.Message) error {
	return nil
}
