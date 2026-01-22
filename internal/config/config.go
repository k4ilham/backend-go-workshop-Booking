package config

import "time"

type Config struct {
	PostgresDSN    string
	JWTSecret      string
	TursoURL       string
	TursoToken     string
	N8NWebhookURL  string
	ServerAddr     string
	TokenTTL       time.Duration
	AdminOnlyPaths []string
}
