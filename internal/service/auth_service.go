package service

import (
	"errors"

	"maxsasi/internal/auth"
	"maxsasi/internal/repository"
	"maxsasi/internal/user"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUsernameRequired = errors.New("username is required")
var ErrPasswordRequired = errors.New("password is required")

type AuthService interface {
	Register(input user.RegisterInput) (user.User, error)
	Login(input user.LoginInput) (auth.TokenPair, error)
	Refresh(refreshToken string) (auth.TokenPair, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(input user.RegisterInput) (user.User, error) {
	if input.Username == "" {
		return user.User{}, ErrUsernameRequired
	}
	if input.Password == "" {
		return user.User{}, ErrPasswordRequired
	}

	return s.userRepo.Create(input.Username, input.Password)
}

func (s *authService) Login(input user.LoginInput) (auth.TokenPair, error) {
	if input.Username == "" {
		return auth.TokenPair{}, ErrUsernameRequired
	}
	if input.Password == "" {
		return auth.TokenPair{}, ErrPasswordRequired
	}

	u, err := s.userRepo.GetByUsername(input.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return auth.TokenPair{}, ErrInvalidCredentials
		}
		return auth.TokenPair{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		return auth.TokenPair{}, ErrInvalidCredentials
	}

	return auth.GenerateTokenPair(u.ID, u.Username, s.jwtSecret)
}

func (s *authService) Refresh(refreshToken string) (auth.TokenPair, error) {
	claims, err := auth.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return auth.TokenPair{}, ErrInvalidCredentials
	}

	return auth.GenerateTokenPair(claims.UserID, claims.Username, s.jwtSecret)
}
