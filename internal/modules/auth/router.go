package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
)

type Module struct {
	handler *Handler
}

func NewModule(redisClient *redis.Client, userAuthReader UserAuthReader, menuReader MenuReader) *Module {
	service := NewService(redisClient, userAuthReader, menuReader)
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
		private.GET("/getRouters", m.handler.GetRouters)
		private.POST("/logout", m.handler.Logout)
		private.GET("/online", middleware.Permission("monitor:online:query"), m.handler.OnlineUsers)
		private.DELETE("/online/:token", middleware.Permission("monitor:online:forceLogout"), m.handler.ForceLogout)
	}
}
