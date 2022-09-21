package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

var ErrEnvironmentVariableIsEmpty = errors.New("environment variable is empty")

func Env(key string) (string, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return "", ErrEnvironmentVariableIsEmpty
	}

	return env, nil
}

func MustEnv(key string) string {
	env, err := Env(key)
	if err != nil {
		panic(err)
	}

	return env
}

func EnvInt64(key string) (int64, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, ErrEnvironmentVariableIsEmpty
	}

	value, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return value, nil
}

func MustEnvInt64(key string) int64 {
	env, err := EnvInt64(key)
	if err != nil {
		panic(err)
	}

	return env
}

func EnvFloat64(key string) (float64, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, ErrEnvironmentVariableIsEmpty
	}

	value, err := strconv.ParseFloat(env, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseFloat: %w", err)
	}

	return value, nil
}

func MustEnvFloat64(key string) float64 {
	env, err := EnvFloat64(key)
	if err != nil {
		panic(err)
	}

	return env
}

func EnvSecond(key string) (time.Duration, error) {
	env, err := EnvInt64(key)
	if err != nil {
		return 0, fmt.Errorf("EnvInt64: %w", err)
	}

	return time.Duration(env) * time.Second, nil
}

func MustEnvSecond(key string) time.Duration {
	env, err := EnvSecond(key)
	if err != nil {
		panic(err)
	}

	return env
}
