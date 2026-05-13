package response

import (
	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
