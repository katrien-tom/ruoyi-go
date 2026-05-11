package role

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Children: []*route.Group{
			{
				Prefix:     "/system/role",
				Middlewares: []gin.HandlerFunc{middleware.JWTAuth()},
				Routes: []route.Route{
					{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "角色列表", Permission: "system:role:query"}},
					{Method: "GET", Path: "/:roleId", Handler: h.GetInfo, Meta: route.Meta{Name: "角色详情", Permission: "system:role:query"}},
					{Method: "GET", Path: "/deptTree/:roleId", Handler: h.DeptTree, Meta: route.Meta{Name: "部门树"}},
				},
				Children: []*route.Group{{
					Middlewares: []gin.HandlerFunc{middleware.Permission("system:role:manage")},
					Routes: []route.Route{
						{Method: "POST", Path: "/", Handler: h.Add, Meta: route.Meta{Name: "新增角色"}},
						{Method: "PUT", Path: "/", Handler: h.Edit, Meta: route.Meta{Name: "修改角色"}},
						{Method: "DELETE", Path: "/:roleIds", Handler: h.Delete, Meta: route.Meta{Name: "删除角色"}},
						{Method: "PUT", Path: "/changeStatus", Handler: h.ChangeStatus, Meta: route.Meta{Name: "状态修改"}},
						{Method: "PUT", Path: "/dataScope", Handler: h.DataScope, Meta: route.Meta{Name: "数据权限"}},
					},
				}},
			},
			{
				Prefix:     "/system/role/authUser",
				Middlewares: []gin.HandlerFunc{middleware.JWTAuth()},
				Children: []*route.Group{{
					Middlewares: []gin.HandlerFunc{middleware.Permission("system:role:manage")},
					Routes: []route.Route{
						{Method: "GET", Path: "/allocatedList", Handler: h.AllocatedList, Meta: route.Meta{Name: "已分配用户"}},
						{Method: "GET", Path: "/unallocatedList", Handler: h.UnallocatedList, Meta: route.Meta{Name: "未分配用户"}},
						{Method: "PUT", Path: "/cancel", Handler: h.CancelAuthUser, Meta: route.Meta{Name: "取消授权"}},
						{Method: "PUT", Path: "/cancelAll", Handler: h.CancelAuthAll, Meta: route.Meta{Name: "批量取消授权"}},
						{Method: "PUT", Path: "/selectAll", Handler: h.SelectAuthAll, Meta: route.Meta{Name: "批量授权"}},
					},
				}},
			},
		},
	}
}
