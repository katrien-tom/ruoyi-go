package dict

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Children: []*route.Group{
			{
				Prefix:     "/system/dict/type",
				Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
				Children: []*route.Group{
					{
						Middlewares: []gin.HandlerFunc{middleware.Permission("system:dict:query")},
						Routes: []route.Route{
							{Method: "GET", Path: "/list", Handler: h.TypeList, Meta: route.Meta{Name: "字典类型列表"}},
							{Method: "GET", Path: "/:dictId", Handler: h.TypeGetInfo, Meta: route.Meta{Name: "字典类型详情"}},
						},
					},
					{
						Middlewares: []gin.HandlerFunc{middleware.Permission("system:dict:manage")},
						Routes: []route.Route{
							{Method: "POST", Path: "/", Handler: h.TypeAdd, Meta: route.Meta{Name: "新增字典类型"}},
							{Method: "PUT", Path: "/", Handler: h.TypeEdit, Meta: route.Meta{Name: "修改字典类型"}},
							{Method: "DELETE", Path: "/:dictIds", Handler: h.TypeDelete, Meta: route.Meta{Name: "删除字典类型"}},
							{Method: "DELETE", Path: "/refreshCache", Handler: h.TypeRefreshCache, Meta: route.Meta{Name: "刷新缓存"}},
						},
					},
				},
			},
			{
				Prefix:     "/system/dict/data",
				Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
				Routes: []route.Route{
					{Method: "GET", Path: "/type/:dictType", Handler: h.DataByType, Meta: route.Meta{Name: "字典数据(类型)"}},
				},
				Children: []*route.Group{
					{
						Middlewares: []gin.HandlerFunc{middleware.Permission("system:dict:query")},
						Routes: []route.Route{
							{Method: "GET", Path: "/list", Handler: h.DataList, Meta: route.Meta{Name: "字典数据列表"}},
							{Method: "GET", Path: "/:dictCode", Handler: h.DataGetInfo, Meta: route.Meta{Name: "字典数据详情"}},
						},
					},
					{
						Middlewares: []gin.HandlerFunc{middleware.Permission("system:dict:manage")},
						Routes: []route.Route{
							{Method: "POST", Path: "/", Handler: h.DataAdd, Meta: route.Meta{Name: "新增字典数据"}},
							{Method: "PUT", Path: "/", Handler: h.DataEdit, Meta: route.Meta{Name: "修改字典数据"}},
							{Method: "DELETE", Path: "/:dictCodes", Handler: h.DataDelete, Meta: route.Meta{Name: "删除字典数据"}},
						},
					},
				},
			},
		},
	}
}
