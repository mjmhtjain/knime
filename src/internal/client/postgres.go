package client

import (
	"time"

	"github.com/mjmhtjain/knime/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDB *gorm.DB = nil

func NewDBClient(config *config.OutboxDBConfig) (*gorm.DB, error) {
	if postgresDB != nil {
		return postgresDB, nil
	}

	var err error
	postgresDB, err = gorm.Open(postgres.Open(config.GetConnectionString()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Test the connection
	sqlDB, err := postgresDB.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return postgresDB, nil
}
