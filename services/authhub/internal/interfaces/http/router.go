package http

import (
	"github.com/gin-gonic/gin"
	"github.com/napryag/eventflow-platform/pkg/errs"
	"github.com/napryag/eventflow-platform/services/authhub/config"
)

type GinRouter struct {
	engine *gin.Engine
}

func New() *GinRouter {
	engine := gin.Default()

	r := &GinRouter{
		engine: engine,
	}

	r.setupRoutes()

	return r
}

func (g *GinRouter) Run(address config.HTTPConfig) error {
	if err := g.engine.Run(address.Host + ":" + address.Port); err != nil {
		return errs.New("faield to attach the router to http.Server").Wrap(err)
	}

	return nil
}

func (g *GinRouter) setupRoutes() {
	g.engine.GET("/health", HealthHandler)
}
