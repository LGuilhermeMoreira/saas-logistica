package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Env struct {
	PORT       string
	JWT_SECRET string
	JWT_TTL    time.Duration

	DATABASE_NAME     string
	DATABASE_USER     string
	DATABASE_PASSWORD string
	DATABASE_HOST     string
	DATABASE_SSL_MODE string
	DATABASE_PORT     string

	BUCKET_NAME           string
	URI                   string
	ACCESS_KEY_ID         string
	SECRECT_ACCESS_KEY    string
	USE_SSL               bool
	REGION                string
	SIGNED_URL_EXPIRATION time.Duration
}

func NewEnv() (*Env, error) {
	jwtTTL, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL: %w", err)
	}

	signedURLExpiration, err := time.ParseDuration(os.Getenv("SIGNED_URL_EXPIRATION"))
	if err != nil {
		return nil, fmt.Errorf("invalid SIGNED_URL_EXPIRATION: %w", err)
	}

	useSSL, err := strconv.ParseBool(os.Getenv("USE_SSL"))
	if err != nil {
		return nil, fmt.Errorf("invalid USE_SSL: %w", err)
	}

	return &Env{
		PORT:       os.Getenv("PORT"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		JWT_TTL:    jwtTTL,

		DATABASE_NAME:     os.Getenv("DATABASE_NAME"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
		DATABASE_HOST:     os.Getenv("DATABASE_HOST"),
		DATABASE_SSL_MODE: os.Getenv("DATABASE_SSL_MODE"),
		DATABASE_PORT:     os.Getenv("DATABASE_PORT"),

		BUCKET_NAME:           os.Getenv("BUCKET_NAME"),
		URI:                   os.Getenv("URI"),
		ACCESS_KEY_ID:         os.Getenv("ACCESS_KEY_ID"),
		SECRECT_ACCESS_KEY:    os.Getenv("SECRECT_ACCESS_KEY"),
		USE_SSL:               useSSL,
		REGION:                os.Getenv("REGION"),
		SIGNED_URL_EXPIRATION: signedURLExpiration,
	}, nil
}

func (e *Env) PostgresURI() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		e.DATABASE_USER,
		e.DATABASE_PASSWORD,
		e.DATABASE_HOST,
		e.DATABASE_PORT,
		e.DATABASE_NAME,
		e.DATABASE_SSL_MODE,
	)
}
