package main

import (
	"github.com/mjmhtjain/knime/src/router"
)

func main() {
	// Initialize Gin router
	router := router.NewRouter()

	// Start the server
	router.Run(":8080")
}
