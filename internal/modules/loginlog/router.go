package loginlog

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Prefix:     "/monitor/logininfor",
		Middlewares: []gin.HandlerFunc{middleware.JWTAuth()},
		Children: []*route.Group{
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("monitor:logininfor:query")},
				Routes: []route.Route{
					{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "登录日志列表"}},
				},
			},
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("monitor:logininfor:manage")},
				Routes: []route.Route{
					{Method: "DELETE", Path: "/:infoIds", Handler: h.Delete, Meta: route.Meta{Name: "删除登录日志"}},
					{Method: "DELETE", Path: "/clean", Handler: h.Clean, Meta: route.Meta{Name: "清空登录日志"}},
				},
			},
		},
	}
}
