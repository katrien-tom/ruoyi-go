package job

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
)

type Module struct {
	handler *Handler
}

func NewModule(service *Service) *Module {
	handler := NewHandler(service)

	return &Module{handler: handler}
}

func (m *Module) Register(rg *gin.RouterGroup) {
	private := rg.Group("/job")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/list", m.handler.JobList)
		private.GET("/:jobId", m.handler.GetInfo)
		private.GET("/log/list", m.handler.JobLogList)
	}

	admin := private.Group("")
	admin.Use(middleware.Permission("monitor:job:manage"))
	{
		admin.POST("/", m.handler.Add)
		admin.PUT("/", m.handler.Edit)
		admin.DELETE("/:jobId", m.handler.Delete)
		admin.PUT("/run", m.handler.Run)
	}
}
