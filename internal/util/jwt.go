package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret []byte
}

func NewJWT(secret string) *JWT {
	return &JWT{secret: []byte(secret)}
}

func (j *JWT) Generate(claims map[string]any, ttl time.Duration) (string, error) {
	c := jwt.MapClaims{}
	for k, v := range claims {
		c[k] = v
	}
	c["exp"] = time.Now().Add(ttl).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(j.secret)
}

func (j *JWT) Parse(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return j.secret, nil
	})
}
