package auth

import (
	"errors"

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
	var req LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, "invalid login request")
		return
	}

	data, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			response.Fail(c, response.LoginFailed, "username or password is incorrect")
		case errors.Is(err, ErrAccountDisabled):
			response.Fail(c, response.AccountDisabled, "account is disabled")
		default:
			response.Fail(c, response.InternalError, "failed to login")
		}
		return
	}

	response.SuccessNoMsg(c, data)
}

func (h *Handler) Captcha(c *gin.Context) {
	data, err := h.service.GetCaptcha(c.Request.Context())
	if err != nil {
		response.Fail(c, response.InternalError, "failed to generate captcha")
		return
	}

	response.SuccessNoMsg(c, data)
}
