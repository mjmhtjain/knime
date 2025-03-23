package repository

import (
	"fmt"

	"github.com/mjmhtjain/knime/src/internal/obj"
)

// Example shows how to use the repository layer
func Example() {
	// Create database configuration
	dbConfig := NewDBConfig(
		"localhost", // host
		"5432",      // port
		"postgres",  // user
		"password",  // password
		"knime_db",  // dbname
		"disable",   // sslmode for local development
	)

	// Connect to database
	db, err := Connect(dbConfig)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	// Migrate database schema
	if err := Migrate(db); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		return
	}

	// Create repository
	repo := NewOutboxMessageRepository(db)

	// Create a message
	message := obj.NewMessage("Hello", "World")
	if err := repo.Create(message); err != nil {
		fmt.Printf("Failed to create message: %v\n", err)
		return
	}

	fmt.Println("Message created successfully")
}
