package user

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/app"
	"github.com/banyejiu/ruoyi-go/internal/middleware"
)

type Module struct {
	handler *Handler
}

func NewModule() *Module {
	repo := NewRepository(app.Global.DB)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{handler: handler}
}

// ⭐ 所有路由在这里
func (m *Module) Register(rg *gin.RouterGroup) {

	// ========================
	// 1️⃣ 公共接口（不需要权限）
	// ========================
	public := rg.Group("/public/user")
	{
		public.POST("/login", m.handler.Login)
		public.POST("/register", m.handler.Register)
	}

	// ========================
	// 2️⃣ 需要登录
	// ========================
	private := rg.Group("/user")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/profile", m.handler.Profile)
	}

	// ========================
	// 3️⃣ 权限控制
	// ========================
	admin := private.Group("")
	admin.Use(middleware.Permission("user:manage"))
	{
		admin.DELETE("/:id", m.handler.Delete)
	}
}
