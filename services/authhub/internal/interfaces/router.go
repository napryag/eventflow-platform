package http

import (
	"github.com/gin-gonic/gin"
	"github.com/napryag/eventflow-platform/services/authhub/config"
)

type GinRouter struct {
	router *gin.Engine
}

func SetupRouter() GinRouter {
	router := gin.Default()
	router.GET("/health", HealthHandler)
	return GinRouter{router: router}
}

func (g *GinRouter) Run(address config.HTTPConfig) error {
	if err := g.router.Run(address.Host + ":" + address.Port); err != nil {
		return err
	}
	return nil
}
