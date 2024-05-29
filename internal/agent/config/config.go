package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	envErrorStr = "The COMPUTING_POWER environment variable is not set or has an incorrect value."
)

type Config struct {
	ComputingPower int
}

func NewConfigFromEnv() (*Config, error) {
	cp, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || cp <= 0 {
		return nil, fmt.Errorf(envErrorStr)
	}
	return &Config{ComputingPower: cp}, nil
}
