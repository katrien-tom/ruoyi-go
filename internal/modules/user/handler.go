package user

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(c *gin.Context) {
	response.SuccessNoMsg(c, "success register")
}

func (h *Handler) Profile(c *gin.Context) {
	response.SuccessNoMsg(c, "success profile")
}

func (h *Handler) Delete(c *gin.Context) {
	response.SuccessNoMsg(c, "success delete")
}
