package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/banyejiu/ruoyi-go/internal/app"
	"github.com/banyejiu/ruoyi-go/internal/router"
	"github.com/banyejiu/ruoyi-go/pkg/database"
	"github.com/banyejiu/ruoyi-go/pkg/logger"
)

func main() {

	// 1️⃣ 初始化日志
	logger.Init("dev")

	// 2️⃣ 加载配置
	if err := initConfig(); err != nil {
		panic(err)
	}

	// 3️⃣ 初始化数据库
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	// 4️⃣ 初始化Redis
	rdb, err := database.InitRedis()
	if err != nil {
		panic(err)
	}

	// 5️⃣ 初始化全局依赖容器
	app.Init(db, rdb)

	// 6️⃣ 初始化路由
	r := router.InitRouter()
	r.Run(":8080")
}

func initConfig() error {
	viper.SetConfigFile("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	return nil
}
