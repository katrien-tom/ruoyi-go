package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/app"
	"github.com/banyejiu/ruoyi-go/internal/modules/user"
)

type Module struct {
	handler *Handler
}

func NewModule() *Module {
	repo := user.NewRepository(app.Global.DB)
	service := NewService(app.Global.Redis, repo)
	handler := NewHandler(service)

	return &Module{handler: handler}
}

func (m *Module) Register(rg *gin.RouterGroup) {
	public := rg.Group("/public/auth")
	{
		public.GET("/captchaImage", m.handler.Captcha)
		public.POST("/login", m.handler.Login)
	}
}
