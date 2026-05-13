package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/napryag/eventflow-platform/pkg/response"
)

func HealthHandler(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.HealthResponse{
		Status:  "ok",
		Service: "authhub",
	})
}
