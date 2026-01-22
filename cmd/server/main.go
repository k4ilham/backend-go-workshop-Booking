package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"be-golang/internal/app"
	"be-golang/internal/config"
)

func main() {
	loadEnvFile(".env")
	cfg := config.Config{
		PostgresDSN:    firstNonEmpty(os.Getenv("POSTGRES_DSN"), os.Getenv("DATABASE_URL")),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		TursoURL:       os.Getenv("TURSO_URL"),
		TursoToken:     os.Getenv("TURSO_TOKEN"),
		N8NWebhookURL:  os.Getenv("N8N_WEBHOOK_URL"),
		ServerAddr:     envString("SERVER_ADDR", ":8080"),
		TokenTTL:       envDuration("TOKEN_TTL", time.Hour*24),
		AdminOnlyPaths: []string{"/admin", "/services"},
	}
	if p := os.Getenv("PORT"); p != "" {
		cfg.ServerAddr = ":" + p
	}
	if cfg.PostgresDSN == "" || cfg.JWTSecret == "" {
		log.Fatal("missing required environment variables")
	}
	err := app.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func envString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func envDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return time.Duration(d) * time.Second
}

func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		i := strings.Index(line, "=")
		if i <= 0 {
			continue
		}
		k := strings.TrimSpace(line[:i])
		v := strings.TrimSpace(line[i+1:])
		if len(v) >= 2 && ((v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'')) {
			v = v[1 : len(v)-1]
		}
		_ = os.Setenv(k, v)
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
