package config

import (
	"fmt"
)

// OutboxDBConfig holds database connection configuration
type OutboxDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewOutboxDBConfig creates a new database configuration
func NewOutboxDBConfig(host, port, user, password, dbName string) *OutboxDBConfig {
	return &OutboxDBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}
}

func (c *OutboxDBConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}
