package config

import (
	"fmt"
	"os"
	"time"
)

type Env struct {
	PORT       string
	JWT_SECRET string
	JWT_TTL    time.Duration

	DATABASE_NAME,
	DATABASE_USER,
	DATABASE_PASSWORD,
	DATABASE_HOST,
	DATABASE_SSL_MODE,
	DATABASE_PORT string
}

func NewEnv() (*Env, error) {
	jwtTTLDuration, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		return nil, err
	}
	return &Env{
		PORT:       os.Getenv("PORT"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		JWT_TTL:    jwtTTLDuration,
	}, nil
}

func (e *Env) PostgresURI() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		e.DATABASE_USER,
		e.DATABASE_PASSWORD,
		e.DATABASE_HOST,
		e.DATABASE_PORT,
		e.DATABASE_NAME,
		e.DATABASE_SSL_MODE)
}
