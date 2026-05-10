package user

import (
	"errors"
	"strconv"
	"strings"

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

func (h *Handler) List(c *gin.Context) {
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	data, err := h.service.ListUsers(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to query users")
		return
	}

	response.SuccessNoMsg(c, data)
}

func (h *Handler) GetInfo(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid userId")
		return
	}

	data, err := h.service.GetUserDetail(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.Fail(c, response.NotFound, "user not found")
			return
		}
		response.Fail(c, response.InternalError, "failed to get user")
		return
	}

	response.SuccessNoMsg(c, data)
}

func (h *Handler) Add(c *gin.Context) {
	var req AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	loginUser := getLoginUser(c)

	if err := h.service.AddUser(c.Request.Context(), req, loginUser.UserID); err != nil {
		switch {
		case errors.Is(err, ErrUserNameDuplicate):
			response.Fail(c, response.DataAlreadyExist, "user name already exists")
		case errors.Is(err, ErrPhoneDuplicate):
			response.Fail(c, response.DataAlreadyExist, "phone number already exists")
		case errors.Is(err, ErrEmailDuplicate):
			response.Fail(c, response.DataAlreadyExist, "email already exists")
		default:
			response.Fail(c, response.InternalError, err.Error())
		}
		return
	}

	response.Success[any](c, nil, "add user success")
}

func (h *Handler) Edit(c *gin.Context) {
	var req EditUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	loginUser := getLoginUser(c)

	if err := h.service.EditUser(c.Request.Context(), req, loginUser.UserID); err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			response.Fail(c, response.NotFound, "user not found")
		case errors.Is(err, ErrSelfOperation):
			response.Fail(c, response.DataInvalid, "cannot edit yourself")
		case errors.Is(err, ErrPhoneDuplicate):
			response.Fail(c, response.DataAlreadyExist, "phone number already exists")
		case errors.Is(err, ErrEmailDuplicate):
			response.Fail(c, response.DataAlreadyExist, "email already exists")
		default:
			response.Fail(c, response.InternalError, err.Error())
		}
		return
	}

	response.Success[any](c, nil, "edit user success")
}

func (h *Handler) Delete(c *gin.Context) {
	ids := strings.Split(c.Param("userIds"), ",")
	if len(ids) == 0 || (len(ids) == 1 && ids[0] == "") {
		response.Fail(c, response.ParamError, "userIds cannot be empty")
		return
	}

	userIDs := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			response.Fail(c, response.ParamError, "invalid user id: "+s)
			return
		}
		userIDs = append(userIDs, id)
	}

	loginUser := getLoginUser(c)

	if err := h.service.DeleteUsers(c.Request.Context(), userIDs, loginUser.UserID); err != nil {
		if errors.Is(err, ErrSelfOperation) {
			response.Fail(c, response.DataInvalid, "cannot delete yourself")
			return
		}
		response.Fail(c, response.InternalError, "failed to delete users")
		return
	}

	response.Success[any](c, nil, "delete users success")
}

func (h *Handler) ResetPwd(c *gin.Context) {
	var req ResetPwdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req); err != nil {
		if errors.Is(err, ErrPasswordEmpty) {
			response.Fail(c, response.DataInvalid, "password cannot be empty")
			return
		}
		response.Fail(c, response.InternalError, "failed to reset password")
		return
	}

	response.Success[any](c, nil, "reset password success")
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.ChangeStatus(c.Request.Context(), req); err != nil {
		response.Fail(c, response.InternalError, "failed to change status")
		return
	}

	response.Success[any](c, nil, "change status success")
}

func (h *Handler) GetAuthRole(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid userId")
		return
	}

	data, err := h.service.GetAuthRole(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get roles")
		return
	}

	response.SuccessNoMsg(c, data)
}

func (h *Handler) SaveAuthRole(c *gin.Context) {
	var req AuthRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.SaveAuthRole(c.Request.Context(), req); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.Fail(c, response.NotFound, "user not found")
			return
		}
		response.Fail(c, response.InternalError, "failed to save role assignment")
		return
	}

	response.Success[any](c, nil, "save role assignment success")
}

func (h *Handler) Profile(c *gin.Context) {
	loginUser := getLoginUser(c)
	data, err := h.service.GetUserDetail(c.Request.Context(), loginUser.UserID)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get profile")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) Register(c *gin.Context) {
	response.SuccessNoMsg(c, "success register")
}

func getLoginUser(c *gin.Context) *security.LoginUser {
	if lu, exists := c.Get("login_user"); exists {
		return lu.(*security.LoginUser)
	}
	return &security.LoginUser{}
}
