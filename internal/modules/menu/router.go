package menu

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Prefix:     "/system/menu",
		Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
		Routes: []route.Route{
			{Method: "GET", Path: "/treeselect", Handler: h.TreeSelect, Meta: route.Meta{Name: "菜单树"}},
			{Method: "GET", Path: "/roleMenuTreeselect/:roleId", Handler: h.RoleMenuTreeSelect, Meta: route.Meta{Name: "角色菜单树"}},
		},
		Children: []*route.Group{
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("system:menu:query")},
				Routes: []route.Route{
					{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "菜单列表"}},
					{Method: "GET", Path: "/:menuId", Handler: h.GetInfo, Meta: route.Meta{Name: "菜单详情"}},
				},
			},
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("system:menu:manage")},
				Routes: []route.Route{
					{Method: "POST", Path: "/", Handler: h.Add, Meta: route.Meta{Name: "新增菜单"}},
					{Method: "PUT", Path: "/", Handler: h.Edit, Meta: route.Meta{Name: "修改菜单"}},
					{Method: "DELETE", Path: "/:menuId", Handler: h.Delete, Meta: route.Meta{Name: "删除菜单"}},
				},
			},
		},
	}
}
