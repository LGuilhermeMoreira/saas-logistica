package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Env struct {
	PORT       string
	LOG_MODE   string
	JWT_SECRET string
	JWT_TTL    time.Duration

	OPA_ADDRESS string


	DATABASE_NAME     string
	DATABASE_USER     string
	DATABASE_PASSWORD string
	DATABASE_HOST     string
	DATABASE_SSL_MODE string
	DATABASE_PORT     string

	STORAGE_BUCKET_NAME           string
	STORAGE_URI                   string
	STORAGE_ACCESS_KEY_ID         string
	STORAGE_SECRECT_ACCESS_KEY    string
	STORAGE_USE_SSL               bool
	STORAGE_REGION                string
	STORAGE_SIGNED_URL_EXPIRATION time.Duration
}

func NewEnv() (*Env, error) {
	jwtTTL, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL: %w", err)
	}

	signedURLExpiration, err := time.ParseDuration(os.Getenv("STORAGE_SIGNED_URL_EXPIRATION"))
	if err != nil {
		return nil, fmt.Errorf("invalid SIGNED_URL_EXPIRATION: %w", err)
	}

	useSSL, err := strconv.ParseBool(os.Getenv("STORAGE_USE_SSL"))
	if err != nil {
		return nil, fmt.Errorf("invalid USE_SSL: %w", err)
	}

	opaAddress := os.Getenv("OPA_ADDRESS")
	if opaAddress == "" {
		opaAddress = "http://localhost:8181/v1/data/authz/allow"
	}

	logMode := os.Getenv("LOG_MODE")
	if logMode == "" {
		logMode = "dev"
	}

	return &Env{
		PORT:       os.Getenv("PORT"),
		LOG_MODE:   logMode,
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		JWT_TTL:    jwtTTL,

		OPA_ADDRESS: opaAddress,

		DATABASE_NAME:     os.Getenv("DATABASE_NAME"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
		DATABASE_HOST:     os.Getenv("DATABASE_HOST"),
		DATABASE_SSL_MODE: os.Getenv("DATABASE_SSL_MODE"),
		DATABASE_PORT:     os.Getenv("DATABASE_PORT"),

		STORAGE_BUCKET_NAME:           os.Getenv("STORAGE_BUCKET_NAME"),
		STORAGE_URI:                   os.Getenv("STORAGE_URI"),
		STORAGE_ACCESS_KEY_ID:         os.Getenv("STORAGE_ACCESS_KEY_ID"),
		STORAGE_SECRECT_ACCESS_KEY:    os.Getenv("STORAGE_SECRECT_ACCESS_KEY"),
		STORAGE_USE_SSL:               useSSL,
		STORAGE_REGION:                os.Getenv("STORAGE_REGION"),
		STORAGE_SIGNED_URL_EXPIRATION: signedURLExpiration,
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
