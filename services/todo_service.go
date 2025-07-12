package services

import (
	"context"
	"errors"
	"fmt"
	"todolist/models"
	"todolist/repositories"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoService interface {
	CreateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error)
	GetAllTodos(ctx context.Context, title string, complete *bool, page, pagSize int) ([]models.Todo, int64, error)
	GetTodoByID(ctx context.Context, id string) (*models.Todo, error)
	UpdateTodo(ctx context.Context, id string, todo *models.Todo) (*models.Todo, error)
	DeleteTodo(ctx context.Context, id string) error
	ToggleTodoComplete(ctx context.Context, id string) (*models.Todo, error)
}

type todoService struct {
	repo repositories.TodoRepository
	validate *validator.Validate
}

func NewTodoService(repo repositories.TodoRepository) TodoService {
	return &todoService {
		repo: repo,
		validate: validator.New(),
	}
}

func (s *todoService) CreateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error)  {
	if err:= s.validate.StructPartial(todo); err != nil {
		return nil, fmt.Errorf("Error the update task");
	}

	return s.repo.CreateTodo(ctx, todo);
}

func (s *todoService) GetAllTodos(ctx context.Context, title string, complete *bool, page, pagSize int) ([]models.Todo, int64, error) {
	filter := make(primitive.M)

	if title != "" {
		filter["title"] = primitive.Regex{Pattern: title, Options: "i"}
	}

	if complete != nil {
		filter["complete"] = *complete
	}

	skip := int64((page - 1) * pagSize)
	limit := int64(pagSize);

	todos, err := s.repo.GetAll(ctx, filter, limit, skip);
	if err != nil {
		return nil, 0, err;
	}

	total, err := s.repo.CountTodos(ctx, filter);
	if err != nil {
		return nil, 0, err;
	}

	return todos, total, nil
}

func (s *todoService) GetTodoByID(ctx context.Context, id string) (*models.Todo, error)  {
	objectID, err := primitive.ObjectIDFromHex(id);
	if err != nil {
		return nil, fmt.Errorf("Id is Invalid")
	}

	todo, err := s.repo.GetTodoByID(ctx, objectID);
	if err != nil {
		return nil, err;
	}

	if todo == nil {
		return nil, errors.New("Task not found");
	}

	return todo, nil;
}

func (s *todoService) UpdateTodo(ctx context.Context, id string, todo *models.Todo) (*models.Todo, error)  {
	objectID, err := primitive.ObjectIDFromHex(id);
	if err != nil {
		return nil, fmt.Errorf("Id is Invalid")
	}

	if err := s.validate.StructPartial(todo, "Title", "Description", "Complete"); err != nil {
		return nil, fmt.Errorf("Validation error: %w", err);
	}

	checkTodo, err := s.repo.GetTodoByID(ctx, objectID);
	if err != nil {
		return nil, fmt.Errorf("Error the check if todo exists!");
	}

	if checkTodo == nil {
		return nil, fmt.Errorf("Todo not found");
	}

	if todo.Title != "" {
		checkTodo.Title = todo.Title
	}
	if todo.Description != "" {
		checkTodo.Description = todo.Description
	}

	checkTodo.Complete = todo.Complete;

	updateTodo, err := s.repo.UpdateTodo(ctx, objectID, checkTodo)
	if err != nil {
		return nil, fmt.Errorf("Error the update if todo exists!");
	}

	return updateTodo, nil;
}

func (s *todoService) DeleteTodo(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id);
	if err != nil {
		return fmt.Errorf("Id is invalid! Error: %w", err)
	}

	if err := s.repo.DeleteTodo(ctx, objectID); err != nil {
		return err;
	}

	return err;
}

func (s *todoService) ToggleTodoComplete(ctx context.Context, id string) (*models.Todo, error) {
	objectID, err := primitive.ObjectIDFromHex(id);
	if err != nil {
		return nil, fmt.Errorf("Id is Invalid")
	}

	todo, err := s.repo.GetTodoByID(ctx, objectID);
	if err != nil {
		return nil, fmt.Errorf("Error the check if todo exists before change status complete")
	}
	
	if todo == nil {
		return nil, errors.New("Task not found");
	}

	todo.Complete = !todo.Complete;

	changeTodo, err := s.repo.UpdateTodo(ctx, objectID, todo);
	if err != nil {
		return nil, err
	}

	return changeTodo, nil
}