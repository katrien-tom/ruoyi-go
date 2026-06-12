package user

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Routes: []route.Route{
			{Method: "POST", Path: "/public/user/register", Handler: h.Register},
		},
		Children: []*route.Group{
			{
				Prefix:     "/user",
				Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
				Routes: []route.Route{
					{Method: "GET", Path: "/profile", Handler: h.Profile, Meta: route.Meta{Name: "个人信息"}},
				},
			},
			{
				Prefix:     "/system/user",
				Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
				Routes: []route.Route{
					{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "用户列表", Permission: "system:user:query"}},
					{Method: "GET", Path: "/:userId", Handler: h.GetInfo, Meta: route.Meta{Name: "用户详情", Permission: "system:user:query"}},
				},
				Children: []*route.Group{{
					Middlewares: []gin.HandlerFunc{middleware.Permission("system:user:manage")},
					Routes: []route.Route{
						{Method: "POST", Path: "/", Handler: h.Add, Meta: route.Meta{Name: "新增用户"}},
						{Method: "PUT", Path: "/", Handler: h.Edit, Meta: route.Meta{Name: "修改用户"}},
						{Method: "DELETE", Path: "/:userIds", Handler: h.Delete, Meta: route.Meta{Name: "删除用户"}},
						{Method: "PUT", Path: "/resetPwd", Handler: h.ResetPwd, Meta: route.Meta{Name: "重置密码"}},
						{Method: "PUT", Path: "/changeStatus", Handler: h.ChangeStatus, Meta: route.Meta{Name: "状态修改"}},
						{Method: "GET", Path: "/authRole/:userId", Handler: h.GetAuthRole, Meta: route.Meta{Name: "分配角色"}},
						{Method: "PUT", Path: "/authRole", Handler: h.SaveAuthRole, Meta: route.Meta{Name: "保存角色"}},
					},
				}},
			},
		},
	}
}
