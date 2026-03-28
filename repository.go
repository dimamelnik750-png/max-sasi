package main

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrTodoNotFound = errors.New("todo not found")

type TodoRepository interface {
	GetAll() []Todo
	GetByID(id string) (Todo, error)
	Create(todo Todo) Todo
	Update(todo Todo) (Todo, error)
	Delete(id string) error
	NextID() string
}

type InMemoryTodoRepository struct {
	mu    sync.RWMutex
	todos map[string]Todo
}

func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[string]Todo),
	}
}

func (r *InMemoryTodoRepository) GetAll() []Todo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		list = append(list, todo)
	}

	return list
}

func (r *InMemoryTodoRepository) GetByID(id string) (Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, ok := r.todos[id]
	if !ok {
		return Todo{}, ErrTodoNotFound
	}

	return todo, nil
}

func (r *InMemoryTodoRepository) Create(todo Todo) Todo {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.todos[todo.ID] = todo
	return todo
}

func (r *InMemoryTodoRepository) Update(todo Todo) (Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[todo.ID]; !ok {
		return Todo{}, ErrTodoNotFound
	}

	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *InMemoryTodoRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[id]; !ok {
		return ErrTodoNotFound
	}

	delete(r.todos, id)
	return nil
}

func (r *InMemoryTodoRepository) NextID() string {
	return uuid.NewString()
}
