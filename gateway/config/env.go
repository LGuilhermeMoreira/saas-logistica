package config

import (
	"os"
	"time"
)

type Env struct {
	PORT          string
	LOG_MODE      string
	AUTH_BASE_URL string
	OPA_ADDRESS   string
	JWT_SECRET    string
	JWT_TTL       time.Duration
}

func NewENV() *Env {

	jwtTTL, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		jwtTTL = time.Hour * 8
	}

	logMode := os.Getenv("LOG_MODE")
	if logMode == "" {
		logMode = "dev"
	}

	opaAddress := os.Getenv("OPA_ADDRESS")
	if opaAddress == "" {
		opaAddress = "http://localhost:8181/v1/data/authz/allow"
	}

	return &Env{
		PORT:          os.Getenv("PORT"),
		LOG_MODE:      logMode,
		AUTH_BASE_URL: os.Getenv("AUTH_BASE_URL"),
		OPA_ADDRESS:   opaAddress,
		JWT_SECRET:    os.Getenv("JWT_SECRET"),
		JWT_TTL:       jwtTTL,
	}
}
