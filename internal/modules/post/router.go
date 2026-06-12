package post

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Prefix:     "/system/post",
		Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
		Routes: []route.Route{
			{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "岗位列表"}},
			{Method: "GET", Path: "/optionselect", Handler: h.OptionSelect, Meta: route.Meta{Name: "岗位选项"}},
			{Method: "GET", Path: "/:postId", Handler: h.GetInfo, Meta: route.Meta{Name: "岗位详情"}},
		},
		Children: []*route.Group{{
			Middlewares: []gin.HandlerFunc{middleware.Permission("system:post:manage")},
			Routes: []route.Route{
				{Method: "POST", Path: "/", Handler: h.Add, Meta: route.Meta{Name: "新增岗位"}},
				{Method: "PUT", Path: "/", Handler: h.Edit, Meta: route.Meta{Name: "修改岗位"}},
				{Method: "DELETE", Path: "/:postIds", Handler: h.Delete, Meta: route.Meta{Name: "删除岗位"}},
			},
		}},
	}
}
