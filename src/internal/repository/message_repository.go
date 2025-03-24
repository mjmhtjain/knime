package repository

import (
	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/internal/client"
	"github.com/mjmhtjain/knime/src/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IOutboxMessageRepository defines the interface for message storage operations
type IOutboxMessageRepository interface {
	Create(messageEntity *model.OutboxMessageEntity) error
}

// OutboxMessageRepository implements MessageRepository using GORM
type OutboxMessageRepository struct {
	db *gorm.DB
}

// NewOutboxMessageRepository creates a new message repository
func NewOutboxMessageRepository(outboxDBConfig *config.OutboxDBConfig) IOutboxMessageRepository {
	db, err := client.NewDBClient(outboxDBConfig)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	if err := model.Migrate(db); err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
	}

	return &OutboxMessageRepository{db: db}
}

// Create stores a new message in the database
func (r *OutboxMessageRepository) Create(messageEntity *model.OutboxMessageEntity) error {
	result := r.db.Create(messageEntity)
	if result.Error != nil {
		logrus.
			WithFields(logrus.Fields{
				"Repository": "OutboxMessageRepository",
				"Method":     "Create",
				"Error":      result.Error,
			}).Errorf("Failed to create message: %v", result.Error)
		return result.Error
	}

	return nil
}
