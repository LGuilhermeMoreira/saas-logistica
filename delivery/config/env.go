package config

import (
	"fmt"
	"os"
	"time"
)

type Env struct {
	PORT        string
	LOG_MODE    string
	OPA_ADDRESS string
	JWT_SECRET  string
	JWT_TTL     time.Duration

	DATABASE_HOST     string
	DATABASE_PORT     string
	DATABASE_NAME     string
	DATABASE_PASSWORD string
	DATABASE_USER     string
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
		PORT:              os.Getenv("PORT"),
		LOG_MODE:          logMode,
		OPA_ADDRESS:       opaAddress,
		JWT_SECRET:        os.Getenv("JWT_SECRET"),
		DATABASE_HOST:     os.Getenv("DATABASE_HOST"),
		DATABASE_PORT:     os.Getenv("DATABASE_PORT"),
		DATABASE_NAME:     os.Getenv("DATABASE_NAME"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		JWT_TTL:           jwtTTL,
	}
}

func (e Env) MongoURI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s",
		e.DATABASE_USER,
		e.DATABASE_PASSWORD,
		e.DATABASE_HOST,
		e.DATABASE_PORT,
		e.DATABASE_NAME,
	)
}
