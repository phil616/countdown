package handler

import (
	"ops-timer-backend/internal/dto"
	"ops-timer-backend/internal/pkg/response"
	"ops-timer-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifService *service.NotificationService
}

func NewNotificationHandler(notifService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

func (h *NotificationHandler) List(c *gin.Context) {
	var params dto.NotificationQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	notifications, total, err := h.notifService.List(&params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMeta(c, notifications, &response.Meta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: response.CalculateTotalPages(total, params.PageSize),
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if err := h.notifService.MarkAsRead(id); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	if err := h.notifService.MarkAllAsRead(); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	count, err := h.notifService.UnreadCount()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, dto.UnreadCountResponse{Count: count})
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.notifService.Delete(id); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.NoContent(c)
}
