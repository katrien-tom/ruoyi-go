package job

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
)

func (h *Handler) Routes() route.Group {
	return route.Group{
		Prefix:     "/monitor/job",
		Middlewares: []gin.HandlerFunc{middleware.JWTAuth()},
		Children: []*route.Group{
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("monitor:job:query")},
				Routes: []route.Route{
					{Method: "GET", Path: "/list", Handler: h.JobList, Meta: route.Meta{Name: "定时任务列表"}},
					{Method: "GET", Path: "/:jobId", Handler: h.GetInfo, Meta: route.Meta{Name: "定时任务详情"}},
					{Method: "GET", Path: "/log/list", Handler: h.JobLogList, Meta: route.Meta{Name: "任务日志列表"}},
				},
			},
			{
				Middlewares: []gin.HandlerFunc{middleware.Permission("monitor:job:manage")},
				Routes: []route.Route{
					{Method: "POST", Path: "/", Handler: h.Add, Meta: route.Meta{Name: "新增定时任务"}},
					{Method: "PUT", Path: "/", Handler: h.Edit, Meta: route.Meta{Name: "修改定时任务"}},
					{Method: "DELETE", Path: "/:jobId", Handler: h.Delete, Meta: route.Meta{Name: "删除定时任务"}},
					{Method: "PUT", Path: "/run", Handler: h.Run, Meta: route.Meta{Name: "执行一次"}},
				},
			},
		},
	}
}
