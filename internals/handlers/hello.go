package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloWorld handles GET request and returns a hello world message
func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}
