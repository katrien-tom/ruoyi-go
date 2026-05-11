package operlog

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Prefix:     "/monitor/operlog",
		Middlewares: []gin.HandlerFunc{middleware.JWTAuth()},
		Children: []*route.Group{
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("monitor:operlog:query")},
				Routes: []route.Route{
					{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "操作日志列表"}},
				},
			},
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("monitor:operlog:manage")},
				Routes: []route.Route{
					{Method: "DELETE", Path: "/:operIds", Handler: h.Delete, Meta: route.Meta{Name: "删除操作日志"}},
					{Method: "DELETE", Path: "/clean", Handler: h.Clean, Meta: route.Meta{Name: "清空操作日志"}},
				},
			},
		},
	}
}
