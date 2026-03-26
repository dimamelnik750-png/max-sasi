package main

import (
	"errors"
	"strings"
	"time"
)

var ErrTitleRequired = errors.New("title is required")

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTodoInput struct {
	Title string `json:"title"`
}

type UpdateTodoInput struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoService interface {
	GetAllTodos() []Todo
	GetTodoByID(id string) (Todo, error)
	CreateTodo(input CreateTodoInput) (Todo, error)
	UpdateTodo(id string, input UpdateTodoInput) (Todo, error)
	DeleteTodo(id string) error
}

type todoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) GetAllTodos() []Todo {
	return s.repo.GetAll()
}

func (s *todoService) GetTodoByID(id string) (Todo, error) {
	return s.repo.GetByID(id)
}

func (s *todoService) CreateTodo(input CreateTodoInput) (Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Todo{}, ErrTitleRequired
	}

	todo := Todo{
		ID:        s.repo.NextID(),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(todo), nil
}

func (s *todoService) UpdateTodo(id string, input UpdateTodoInput) (Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Todo{}, ErrTitleRequired
	}

	todo, err := s.repo.GetByID(id)
	if err != nil {
		return Todo{}, err
	}

	todo.Title = title
	todo.Completed = input.Completed

	return s.repo.Update(todo)
}

func (s *todoService) DeleteTodo(id string) error {
	return s.repo.Delete(id)
}
