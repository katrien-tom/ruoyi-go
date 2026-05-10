package role

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

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
	var req RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.ListRoles(c.Request.Context(), req, pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get role list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) GetInfo(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid roleId")
		return
	}

	data, err := h.service.GetByID(c.Request.Context(), roleID)
	if err != nil {
		response.Fail(c, response.NotFound, "role not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) Add(c *gin.Context) {
	var req AddRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.AddRole(c.Request.Context(), req); err != nil {
		switch {
		case errors.Is(err, errRoleNameDuplicate):
			response.Fail(c, response.DataAlreadyExist, "role name already exists")
		case errors.Is(err, errRoleKeyDuplicate):
			response.Fail(c, response.DataAlreadyExist, "role key already exists")
		default:
			response.Fail(c, response.InternalError, "failed to add role")
		}
		return
	}
	response.Success[any](c, nil, "add role success")
}

func (h *Handler) Edit(c *gin.Context) {
	var req EditRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.EditRole(c.Request.Context(), req); err != nil {
		switch {
		case errors.Is(err, errRoleNameDuplicate):
			response.Fail(c, response.DataAlreadyExist, "role name already exists")
		case errors.Is(err, errRoleKeyDuplicate):
			response.Fail(c, response.DataAlreadyExist, "role key already exists")
		default:
			response.Fail(c, response.InternalError, "failed to edit role")
		}
		return
	}
	response.Success[any](c, nil, "edit role success")
}

func (h *Handler) Delete(c *gin.Context) {
	ids := strings.Split(c.Param("roleIds"), ",")
	roleIDs := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			response.Fail(c, response.ParamError, "invalid role id: "+s)
			return
		}
		roleIDs = append(roleIDs, id)
	}

	if err := h.service.DeleteRoles(c.Request.Context(), roleIDs); err != nil {
		response.Fail(c, response.InternalError, "failed to delete roles")
		return
	}
	response.Success[any](c, nil, "delete roles success")
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	var req ChangeRoleStatusRequest
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

func (h *Handler) DataScope(c *gin.Context) {
	var req DataScopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.UpdateDataScope(c.Request.Context(), req); err != nil {
		response.Fail(c, response.InternalError, "failed to update data scope")
		return
	}
	response.Success[any](c, nil, "update data scope success")
}

func (h *Handler) DeptTree(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid roleId")
		return
	}

	data, err := h.service.GetDeptTreeForRole(c.Request.Context(), roleID)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get dept tree")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) AllocatedList(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.DefaultQuery("roleId", "0"), 10, 64)
	if err != nil || roleID == 0 {
		response.Fail(c, response.ParamError, "invalid roleId")
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.GetAllocatedUsers(c.Request.Context(), roleID,
		c.Query("userName"), c.Query("phonenumber"), pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get allocated users")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) UnallocatedList(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.DefaultQuery("roleId", "0"), 10, 64)
	if err != nil || roleID == 0 {
		response.Fail(c, response.ParamError, "invalid roleId")
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.GetUnallocatedUsers(c.Request.Context(), roleID,
		c.Query("userName"), c.Query("phonenumber"), pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get unallocated users")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) CancelAuthUser(c *gin.Context) {
	var req AuthUserCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.CancelAuthUser(c.Request.Context(), req.RoleID, req.UserID); err != nil {
		response.Fail(c, response.InternalError, "failed to cancel auth user")
		return
	}
	response.Success[any](c, nil, "cancel auth user success")
}

func (h *Handler) CancelAuthAll(c *gin.Context) {
	var req AuthUserCancelAllRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.CancelAuthUsers(c.Request.Context(), req.RoleID, req.UserIDs); err != nil {
		response.Fail(c, response.InternalError, "failed to cancel auth users")
		return
	}
	response.Success[any](c, nil, "cancel auth users success")
}

func (h *Handler) SelectAuthAll(c *gin.Context) {
	var req AuthUserSelectAllRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.SelectAuthUsers(c.Request.Context(), req.RoleID, req.UserIDs); err != nil {
		response.Fail(c, response.InternalError, "failed to select auth users")
		return
	}
	response.Success[any](c, nil, "select auth users success")
}
