package post

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
	var req PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	data, err := h.service.ListPosts(c.Request.Context(), req, pageNum, pageSize)
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get post list")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) GetInfo(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		response.Fail(c, response.ParamError, "invalid postId")
		return
	}

	data, err := h.service.GetByID(c.Request.Context(), postID)
	if err != nil {
		response.Fail(c, response.NotFound, "post not found")
		return
	}
	response.SuccessNoMsg(c, data)
}

func (h *Handler) Add(c *gin.Context) {
	var req AddPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.AddPost(c.Request.Context(), req); err != nil {
		switch {
		case errors.Is(err, errPostCodeDuplicate):
			response.Fail(c, response.DataAlreadyExist, "post code already exists")
		case errors.Is(err, errPostNameDuplicate):
			response.Fail(c, response.DataAlreadyExist, "post name already exists")
		default:
			response.Fail(c, response.InternalError, "failed to add post")
		}
		return
	}
	response.Success[any](c, nil, "add post success")
}

func (h *Handler) Edit(c *gin.Context) {
	var req EditPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError, validation.TranslateError(err))
		return
	}

	if err := h.service.EditPost(c.Request.Context(), req); err != nil {
		switch {
		case errors.Is(err, errPostCodeDuplicate):
			response.Fail(c, response.DataAlreadyExist, "post code already exists")
		case errors.Is(err, errPostNameDuplicate):
			response.Fail(c, response.DataAlreadyExist, "post name already exists")
		default:
			response.Fail(c, response.InternalError, "failed to edit post")
		}
		return
	}
	response.Success[any](c, nil, "edit post success")
}

func (h *Handler) Delete(c *gin.Context) {
	ids := strings.Split(c.Param("postIds"), ",")
	postIDs := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			response.Fail(c, response.ParamError, "invalid post id: "+s)
			return
		}
		postIDs = append(postIDs, id)
	}

	if err := h.service.DeletePosts(c.Request.Context(), postIDs); err != nil {
		response.Fail(c, response.InternalError, "failed to delete posts")
		return
	}
	response.Success[any](c, nil, "delete posts success")
}

func (h *Handler) OptionSelect(c *gin.Context) {
	data, err := h.service.OptionSelect(c.Request.Context())
	if err != nil {
		response.Fail(c, response.InternalError, "failed to get post options")
		return
	}
	response.SuccessNoMsg(c, data)
}
