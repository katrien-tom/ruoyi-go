package menu

import (
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
	data, err := h.service.GetMenuList(c.Request.Context())
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get menu list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) TreeSelect(c *gin.Context) {
	data, err := h.service.GetMenuTreeSelect(c.Request.Context())
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get menu tree")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) GetInfo(c *gin.Context) {
	menuID, err := strconv.ParseInt(c.Param("menuId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid menuId")
		return
	}

	data, err := h.service.GetByID(c.Request.Context(), menuID)
	if err != nil {
		response.Fail(c, response.NotFound, "menu not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) Add(c *gin.Context) {
	var req AddMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.AddMenu(c.Request.Context(), req); err != nil {
		response.Fail(c, response.InternalError, "failed to add menu")
		return
	}
	response.Success[any](c, nil, "add menu success")
}

func (h *Handler) Edit(c *gin.Context) {
	var req EditMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.EditMenu(c.Request.Context(), req); err != nil {
		response.Fail(c, response.InternalError, "failed to edit menu")
		return
	}
	response.Success[any](c, nil, "edit menu success")
}

func (h *Handler) Delete(c *gin.Context) {
	menuID, err := strconv.ParseInt(c.Param("menuId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid menuId")
		return
	}

	if err := h.service.DeleteMenu(c.Request.Context(), menuID); err != nil {
		response.Fail(c, response.InternalError, err.Error())
		return
	}
	response.Success[any](c, nil, "delete menu success")
}

func (h *Handler) RoleMenuTreeSelect(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid roleId")
		return
	}

	data, err := h.service.GetRoleMenuTreeSelect(c.Request.Context(), roleID)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get role menu tree")
		return
	}
	response.SuccessNoMsg(c, data)
}
