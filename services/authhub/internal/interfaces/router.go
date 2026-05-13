package http

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	return gin.Default()
}
