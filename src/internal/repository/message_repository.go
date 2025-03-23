package repository

import (
	"encoding/json"

	"github.com/mjmhtjain/knime/src/internal/obj"
	"gorm.io/gorm"
)

// IOutboxMessageRepository defines the interface for message storage operations
type IOutboxMessageRepository interface {
	Create(message *obj.Message) error
	FindByID(id uint) (*obj.Message, error)
}

// OutboxMessageRepository implements MessageRepository using GORM
type OutboxMessageRepository struct {
	db *gorm.DB
}

// NewOutboxMessageRepository creates a new message repository
func NewOutboxMessageRepository(db *gorm.DB) IOutboxMessageRepository {
	return &OutboxMessageRepository{db: db}
}

// Create stores a new message in the database
func (r *OutboxMessageRepository) Create(message *obj.Message) error {
	// Convert domain object to entity
	entity := OutboxMessageEntity{
		Subject: message.Subject,
	}

	// Assuming message.Body is interface{}
	jsonBytes, err := json.Marshal(message.Body)
	if err != nil {
		return err
	}
	entity.Body = json.RawMessage(jsonBytes)

	// Save to database
	result := r.db.Create(&entity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindByID retrieves a message by its ID
func (r *OutboxMessageRepository) FindByID(id uint) (*obj.Message, error) {
	var entity OutboxMessageEntity
	result := r.db.First(&entity, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert entity to domain object
	message := &obj.Message{
		Subject: entity.Subject,
		Body:    entity.Body,
	}

	return message, nil
}
