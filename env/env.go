package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DSN  string
	PORT string
}

func Load() *Env {
	_ = godotenv.Load()

	e := &Env{
		DSN:  getEnv("PG_DSN", ""),
		PORT: getEnv("PG_PORT", "5432"),
	}

	if e.DSN == "" {
		return nil
	}

	return e
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
