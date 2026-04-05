package main

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type PostgresTodoRepository struct {
	db *sql.DB
}

func NewPostgresTodoRepository(db *sql.DB) *PostgresTodoRepository {
	return &PostgresTodoRepository{db: db}
}

func (r *PostgresTodoRepository) GetAll() ([]Todo, error) {
	rows, err := r.db.Query(`
		SELECT id, title, completed, created_at
		FROM todos
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *PostgresTodoRepository) GetByID(id string) (Todo, error) {
	var todo Todo

	err := r.db.QueryRow(`
		SELECT id, title, completed, created_at
		FROM todos
		WHERE id = $1
	`, id).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, ErrTodoNotFound
		}

		return Todo{}, err
	}

	return todo, nil
}

func (r *PostgresTodoRepository) Create(todo Todo) (Todo, error) {
	created := todo

	err := r.db.QueryRow(`
		INSERT INTO todos (id, title, completed, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, completed, created_at
	`, todo.ID, todo.Title, todo.Completed, todo.CreatedAt).
		Scan(&created.ID, &created.Title, &created.Completed, &created.CreatedAt)
	if err != nil {
		return Todo{}, err
	}

	return created, nil
}

func (r *PostgresTodoRepository) Update(todo Todo) (Todo, error) {
	updated := todo

	err := r.db.QueryRow(`
		UPDATE todos
		SET title = $2, completed = $3
		WHERE id = $1
		RETURNING id, title, completed, created_at
	`, todo.ID, todo.Title, todo.Completed).
		Scan(&updated.ID, &updated.Title, &updated.Completed, &updated.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, ErrTodoNotFound
		}

		return Todo{}, err
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

func (r *PostgresTodoRepository) NextID() string {
	return uuid.NewString()
}
