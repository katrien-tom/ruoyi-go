package menu

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
	private := rg.Group("/system/menu")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/list", middleware.Permission("system:menu:query"), m.handler.List)
		private.GET("/treeselect", m.handler.TreeSelect)
		private.GET("/:menuId", middleware.Permission("system:menu:query"), m.handler.GetInfo)
		private.GET("/roleMenuTreeselect/:roleId", m.handler.RoleMenuTreeSelect)
	}

	admin := private.Group("")
	admin.Use(middleware.Permission("system:menu:manage"))
	{
		admin.POST("/", m.handler.Add)
		admin.PUT("/", m.handler.Edit)
		admin.DELETE("/:menuId", m.handler.Delete)
	}
}
