package post

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
	private := rg.Group("/system/post")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/list", m.handler.List)
		private.GET("/optionselect", m.handler.OptionSelect)
		private.GET("/:postId", m.handler.GetInfo)
	}

	admin := private.Group("")
	admin.Use(middleware.Permission("system:post:manage"))
	{
		admin.POST("/", m.handler.Add)
		admin.PUT("/", m.handler.Edit)
		admin.DELETE("/:postIds", m.handler.Delete)
	}
}
