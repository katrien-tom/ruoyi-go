package router

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/modules/user"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.RequestID(), // ⭐ 放最前
		middleware.Logger(),    // ⭐ 记录日志
		middleware.Recovery(),  // ⭐ 捕获panic
	)

	api := r.Group("/api")

	// 注册模块
	modules := []Module{
		user.NewModule(),
	}

	for _, m := range modules {
		m.Register(api)
	}

	return r
}
