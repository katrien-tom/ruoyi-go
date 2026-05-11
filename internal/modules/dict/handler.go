package dict

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

// DictType handlers

func (h *Handler) TypeList(c *gin.Context) {
	var req DictTypeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.ListDictTypes(c.Request.Context(), req, pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get dict type list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) TypeGetInfo(c *gin.Context) {
	dictID, err := strconv.ParseInt(c.Param("dictId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid dictId")
		return
	}

	data, err := h.service.GetDictTypeByID(c.Request.Context(), dictID)
	if err != nil {
		response.Fail(c, response.NotFound, "dict type not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) TypeAdd(c *gin.Context) {
	var req AddDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.AddDictType(c.Request.Context(), req); err != nil {
		if errors.Is(err, errDictTypeDuplicate) {
			response.Fail(c, response.DataAlreadyExist, "dict type already exists")
			return
		}
		response.Fail(c, response.InternalError, "failed to add dict type")
		return
	}
	response.Success[any](c, nil, "add dict type success")
}

func (h *Handler) TypeEdit(c *gin.Context) {
	var req EditDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.EditDictType(c.Request.Context(), req); err != nil {
		if errors.Is(err, errDictTypeDuplicate) {
			response.Fail(c, response.DataAlreadyExist, "dict type already exists")
			return
		}
		response.Fail(c, response.InternalError, "failed to edit dict type")
		return
	}
	response.Success[any](c, nil, "edit dict type success")
}

func (h *Handler) TypeDelete(c *gin.Context) {
	ids := strings.Split(c.Param("dictIds"), ",")
	dictIDs := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			response.Fail(c, response.ParamError, "invalid dict id: "+s)
			return
		}
		dictIDs = append(dictIDs, id)
	}

	if err := h.service.DeleteDictTypes(c.Request.Context(), dictIDs); err != nil {
		response.Fail(c, response.InternalError, "failed to delete dict types")
		return
	}
	response.Success[any](c, nil, "delete dict types success")
}

func (h *Handler) TypeRefreshCache(c *gin.Context) {
	response.Success[any](c, nil, "refresh cache success")
}

// DictData handlers

func (h *Handler) DataList(c *gin.Context) {
	var req DictDataListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.ListDictDatas(c.Request.Context(), req, pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get dict data list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) DataByType(c *gin.Context) {
	dictType := c.Param("dictType")
	data, err := h.service.GetDictDataByType(c.Request.Context(), dictType)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get dict data")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) DataGetInfo(c *gin.Context) {
	dictCode, err := strconv.ParseInt(c.Param("dictCode"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid dictCode")
		return
	}

	data, err := h.service.GetDictDataByID(c.Request.Context(), dictCode)
	if err != nil {
		response.Fail(c, response.NotFound, "dict data not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) DataAdd(c *gin.Context) {
	var req AddDictDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.AddDictData(c.Request.Context(), req); err != nil {
		response.Fail(c, response.InternalError, "failed to add dict data")
		return
	}
	response.Success[any](c, nil, "add dict data success")
}

func (h *Handler) DataEdit(c *gin.Context) {
	var req EditDictDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.EditDictData(c.Request.Context(), req); err != nil {
		response.Fail(c, response.InternalError, "failed to edit dict data")
		return
	}
	response.Success[any](c, nil, "edit dict data success")
}

func (h *Handler) DataDelete(c *gin.Context) {
	ids := strings.Split(c.Param("dictCodes"), ",")
	dictCodes := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			response.Fail(c, response.ParamError, "invalid dict code: "+s)
			return
		}
		dictCodes = append(dictCodes, id)
	}

	if err := h.service.DeleteDictDatas(c.Request.Context(), dictCodes); err != nil {
		response.Fail(c, response.InternalError, "failed to delete dict data")
		return
	}
	response.Success[any](c, nil, "delete dict data success")
}
