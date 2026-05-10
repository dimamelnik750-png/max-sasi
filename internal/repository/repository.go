package repository

import (
	"errors"
	"sync"

	"maxsasi/internal/todo"
)

var ErrTodoNotFound = errors.New("todo not found")

type TodoRepository interface {
	GetAll() ([]todo.Todo, error)
	GetByID(id string) (todo.Todo, error)
	Create(item todo.Todo) (todo.Todo, error)
	Update(item todo.Todo) (todo.Todo, error)
	Delete(id string) error
}

type InMemoryTodoRepository struct {
	mu    sync.RWMutex
	todos map[string]todo.Todo
}

func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[string]todo.Todo),
	}
}

func (r *InMemoryTodoRepository) GetAll() ([]todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]todo.Todo, 0, len(r.todos))
	for _, item := range r.todos {
		list = append(list, item)
	}

	return list, nil
}

func (r *InMemoryTodoRepository) GetByID(id string) (todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.todos[id]
	if !ok {
		return todo.Todo{}, ErrTodoNotFound
	}

	return item, nil
}

func (r *InMemoryTodoRepository) Create(item todo.Todo) (todo.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.todos[item.ID] = item
	return item, nil
}

func (r *InMemoryTodoRepository) Update(item todo.Todo) (todo.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[item.ID]; !ok {
		return todo.Todo{}, ErrTodoNotFound
	}

	r.todos[item.ID] = item
	return item, nil
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
