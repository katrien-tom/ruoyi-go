package config

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Prefix:     "/system/config",
		Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
		Routes: []route.Route{
			{Method: "GET", Path: "/configKey/:configKey", Handler: h.GetByKey, Meta: route.Meta{Name: "参数键值"}},
		},
		Children: []*route.Group{
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("system:config:query")},
				Routes: []route.Route{
					{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "参数列表"}},
					{Method: "GET", Path: "/:configId", Handler: h.GetInfo, Meta: route.Meta{Name: "参数详情"}},
				},
			},
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("system:config:manage")},
				Routes: []route.Route{
					{Method: "PUT", Path: "/", Handler: h.Edit, Meta: route.Meta{Name: "修改参数"}},
					{Method: "DELETE", Path: "/:configIds", Handler: h.Delete, Meta: route.Meta{Name: "删除参数"}},
					{Method: "DELETE", Path: "/refreshCache", Handler: h.RefreshCache, Meta: route.Meta{Name: "刷新缓存"}},
				},
			},
		},
	}
}
