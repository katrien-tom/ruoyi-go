package dept

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
	private := rg.Group("/system/dept")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/list", m.handler.List)
		private.GET("/treeselect", m.handler.TreeSelect)
		private.GET("/:deptId", m.handler.GetInfo)
	}

	admin := private.Group("")
	admin.Use(middleware.Permission("system:dept:manage"))
	{
		admin.POST("/", m.handler.Add)
		admin.PUT("/", m.handler.Edit)
		admin.DELETE("/:deptId", m.handler.Delete)
	}
}
