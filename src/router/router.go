package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mjmhtjain/knime/internals/handlers"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/hello", handlers.HelloWorld)
	return router
}
