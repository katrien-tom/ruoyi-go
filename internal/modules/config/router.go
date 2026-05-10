package config

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
)

type Module struct {
	handler *Handler
}

func NewModule(service *Service) *Module {
	handler := NewHandler(service)

	return &Module{handler: handler}
}

func (m *Module) Register(rg *gin.RouterGroup) {
	private := rg.Group("/system/config")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/list", middleware.Permission("system:config:query"), m.handler.List)
		private.GET("/configKey/:configKey", m.handler.GetByKey)
		private.GET("/:configId", middleware.Permission("system:config:query"), m.handler.GetInfo)
	}

	admin := private.Group("")
	admin.Use(middleware.Permission("system:config:manage"))
	{
		admin.PUT("/", m.handler.Edit)
		admin.DELETE("/:configIds", m.handler.Delete)
		admin.DELETE("/refreshCache", m.handler.RefreshCache)
	}
}
