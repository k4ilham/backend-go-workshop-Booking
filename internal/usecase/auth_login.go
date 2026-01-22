package usecase

import (
	"errors"
	"time"

	"be-golang/internal/ports"
	"be-golang/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AuthLogin struct {
	users    ports.UserRepository
	logger   ports.Logger
	jwt      *util.JWT
	tokenTTL time.Duration
}

func NewAuthLogin(users ports.UserRepository, logger ports.Logger, jwt *util.JWT, ttl time.Duration) *AuthLogin {
	return &AuthLogin{users: users, logger: logger, jwt: jwt, tokenTTL: ttl}
}

func (a *AuthLogin) Exec(email, password string) (string, error) {
	u, err := a.users.GetByEmail(email)
	if err != nil || u == nil {
		return "", errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	claims := map[string]any{"sub": u.ID, "email": u.Email}
	token, err := a.jwt.Generate(claims, a.tokenTTL)
	if err != nil {
		return "", err
	}
	_ = a.logger.Log("admin_login", u.Email, time.Now().UTC())
	return token, nil
}
