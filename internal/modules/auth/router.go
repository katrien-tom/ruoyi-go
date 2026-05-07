package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/app"
	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/modules/user"
)

type Module struct {
	handler *Handler
}

func NewModule() *Module {
	repo := user.NewRepository(app.Global.DB)
	userService := user.NewService(repo)
	service := NewService(app.Global.Redis, userService)
	handler := NewHandler(service)

	return &Module{handler: handler}
}

func (m *Module) Register(rg *gin.RouterGroup) {
	public := rg.Group("/public/auth")
	{
		public.GET("/captchaImage", m.handler.Captcha)
		public.POST("/login", m.handler.Login)
	}

	private := rg.Group("/auth")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/getInfo", m.handler.GetInfo)
		private.POST("/logout", m.handler.Logout)
		private.GET("/online", middleware.Permission("monitor:online:query"), m.handler.OnlineUsers)
		private.DELETE("/online/:token", middleware.Permission("monitor:online:forceLogout"), m.handler.ForceLogout)
	}
}
