package config

import (
	"os"
	"strconv"
)

type Env struct {
	PORT       uint64
	JWT_SECRET string
}

func NewEnv() (*Env, error) {
	value, err := strconv.ParseUint(os.Getenv("PORT"), 10, 64)
	if err != nil {
		return nil, err
	}
	return &Env{
		PORT:       value,
		JWT_SECRET: os.Getenv("JWT_SECRET"),
	}, nil
}
