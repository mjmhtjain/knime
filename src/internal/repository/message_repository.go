package repository

import (
	"errors"

	"github.com/mjmhtjain/knime/src/config"
	"github.com/mjmhtjain/knime/src/internal/client"
	"github.com/mjmhtjain/knime/src/internal/model"
	"github.com/mjmhtjain/knime/src/internal/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IOutboxMessageRepository defines the interface for message storage operations
type IOutboxMessageRepository interface {
	Create(messageEntity *model.OutboxMessageEntity) error
	ReadLatestPendingMessages() ([]model.OutboxMessageEntity, error)
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

// ReadLatestMessage reads the latest message from the database with status pending
func (r *OutboxMessageRepository) ReadLatestPendingMessages() ([]model.OutboxMessageEntity, error) {
	var messages []model.OutboxMessageEntity

	tx := r.db.Begin()
	// get the latest message with status pending
	result := tx.Order("created_at ASC").
		Where("status = ?", util.MessageStatusPending).
		Limit(10).
		Find(&messages)
	if result.Error != nil {
		logrus.
			WithFields(logrus.Fields{
				"Repository": "OutboxMessageRepository",
				"Method":     "ReadLatestPendingMessages",
				"Error":      result.Error,
			}).Errorf("Failed to read latest pending messages: %v", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		err := errors.New("no messages found")
		logrus.
			WithFields(logrus.Fields{
				"Repository": "OutboxMessageRepository",
				"Method":     "ReadLatestPendingMessages",
				"Error":      err,
			}).Error(err)

		return nil, err
	}

	// update all the messages status to sent
	tx.Model(&messages).Update("status", util.MessageStatusSent)
	if tx.Error != nil {
		logrus.
			WithFields(logrus.Fields{
				"Repository": "OutboxMessageRepository",
				"Method":     "ReadLatestPendingMessages",
				"Error":      tx.Error,
			}).Errorf("Failed to update message status to sent: %v", tx.Error)
		return nil, tx.Error
	}

	tx.Commit()

	return messages, nil
}
