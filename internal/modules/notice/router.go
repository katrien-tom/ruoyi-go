package notice

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Prefix:     "/system/notice",
		Middlewares: []gin.HandlerFunc{middleware.JWTAuth(h.sessionService)},
		Children: []*route.Group{
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("system:notice:query")},
				Routes: []route.Route{
					{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "通知公告列表"}},
					{Method: "GET", Path: "/:noticeId", Handler: h.GetInfo, Meta: route.Meta{Name: "通知公告详情"}},
				},
			},
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("system:notice:manage")},
				Routes: []route.Route{
					{Method: "POST", Path: "/", Handler: h.Add, Meta: route.Meta{Name: "新增通知公告"}},
					{Method: "PUT", Path: "/", Handler: h.Edit, Meta: route.Meta{Name: "修改通知公告"}},
					{Method: "DELETE", Path: "/:noticeIds", Handler: h.Delete, Meta: route.Meta{Name: "删除通知公告"}},
				},
			},
		},
	}
}
