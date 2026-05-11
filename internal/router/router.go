package router

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/app"
	"github.com/banyejiu/ruoyi-go/internal/middleware"
	"github.com/banyejiu/ruoyi-go/internal/route"
	"github.com/banyejiu/ruoyi-go/internal/modules/auth"
	"github.com/banyejiu/ruoyi-go/internal/modules/config"
	"github.com/banyejiu/ruoyi-go/internal/modules/dept"
	"github.com/banyejiu/ruoyi-go/internal/modules/dict"
	"github.com/banyejiu/ruoyi-go/internal/modules/job"
	"github.com/banyejiu/ruoyi-go/internal/modules/loginlog"
	"github.com/banyejiu/ruoyi-go/internal/modules/menu"
	"github.com/banyejiu/ruoyi-go/internal/modules/notice"
	"github.com/banyejiu/ruoyi-go/internal/modules/operlog"
	"github.com/banyejiu/ruoyi-go/internal/modules/post"
	"github.com/banyejiu/ruoyi-go/internal/modules/role"
	"github.com/banyejiu/ruoyi-go/internal/modules/user"
	"github.com/banyejiu/ruoyi-go/pkg/validation"
)

func InitRouter() *gin.Engine {
	if err := validation.Init(); err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
	)

	api := r.Group("/api")

	db, rc := app.Global.DB, app.Global.Redis

	userSvc := user.NewService(user.NewRepository(db))
	menuSvc := menu.NewService(menu.NewRepository(db))
	roleSvc := role.NewService(role.NewRepository(db))
	deptSvc := dept.NewService(dept.NewRepository(db))
	postSvc := post.NewService(post.NewRepository(db))
	dictSvc := dict.NewService(dict.NewRepository(db))
	configSvc := config.NewService(config.NewRepository(db))
	noticeSvc := notice.NewService(notice.NewRepository(db))
	operLogSvc := operlog.NewService(operlog.NewRepository(db))
	loginLogSvc := loginlog.NewService(loginlog.NewRepository(db))
	jobSvc := job.NewService(job.NewRepository(db))

	route.RegisterModules(api,
		auth.NewHandler(auth.NewService(rc, userSvc, menuSvc)),
		user.NewHandler(userSvc),
		role.NewHandler(roleSvc),
		menu.NewHandler(menuSvc),
		dept.NewHandler(deptSvc),
		post.NewHandler(postSvc),
		dict.NewHandler(dictSvc),
		config.NewHandler(configSvc),
		notice.NewHandler(noticeSvc),
		operlog.NewHandler(operLogSvc),
		loginlog.NewHandler(loginLogSvc),
		job.NewHandler(jobSvc),
	)

	return r
}
