package notice

import (
	"github.com/banyejiu/ruoyi-go/internal/security"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/pkg/response"
	"github.com/banyejiu/ruoyi-go/pkg/validation"
)

type Handler struct {
	sessionService *security.SessionService
	service *Service
}

func NewHandler(s *Service, sessionService *security.SessionService) *Handler {
	return &Handler{service: s, sessionService: sessionService}
}

func (h *Handler) List(c *gin.Context) {
	var req NoticeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.ListNotices(c.Request.Context(), req, pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get notice list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) GetInfo(c *gin.Context) {
	noticeID, err := strconv.Atoi(c.Param("noticeId"))
	if err != nil {
		response.Fail(c, response.ParamError, "invalid noticeId")
		return
	}

	data, err := h.service.GetByID(c.Request.Context(), noticeID)
	if err != nil {
		response.Fail(c, response.NotFound, "notice not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) Add(c *gin.Context) {
	var req AddNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.AddNotice(c.Request.Context(), req); err != nil {
		response.Fail(c, response.InternalError, "failed to add notice")
		return
	}
	response.Success[any](c, nil, "add notice success")
}

func (h *Handler) Edit(c *gin.Context) {
	var req EditNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.EditNotice(c.Request.Context(), req); err != nil {
		response.Fail(c, response.InternalError, "failed to edit notice")
		return
	}
	response.Success[any](c, nil, "edit notice success")
}

func (h *Handler) Delete(c *gin.Context) {
	ids := strings.Split(c.Param("noticeIds"), ",")
	noticeIDs := make([]int, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			response.Fail(c, response.ParamError, "invalid notice id: "+s)
			return
		}
		noticeIDs = append(noticeIDs, id)
	}

	if err := h.service.DeleteNotices(c.Request.Context(), noticeIDs); err != nil {
		response.Fail(c, response.InternalError, "failed to delete notices")
		return
	}
	response.Success[any](c, nil, "delete notices success")
}
