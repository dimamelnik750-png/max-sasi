package service

import (
	"context"
	"fmt"
	"time"

	"maxsasi/internal/cache"
	"maxsasi/internal/todo"
)

const (
	todosCacheKey   = "todos:all"
	todoCachePrefix = "todos:"
	cacheTTL        = 5 * time.Minute
)

type cachedTodoService struct {
	service TodoService
	cache   cache.Cache
}

func NewCachedTodoService(service TodoService, cache cache.Cache) TodoService {
	return &cachedTodoService{service: service, cache: cache}
}

func (s *cachedTodoService) GetAllTodos() ([]todo.Todo, error) {
	ctx := context.Background()

	var items []todo.Todo
	if err := s.cache.Get(ctx, todosCacheKey, &items); err == nil {
		return items, nil
	}

	items, err := s.service.GetAllTodos()
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, todosCacheKey, items, cacheTTL)

	return items, nil
}

func (s *cachedTodoService) GetTodoByID(id string) (todo.Todo, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", todoCachePrefix, id)

	var item todo.Todo
	if err := s.cache.Get(ctx, key, &item); err == nil {
		return item, nil
	}

	item, err := s.service.GetTodoByID(id)
	if err != nil {
		return todo.Todo{}, err
	}

	_ = s.cache.Set(ctx, key, item, cacheTTL)

	return item, nil
}

func (s *cachedTodoService) CreateTodo(input todo.CreateTodoInput) (todo.Todo, error) {
	ctx := context.Background()

	item, err := s.service.CreateTodo(input)
	if err != nil {
		return todo.Todo{}, err
	}

	_ = s.cache.Delete(ctx, todosCacheKey)

	return item, nil
}

func (s *cachedTodoService) UpdateTodo(id string, input todo.UpdateTodoInput) (todo.Todo, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", todoCachePrefix, id)

	item, err := s.service.UpdateTodo(id, input)
	if err != nil {
		return todo.Todo{}, err
	}

	_ = s.cache.Delete(ctx, key)
	_ = s.cache.Delete(ctx, todosCacheKey)

	return item, nil
}

func (s *cachedTodoService) DeleteTodo(id string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", todoCachePrefix, id)

	if err := s.service.DeleteTodo(id); err != nil {
		return err
	}

	_ = s.cache.Delete(ctx, key)
	_ = s.cache.Delete(ctx, todosCacheKey)

	return nil
}
