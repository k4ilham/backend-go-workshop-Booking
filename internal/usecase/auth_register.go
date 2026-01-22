package usecase

import (
	"errors"
	"time"

	"be-golang/internal/domain"
	"be-golang/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type AdminRegister struct {
	users  ports.UserRepository
	logger ports.Logger
}

func NewAdminRegister(users ports.UserRepository, logger ports.Logger) *AdminRegister {
	return &AdminRegister{users: users, logger: logger}
}

func (u *AdminRegister) Exec(email, password string) (int64, error) {
	if email == "" || password == "" {
		return 0, errors.New("invalid_input")
	}
	existing, _ := u.users.GetByEmail(email)
	if existing != nil {
		return 0, errors.New("email_exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	now := time.Now().UTC()
	id, err := u.users.Create(domain.User{
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
	})
	if err != nil {
		return 0, err
	}
	_ = u.logger.Log("admin_register", email, now)
	return id, nil
}
