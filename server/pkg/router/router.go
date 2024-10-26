package router

import "github.com/gin-gonic/gin"

func InitializeRoutes() *gin.Engine {
	router := gin.Default()

	return router
}
