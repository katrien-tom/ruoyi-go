package dept

import (
	"errors"
	"strconv"

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
	var req DeptListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	data, err := h.service.GetDeptList(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get dept list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) TreeSelect(c *gin.Context) {
	data, err := h.service.GetDeptTreeSelect(c.Request.Context())
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get dept tree")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) GetInfo(c *gin.Context) {
	deptID, err := strconv.ParseInt(c.Param("deptId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid deptId")
		return
	}

	data, err := h.service.GetByID(c.Request.Context(), deptID)
	if err != nil {
		response.Fail(c, response.NotFound, "dept not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) Add(c *gin.Context) {
	var req AddDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.AddDept(c.Request.Context(), req); err != nil {
		if errors.Is(err, errDeptNameDuplicate) {
			response.Fail(c, response.DataAlreadyExist, "dept name already exists")
			return
		}
		response.Fail(c, response.InternalError, err.Error())
		return
	}
	response.Success[any](c, nil, "add dept success")
}

func (h *Handler) Edit(c *gin.Context) {
	var req EditDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.EditDept(c.Request.Context(), req); err != nil {
		if errors.Is(err, errDeptNameDuplicate) {
			response.Fail(c, response.DataAlreadyExist, "dept name already exists")
			return
		}
		response.Fail(c, response.InternalError, err.Error())
		return
	}
	response.Success[any](c, nil, "edit dept success")
}

func (h *Handler) Delete(c *gin.Context) {
	deptID, err := strconv.ParseInt(c.Param("deptId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid deptId")
		return
	}

	if err := h.service.DeleteDept(c.Request.Context(), deptID); err != nil {
		switch {
		case errors.Is(err, errDeptHasChildren):
			response.Fail(c, response.DataInvalid, "dept has children, please delete children first")
		case errors.Is(err, errDeptHasUsers):
			response.Fail(c, response.DataInvalid, "dept has users assigned")
		default:
			response.Fail(c, response.InternalError, err.Error())
		}
		return
	}
	response.Success[any](c, nil, "delete dept success")
}
