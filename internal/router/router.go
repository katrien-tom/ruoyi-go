package router

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/modules/auth"
	"github.com/banyejiu/ruoyi-go/internal/modules/user"
	"github.com/banyejiu/ruoyi-go/pkg/validation"
)

func InitRouter() *gin.Engine {
	if err := validation.Init(); err != nil {
		panic(err)
	}

	r := gin.New()

	r.Use(
		middleware.RequestID(), // ⭐ 放最前
		middleware.Logger(),    // ⭐ 记录日志
		middleware.Recovery(),  // ⭐ 捕获panic
	)

	api := r.Group("/api")

	// 注册模块
	modules := []Module{
		auth.NewModule(),
		user.NewModule(),
	}

	for _, m := range modules {
		m.Register(api)
	}

	return r
}
