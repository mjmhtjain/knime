package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// OutboxMessageEntity represents a message in the database
type OutboxMessageEntity struct {
	ID        string          `gorm:"primarykey;column:id"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index;column:deleted_at"`
	Subject   string          `gorm:"column:subject"`
	Body      json.RawMessage `gorm:"type:jsonb;column:body"`
	Status    string          `gorm:"column:status"`
}

// TableName specifies the table name for the OutboxMessageEntity
func (OutboxMessageEntity) TableName() string {
	return "outbox_messages"
}

// Migrate creates or updates the database schema
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&OutboxMessageEntity{})
}
