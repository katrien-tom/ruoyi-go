package user

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

	public := rg.Group("/public/user")
	{
		public.POST("/register", m.handler.Register)
	}

	private := rg.Group("/user")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/profile", m.handler.Profile)
	}

	// system user management
	system := rg.Group("/system/user")
	system.Use(middleware.JWTAuth())
	{
		system.GET("/list", m.handler.List)
		system.GET("/:userId", m.handler.GetInfo)
		system.GET("/authRole/:userId", m.handler.GetAuthRole)
		system.PUT("/authRole", m.handler.SaveAuthRole)
	}

	admin := system.Group("")
	admin.Use(middleware.Permission("system:user:manage"))
	{
		admin.POST("/", m.handler.Add)
		admin.PUT("/", m.handler.Edit)
		admin.DELETE("/:userIds", m.handler.Delete)
		admin.PUT("/resetPwd", m.handler.ResetPwd)
		admin.PUT("/changeStatus", m.handler.ChangeStatus)
	}
}
