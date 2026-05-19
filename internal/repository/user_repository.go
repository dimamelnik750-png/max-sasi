package repository

import (
	"database/sql"
	"errors"

	"maxsasi/internal/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository interface {
	Create(username, password string) (user.User, error)
	GetByUsername(username string) (user.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(username, password string) (user.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user.User{}, err
	}

	var u user.User
	err = r.db.QueryRow(`
        INSERT INTO users (id, username, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, username, created_at
    `, uuid.NewString(), username, string(hash)).
		Scan(&u.ID, &u.Username, &u.CreatedAt)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return user.User{}, ErrUserAlreadyExists
		}
		return user.User{}, err
	}

	return u, nil
}

func (r *PostgresUserRepository) GetByUsername(username string) (user.User, error) {
	var u user.User

	err := r.db.QueryRow(`
        SELECT id, username, password_hash, created_at
        FROM users
        WHERE username = $1
    `, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, ErrUserNotFound
		}
		return user.User{}, err
	}

	return u, nil
}
