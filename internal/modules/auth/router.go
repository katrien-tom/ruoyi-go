package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Routes: []route.Route{
			{Method: "GET", Path: "/public/auth/captchaImage", Handler: h.Captcha, Meta: route.Meta{Name: "验证码"}},
			{Method: "POST", Path: "/public/auth/login", Handler: h.Login, Meta: route.Meta{Name: "登录"}},
		},
		Children: []*route.Group{{
			Prefix:      "/auth",
			Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
			Routes: []route.Route{
				{Method: "GET", Path: "/getInfo", Handler: h.GetInfo, Meta: route.Meta{Name: "用户信息"}},
				{Method: "GET", Path: "/getRouters", Handler: h.GetRouters, Meta: route.Meta{Name: "路由信息"}},
				{Method: "POST", Path: "/logout", Handler: h.Logout, Meta: route.Meta{Name: "登出"}},
				{Method: "GET", Path: "/online", Handler: h.OnlineUsers, Meta: route.Meta{Name: "在线用户", Permission: "monitor:online:query"}},
				{Method: "DELETE", Path: "/online/:token", Handler: h.ForceLogout, Meta: route.Meta{Name: "强退用户", Permission: "monitor:online:forceLogout"}},
			},
		}},
	}
}
