package database

import (
	"context"
	"fmt"
	"time"

	"github.com/banyejiu/ruoyi-go/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool-size"`
	MinIdleConns int           `mapstructure:"min-idle-conns"`
	DialTimeout  time.Duration `mapstructure:"dial-timeout"`
	ReadTimeout  time.Duration `mapstructure:"read-timeout"`
	WriteTimeout time.Duration `mapstructure:"write-timeout"`
}

func InitRedis() (*redis.Client, error) {
	var cfg RedisConfig
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		logger.Log.Error("Redis配置读取失败", "err", err)
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Log.Error("Redis连接失败", "err", err)
		return nil, err
	}

	logger.Log.Info("✅ Redis连接成功")
	return client, nil
}
