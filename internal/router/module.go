package router

import "github.com/gin-gonic/gin"

type Module interface {
	Register(rg *gin.RouterGroup)
}
