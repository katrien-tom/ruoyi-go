package auth

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

func (h *Handler) Login(c *gin.Context) {
	response.SuccessNoMsg(c, "success login")
}

func (h *Handler) Captcha(c *gin.Context) {
	data, err := h.service.GetCaptcha(c.Request.Context())
	if err != nil {
		response.Fail(c, response.InternalError, "failed to generate captcha")
		return
	}

	response.SuccessNoMsg(c, data)
}
