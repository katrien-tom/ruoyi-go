package database

import (
	"context"
	"fmt"
	"time"

	"github.com/banyejiu/ruoyi-go/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// DBConfig 数据库配置
type DBConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	Charset         string        `mapstructure:"charset"`
	MaxIdleConns    int           `mapstructure:"max-idle-conns"`
	MaxOpenConns    int           `mapstructure:"max-open-conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn-max-lifetime"`
}

// InitDB 初始化DB，自动读取 config.yml + 接入你的 slog 日志
func InitDB() (*gorm.DB, error) {
	// 1. 读取 Viper 配置
	var cfg DBConfig
	if err := viper.UnmarshalKey("mysql", &cfg); err != nil {
		logger.Log.Error("数据库配置读取失败", "err", err)
		return nil, err
	}

	// 2. 构建DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset,
	)

	// 3. GORM 日志接入你的 slog（关键！统一日志）
	gormLogger := newGormSlogLogger()

	// 4. 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logger.Log.Error("数据库连接失败", "err", err)
		return nil, err
	}

	// 5. 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Log.Error("数据库连接失败", "err", err)
		return nil, err
	}

	logger.Log.Info("✅ 数据库连接成功")
	return db, nil
}

// newGormSlogLogger 将GORM日志输出到你的 slog
func newGormSlogLogger() gormLogger.Interface {
	return gormLogger.New(
		&gormSlogWriter{},
		gormLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormLogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

// gormSlogWriter 对接 slog 核心实现
type gormSlogWriter struct{}

func (w *gormSlogWriter) Printf(msg string, args ...interface{}) {
	logger.Log.Info(fmt.Sprintf(msg, args...))
}
