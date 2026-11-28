package todo

import "time"

type Todo struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
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
