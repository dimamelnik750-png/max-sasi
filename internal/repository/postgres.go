package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"maxsasi/internal/todo"
)

type PostgresTodoRepository struct {
	db *sql.DB
}

func NewPostgresTodoRepository(db *sql.DB) *PostgresTodoRepository {
	return &PostgresTodoRepository{db: db}
}

func (r *PostgresTodoRepository) GetAll() ([]todo.Todo, error) {
	rows, err := r.db.Query(`
		SELECT id, title, completed, created_at
		FROM todos
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]todo.Todo, 0)
	for rows.Next() {
		var item todo.Todo
		if err := rows.Scan(&item.ID, &item.Title, &item.Completed, &item.CreatedAt); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *PostgresTodoRepository) GetByID(id string) (todo.Todo, error) {
	var item todo.Todo

	err := r.db.QueryRow(`
		SELECT id, title, completed, created_at
		FROM todos
		WHERE id = $1
	`, id).Scan(&item.ID, &item.Title, &item.Completed, &item.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return todo.Todo{}, ErrTodoNotFound
		}

		return todo.Todo{}, err
	}

	return item, nil
}

func (r *PostgresTodoRepository) Create(item todo.Todo) (todo.Todo, error) {
	created := item

	err := r.db.QueryRow(`
		INSERT INTO todos (id, title, completed, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, completed, created_at
	`, item.ID, item.Title, item.Completed, item.CreatedAt).
		Scan(&created.ID, &created.Title, &created.Completed, &created.CreatedAt)
	if err != nil {
		return todo.Todo{}, fmt.Errorf("create todo: %w", err)
	}

	return created, nil
}

func (r *PostgresTodoRepository) Update(item todo.Todo) (todo.Todo, error) {
	updated := item

	err := r.db.QueryRow(`
		UPDATE todos
		SET title = $2, completed = $3
		WHERE id = $1
		RETURNING id, title, completed, created_at
	`, item.ID, item.Title, item.Completed).
		Scan(&updated.ID, &updated.Title, &updated.Completed, &updated.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return todo.Todo{}, ErrTodoNotFound
		}

		return todo.Todo{}, fmt.Errorf("update todo: %w", err)
	}

	return updated, nil
}

func (r *PostgresTodoRepository) Delete(id string) error {
	result, err := r.db.Exec(`
		DELETE FROM todos
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}
