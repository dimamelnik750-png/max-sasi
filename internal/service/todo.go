package service

import (
	"errors"
	"strings"
	"time"

	"maxsasi/internal/repository"
	"maxsasi/internal/todo"

	"github.com/google/uuid"
)

var ErrTitleRequired = errors.New("title is required")

type TodoService interface {
	GetAllTodos(userID string) ([]todo.Todo, error)
	GetTodoByID(id, userID string) (todo.Todo, error)
	CreateTodo(userID string, input todo.CreateTodoInput) (todo.Todo, error)
	UpdateTodo(id, userID string, input todo.UpdateTodoInput) (todo.Todo, error)
	DeleteTodo(id, userID string) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) GetAllTodos(userID string) ([]todo.Todo, error) {
	return s.repo.GetAll(userID)
}

func (s *todoService) GetTodoByID(id, userID string) (todo.Todo, error) {
	return s.repo.GetByID(id, userID)
}

func (s *todoService) CreateTodo(userID string, input todo.CreateTodoInput) (todo.Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return todo.Todo{}, ErrTitleRequired
	}

	item := todo.Todo{
		ID:        uuid.NewString(),
		UserID:    userID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(item)
}

func (s *todoService) UpdateTodo(id, userID string, input todo.UpdateTodoInput) (todo.Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return todo.Todo{}, ErrTitleRequired
	}

	item, err := s.repo.GetByID(id, userID)
	if err != nil {
		return todo.Todo{}, err
	}

	item.Title = title
	item.Completed = input.Completed

	return s.repo.Update(item)
}

func (s *todoService) DeleteTodo(id, userID string) error {
	return s.repo.Delete(id, userID)
}
