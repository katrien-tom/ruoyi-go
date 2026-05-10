package dict

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

	dictType := rg.Group("/system/dict/type")
	dictType.Use(middleware.JWTAuth())
	{
		dictType.GET("/list", middleware.Permission("system:dict:query"), m.handler.TypeList)
		dictType.GET("/:dictId", middleware.Permission("system:dict:query"), m.handler.TypeGetInfo)
	}

	dictTypeAdmin := dictType.Group("")
	dictTypeAdmin.Use(middleware.Permission("system:dict:manage"))
	{
		dictTypeAdmin.POST("/", m.handler.TypeAdd)
		dictTypeAdmin.PUT("/", m.handler.TypeEdit)
		dictTypeAdmin.DELETE("/:dictIds", m.handler.TypeDelete)
		dictTypeAdmin.DELETE("/refreshCache", m.handler.TypeRefreshCache)
	}

	dictData := rg.Group("/system/dict/data")
	dictData.Use(middleware.JWTAuth())
	{
		dictData.GET("/list", middleware.Permission("system:dict:query"), m.handler.DataList)
		dictData.GET("/type/:dictType", m.handler.DataByType)
		dictData.GET("/:dictCode", middleware.Permission("system:dict:query"), m.handler.DataGetInfo)
	}

	dictDataAdmin := dictData.Group("")
	dictDataAdmin.Use(middleware.Permission("system:dict:manage"))
	{
		dictDataAdmin.POST("/", m.handler.DataAdd)
		dictDataAdmin.PUT("/", m.handler.DataEdit)
		dictDataAdmin.DELETE("/:dictCodes", m.handler.DataDelete)
	}
}
