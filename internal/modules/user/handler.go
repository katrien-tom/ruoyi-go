package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "login"})
}

func (h *Handler) Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "register"})
}

func (h *Handler) Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "profile"})
}

func (h *Handler) Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "delete"})
}
