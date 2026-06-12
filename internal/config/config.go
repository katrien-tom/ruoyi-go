package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql" validate:"required"`
	Redis RedisConfig `mapstructure:"redis" validate:"required"`
}

type MySQLConfig struct {
	Host            string        `mapstructure:"host" validate:"required,hostname|ip"`
	Port            int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username        string        `mapstructure:"username" validate:"required"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database" validate:"required"`
	Charset         string        `mapstructure:"charset" validate:"required"`
	MaxIdleConns    int           `mapstructure:"max-idle-conns" validate:"min=1"`
	MaxOpenConns    int           `mapstructure:"max-open-conns" validate:"min=1"`
	ConnMaxLifetime time.Duration `mapstructure:"conn-max-lifetime"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host" validate:"required,hostname|ip"`
	Port         int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db" validate:"min=0,max=15"`
	PoolSize     int           `mapstructure:"pool-size" validate:"min=1"`
	MinIdleConns int           `mapstructure:"min-idle-conns" validate:"min=0"`
	DialTimeout  time.Duration `mapstructure:"dial-timeout" validate:"required"`
	ReadTimeout  time.Duration `mapstructure:"read-timeout" validate:"required"`
	WriteTimeout time.Duration `mapstructure:"write-timeout" validate:"required"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}
