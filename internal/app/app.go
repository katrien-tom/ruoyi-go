package app

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	DB    *gorm.DB
	Redis *redis.Client
}

var Global *App

func Init(db *gorm.DB, rdb *redis.Client) {
	Global = &App{
		DB:    db,
		Redis: rdb,
	}
}

//test ceshi hotfix
