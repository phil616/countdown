package handler

import (
	"ops-timer-backend/internal/dto"
	"ops-timer-backend/internal/pkg/response"
	"ops-timer-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	unitService *service.UnitService
}

func NewUnitHandler(unitService *service.UnitService) *UnitHandler {
	return &UnitHandler{unitService: unitService}
}

func (h *UnitHandler) Create(c *gin.Context) {
	var req dto.CreateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	unit, err := h.unitService.Create(&req)
	if err != nil {
		response.BusinessError(c, err.Error())
		return
	}
	response.Created(c, unit)
}

func (h *UnitHandler) Get(c *gin.Context) {
	id := c.Param("id")
	unit, err := h.unitService.GetByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, unit)
}

func (h *UnitHandler) List(c *gin.Context) {
	var params dto.UnitQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	units, total, err := h.unitService.List(&params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMeta(c, units, &response.Meta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: response.CalculateTotalPages(total, params.PageSize),
	})
}

func (h *UnitHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	unit, err := h.unitService.Update(id, &req)
	if err != nil {
		if err == service.ErrUnitNotFound {
			response.NotFound(c, err.Error())
		} else {
			response.BusinessError(c, err.Error())
		}
		return
	}
	response.Success(c, unit)
}

func (h *UnitHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.unitService.Delete(id); err != nil {
		if err == service.ErrUnitNotFound {
			response.NotFound(c, err.Error())
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}
	response.NoContent(c)
}

func (h *UnitHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUnitStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	unit, err := h.unitService.UpdateStatus(id, req.Status)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, unit)
}

func (h *UnitHandler) Step(c *gin.Context) {
	id := c.Param("id")
	var req dto.StepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	unit, err := h.unitService.Step(id, &req)
	if err != nil {
		if err == service.ErrUnitNotFound {
			response.NotFound(c, err.Error())
		} else {
			response.BusinessError(c, err.Error())
		}
		return
	}
	response.Success(c, unit)
}

func (h *UnitHandler) SetValue(c *gin.Context) {
	id := c.Param("id")
	var req dto.SetValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	unit, err := h.unitService.SetValue(id, &req)
	if err != nil {
		if err == service.ErrUnitNotFound {
			response.NotFound(c, err.Error())
		} else {
			response.BusinessError(c, err.Error())
		}
		return
	}
	response.Success(c, unit)
}

func (h *UnitHandler) GetLogs(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.unitService.GetLogs(id, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMeta(c, logs, &response.Meta{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: response.CalculateTotalPages(total, pageSize),
	})
}

func (h *UnitHandler) Summary(c *gin.Context) {
	summary, err := h.unitService.GetSummary()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, summary)
}
