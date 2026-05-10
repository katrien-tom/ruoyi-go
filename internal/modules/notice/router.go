package notice

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
	private := rg.Group("/system/notice")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/list", middleware.Permission("system:notice:query"), m.handler.List)
		private.GET("/:noticeId", m.handler.GetInfo)
	}

	admin := private.Group("")
	admin.Use(middleware.Permission("system:notice:manage"))
	{
		admin.POST("/", m.handler.Add)
		admin.PUT("/", m.handler.Edit)
		admin.DELETE("/:noticeIds", m.handler.Delete)
	}
}
