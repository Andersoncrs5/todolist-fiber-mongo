package handlers

import (
	"strconv"
	"todolist/dtos"
	"todolist/models"
	"todolist/services"
	"todolist/utils"

	"github.com/gofiber/fiber/v2"
)

type TodoHandler interface {
	
}

type todoHandler struct {
	service services.TodoService	
}

func NewTodoHandler(s services.TodoService) TodoHandler {
	return &todoHandler{
		service: s,
	}
}

func (h *todoHandler) CreateTodo(c *fiber.Ctx) error {
	var req dtos.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Parser error: "+err.Error())
	}

	todo := &models.Todo{
		Title: req.Title,
		Description: req.Description,
		Complete: false,
	}

	createdTodo, err := h.service.CreateTodo(c.Context(), todo);
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Fail the to create new todo")
	}

	return utils.SendSuccessResponse(c, fiber.StatusCreated, "Todo created with success!", createdTodo)
}

func (h *todoHandler) GetAllTodos(c *fiber.Ctx) error {
	titleFilter := c.Query("title", "");
	completeStr := c.Query("complete", "");

	var completeFilter *bool
	if completeStr == "" {
		parsedComplete, err := strconv.ParseBool(completeStr);
		if err != nil { return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Parâmetro 'complete' inválido. Use 'true' ou 'false'.") }

		completeFilter = &parsedComplete
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1;
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"));
	if err != nil || pageSize < 1 {
		pageSize = 10;
	}

	todos, total, err := h.service.GetAllTodos(c.Context(), titleFilter, completeFilter, page, pageSize);
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Error the get all todos! details: "+ err.Error())
	}

	meta := fiber.Map{
		"total": total,
		"page": page,
		"pageSize": pageSize,
		"totalPages": (total + int64(pageSize) - 1) / int64(pageSize), 
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Todos listed!", todos, meta);
}

func (h *todoHandler) GetTodoByID(c *fiber.Ctx) error {
	id := c.Params("id");
	todo, err := h.service.GetTodoByID(c.Context(), id)
	if err != nil {
		if err.Error() == "Task not found" || err.Error() == "Id is Invalid" {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, err.Error())
		}

		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Fail the to find task: "+err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Task founded", todo);
}

func (h *todoHandler) UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id");

	var req dtos.UpdateTodoRequest;
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusOK, "Validation errors: " + err.Error());
	}

	todoToUpdate := &models.Todo{}
	if req.Title != "" {
		todoToUpdate.Title = req.Title
	}

	if req.Description != "" {
		todoToUpdate.Description = req.Description
	}

	if req.Complete != nil {
		todoToUpdate.Complete = *req.Complete
	}

	updatedTodo, err := h.service.UpdateTodo(c.Context(), id, todoToUpdate);
	if err != nil {
		if err.Error() == "Todo not found" || err.Error() == "Id is Invalid" || err.Error() == "Error the check if todo exists!" {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Fail the to update todo: "+err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Task updated", updatedTodo)
}

func (h *todoHandler) DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteTodo(c.Context(), id);

	if err != nil {
		if err.Error() == "Error the check if todo exists!" || err.Error() == "Id is Invalid" || err.Error() == "Todo not found" {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Falha ao deletar tarefa: "+err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusNoContent, "Task deleted!", nil)
}

func (h *todoHandler) ToggleTodoComplete(c *fiber.Ctx) error {
	id := c.Params("id")

	updatedTodo, err := h.service.ToggleTodoComplete(c.Context(), id)
	if err != nil {
		if err.Error() == "Task not found" || err.Error() == "Id is Invalid" {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Fail the modificed task status completed : "+err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Status task completed changed", updatedTodo)
}
