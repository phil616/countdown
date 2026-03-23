package handler

import (
	"ops-timer-backend/internal/dto"
	"ops-timer-backend/internal/pkg/response"
	"ops-timer-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService *service.TodoService
}

func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (h *TodoHandler) Create(c *gin.Context) {
	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	todo, err := h.todoService.Create(&req)
	if err != nil {
		response.BusinessError(c, err.Error())
		return
	}
	response.Created(c, todo)
}

func (h *TodoHandler) Get(c *gin.Context) {
	id := c.Param("id")
	todo, err := h.todoService.GetByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, todo)
}

func (h *TodoHandler) List(c *gin.Context) {
	var params dto.TodoQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	todos, total, err := h.todoService.List(&params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMeta(c, todos, &response.Meta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: response.CalculateTotalPages(total, params.PageSize),
	})
}

func (h *TodoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	todo, err := h.todoService.Update(id, &req)
	if err != nil {
		if err == service.ErrTodoNotFound {
			response.NotFound(c, err.Error())
		} else {
			response.BusinessError(c, err.Error())
		}
		return
	}
	response.Success(c, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.todoService.Delete(id); err != nil {
		if err == service.ErrTodoNotFound {
			response.NotFound(c, err.Error())
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}
	response.NoContent(c)
}

func (h *TodoHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateTodoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	todo, err := h.todoService.UpdateStatus(id, req.Status)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, todo)
}

func (h *TodoHandler) BatchAction(c *gin.Context) {
	var req dto.BatchTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	if err := h.todoService.BatchAction(&req); err != nil {
		response.BusinessError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *TodoHandler) CreateGroup(c *gin.Context) {
	var req dto.CreateTodoGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	group, err := h.todoService.CreateGroup(&req)
	if err != nil {
		response.BusinessError(c, err.Error())
		return
	}
	response.Created(c, group)
}

func (h *TodoHandler) ListGroups(c *gin.Context) {
	groups, err := h.todoService.ListGroups()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

func (h *TodoHandler) UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateTodoGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败", nil)
		return
	}

	group, err := h.todoService.UpdateGroup(id, &req)
	if err != nil {
		if err == service.ErrGroupNotFound {
			response.NotFound(c, err.Error())
		} else {
			response.BusinessError(c, err.Error())
		}
		return
	}
	response.Success(c, group)
}

func (h *TodoHandler) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := h.todoService.DeleteGroup(id); err != nil {
		if err == service.ErrGroupNotFound {
			response.NotFound(c, err.Error())
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}
	response.NoContent(c)
}
