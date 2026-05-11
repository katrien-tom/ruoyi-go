package job

import (
	"github.com/gin-gonic/gin"

	"github.com/banyejiu/ruoyi-go/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) JobList(c *gin.Context) {
	response.SuccessNoMsg(c, "success job list")
}

func (h *Handler) JobLogList(c *gin.Context) {
	response.SuccessNoMsg(c, "success job log list")
}

func (h *Handler) GetInfo(c *gin.Context) {
	response.SuccessNoMsg(c, "success getInfo")
}

func (h *Handler) Add(c *gin.Context) {
	response.SuccessNoMsg(c, "success add")
}

func (h *Handler) Edit(c *gin.Context) {
	response.SuccessNoMsg(c, "success edit")
}

func (h *Handler) Delete(c *gin.Context) {
	response.SuccessNoMsg(c, "success delete")
}

func (h *Handler) Run(c *gin.Context) {
	response.SuccessNoMsg(c, "success run")
}
