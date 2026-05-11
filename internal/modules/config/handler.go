package config

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
	var req ConfigListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.ListConfigs(c.Request.Context(), req, pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get config list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) GetByKey(c *gin.Context) {
	key := c.Param("configKey")
	data, err := h.service.GetByKey(c.Request.Context(), key)
	if err != nil {
		response.Fail(c, response.NotFound, "config not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) GetInfo(c *gin.Context) {
	configID, err := strconv.Atoi(c.Param("configId"))
	if err != nil {
		response.Fail(c, response.ParamError, "invalid configId")
		return
	}

	data, err := h.service.GetByID(c.Request.Context(), configID)
	if err != nil {
		response.Fail(c, response.NotFound, "config not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) Edit(c *gin.Context) {
	var req EditConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.EditConfig(c.Request.Context(), req); err != nil {
		if errors.Is(err, errConfigKeyDuplicate) {
			response.Fail(c, response.DataAlreadyExist, "config key already exists")
			return
		}
		response.Fail(c, response.InternalError, "failed to edit config")
		return
	}
	response.Success[any](c, nil, "edit config success")
}

func (h *Handler) Delete(c *gin.Context) {
	ids := strings.Split(c.Param("configIds"), ",")
	configIDs := make([]int, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			response.Fail(c, response.ParamError, "invalid config id: "+s)
			return
		}
		configIDs = append(configIDs, id)
	}

	if err := h.service.DeleteConfigs(c.Request.Context(), configIDs); err != nil {
		response.Fail(c, response.InternalError, "failed to delete configs")
		return
	}
	response.Success[any](c, nil, "delete configs success")
}

func (h *Handler) RefreshCache(c *gin.Context) {
	response.Success[any](c, nil, "refresh cache success")
}
