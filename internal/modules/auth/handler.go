package auth

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/internal/security"
	"github.com/banyejiu/ruoyi-go/pkg/response"
	"github.com/banyejiu/ruoyi-go/pkg/validation"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	meta := BuildLoginMeta(c.Request, nil)
	data, err := h.service.Login(c.Request.Context(), req, meta)
	if err != nil {
		switch {
		case errors.Is(err, ErrCaptchaInvalid):
			response.Fail(c, response.CaptchaInvalid, "captcha is incorrect or expired")
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

func (h *Handler) Logout(c *gin.Context) {
	tokenString, _ := c.Get("login_token")
	if err := h.service.Logout(c.Request.Context(), toString(tokenString)); err != nil {
		switch {
		case errors.Is(err, security.ErrTokenInvalid), errors.Is(err, security.ErrTokenNotFound):
			response.Fail(c, response.Unauthorized, "invalid token")
		default:
			response.Fail(c, response.InternalError, "failed to logout")
		}
		return
	}

	response.SuccessNoMsg[any](c, nil)
}

func (h *Handler) GetInfo(c *gin.Context) {
	tokenString, _ := c.Get("login_token")
	data, err := h.service.GetInfo(c.Request.Context(), toString(tokenString))
	if err != nil {
		switch {
		case errors.Is(err, security.ErrTokenInvalid), errors.Is(err, security.ErrTokenNotFound):
			response.Fail(c, response.Unauthorized, "invalid token")
		default:
			response.Fail(c, response.InternalError, "failed to load user info")
		}
		return
	}

	response.SuccessNoMsg(c, data)
}

func (h *Handler) OnlineUsers(c *gin.Context) {
	data, err := h.service.ListOnlineUsers(c.Request.Context())
	if err != nil {
		response.Fail(c, response.InternalError, "failed to load online users")
		return
	}

	response.SuccessNoMsg(c, data)
}

func (h *Handler) ForceLogout(c *gin.Context) {
	tokenID := c.Param("token")
	if err := h.service.ForceLogout(c.Request.Context(), tokenID); err != nil {
		switch {
		case errors.Is(err, security.ErrTokenInvalid), errors.Is(err, security.ErrTokenNotFound):
			response.Fail(c, response.NotFound, "online user not found")
		default:
			response.Fail(c, response.InternalError, "failed to force logout")
		}
		return
	}

	response.SuccessNoMsg[any](c, nil)
}

func toString(value any) string {
	s, _ := value.(string)
	return s
}
