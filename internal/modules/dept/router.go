package dept

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Prefix:     "/system/dept",
		Middlewares: []gin.HandlerFunc{middleware.JWTAuth()},
		Routes: []route.Route{
			{Method: "GET", Path: "/list", Handler: h.List, Meta: route.Meta{Name: "部门列表"}},
			{Method: "GET", Path: "/treeselect", Handler: h.TreeSelect, Meta: route.Meta{Name: "部门树"}},
			{Method: "GET", Path: "/:deptId", Handler: h.GetInfo, Meta: route.Meta{Name: "部门详情"}},
		},
		Children: []*route.Group{{
			Middlewares: []gin.HandlerFunc{middleware.Permission("system:dept:manage")},
			Routes: []route.Route{
				{Method: "POST", Path: "/", Handler: h.Add, Meta: route.Meta{Name: "新增部门"}},
				{Method: "PUT", Path: "/", Handler: h.Edit, Meta: route.Meta{Name: "修改部门"}},
				{Method: "DELETE", Path: "/:deptId", Handler: h.Delete, Meta: route.Meta{Name: "删除部门"}},
			},
		}},
	}
}
