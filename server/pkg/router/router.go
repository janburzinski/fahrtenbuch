package router

import (
	"server/pkg/handler"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes() *gin.Engine {
	router := gin.Default()
	r1 := router.Group("/api/v1")

	// auth routes
	r1.POST("/auth/register", handler.Register)

	return router
}
