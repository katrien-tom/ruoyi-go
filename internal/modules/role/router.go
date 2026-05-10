package role

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
	private := rg.Group("/system/role")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/list", m.handler.List)
		private.GET("/:roleId", m.handler.GetInfo)
		private.GET("/deptTree/:roleId", m.handler.DeptTree)
	}

	admin := private.Group("")
	admin.Use(middleware.Permission("system:role:manage"))
	{
		admin.POST("/", m.handler.Add)
		admin.PUT("/", m.handler.Edit)
		admin.DELETE("/:roleIds", m.handler.Delete)
		admin.PUT("/changeStatus", m.handler.ChangeStatus)
		admin.PUT("/dataScope", m.handler.DataScope)
	}

	authUser := rg.Group("/system/role/authUser")
	authUser.Use(middleware.JWTAuth())
	{
		authUser.GET("/allocatedList", m.handler.AllocatedList)
		authUser.GET("/unallocatedList", m.handler.UnallocatedList)
	}

	authUserAdmin := authUser.Group("")
	authUserAdmin.Use(middleware.Permission("system:role:manage"))
	{
		authUserAdmin.PUT("/cancel", m.handler.CancelAuthUser)
		authUserAdmin.PUT("/cancelAll", m.handler.CancelAuthAll)
		authUserAdmin.PUT("/selectAll", m.handler.SelectAuthAll)
	}
}
