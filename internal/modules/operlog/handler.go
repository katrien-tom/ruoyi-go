package operlog

import (
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
	var req OperLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.ListLogs(c.Request.Context(), req, pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get oper log list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) Delete(c *gin.Context) {
	ids := strings.Split(c.Param("operIds"), ",")
	logIDs := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			response.Fail(c, response.ParamError, "invalid oper id: "+s)
			return
		}
		logIDs = append(logIDs, id)
	}

	if err := h.service.DeleteLogs(c.Request.Context(), logIDs); err != nil {
		response.Fail(c, response.InternalError, "failed to delete oper logs")
		return
	}
	response.Success[any](c, nil, "delete oper logs success")
}

func (h *Handler) Clean(c *gin.Context) {
	if err := h.service.Clean(c.Request.Context()); err != nil {
		response.Fail(c, response.InternalError, "failed to clean oper logs")
		return
	}
	response.Success[any](c, nil, "clean oper logs success")
}
