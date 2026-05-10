package loginlog

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
	private := rg.Group("/monitor/logininfor")
	private.Use(middleware.JWTAuth(), middleware.Permission("monitor:logininfor:query"))
	{
		private.GET("/list", m.handler.List)
	}

	admin := private.Group("")
	admin.Use(middleware.Permission("monitor:logininfor:manage"))
	{
		admin.DELETE("/:infoIds", m.handler.Delete)
		admin.DELETE("/clean", m.handler.Clean)
	}
}
