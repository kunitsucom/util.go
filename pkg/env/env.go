package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

var ErrEnvironmentVariableIsEmpty = errors.New("env: environment variable is empty")

func String(key string) (string, error) {
	value, found := os.LookupEnv(key)
	if !found {
		return "", fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	return value, nil
}

func StringOrDefault(key string, defaultValue string) string {
	value, err := String(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustString(key string) string {
	value, err := String(key)
	if err != nil {
		panic(err)
	}

	return value
}

func Bool(key string) (bool, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return false, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseBool(env)
	if err != nil {
		return false, fmt.Errorf("strconv.ParseBool: %w", err)
	}

	return value, nil
}

func BoolOrDefault(key string, defaultValue bool) bool {
	value, err := Bool(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustBool(key string) bool {
	env, err := Bool(key)
	if err != nil {
		panic(err)
	}

	return env
}

func Int(key string) (int, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.Atoi(env)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi: %w", err)
	}

	return value, nil
}

func IntOrDefault(key string, defaultValue int) int {
	value, err := Int(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustInt(key string) int {
	env, err := Int(key)
	if err != nil {
		panic(err)
	}

	return env
}

func Int64(key string) (int64, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return value, nil
}

func Int64OrDefault(key string, defaultValue int64) int64 {
	value, err := Int64(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustInt64(key string) int64 {
	env, err := Int64(key)
	if err != nil {
		panic(err)
	}

	return env
}

func Uint(key string) (uint, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseUint(env, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return uint(value), nil
}

func UintOrDefault(key string, defaultValue uint) uint {
	value, err := Uint(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustUint(key string) uint {
	env, err := Uint(key)
	if err != nil {
		panic(err)
	}

	return env
}

func Uint64(key string) (uint64, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseUint(env, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return value, nil
}

func Uint64OrDefault(key string, defaultValue uint64) uint64 {
	value, err := Uint64(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustUint64(key string) uint64 {
	env, err := Uint64(key)
	if err != nil {
		panic(err)
	}

	return env
}

func Float64(key string) (float64, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseFloat(env, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseFloat: %w", err)
	}

	return value, nil
}

func Float64OrDefault(key string, defaultValue float64) float64 {
	value, err := Float64(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustFloat64(key string) float64 {
	env, err := Float64(key)
	if err != nil {
		panic(err)
	}

	return env
}

func Second(key string) (time.Duration, error) {
	env, err := Int64(key)
	if err != nil {
		return 0, fmt.Errorf("Int64: %w", err)
	}

	return time.Duration(env) * time.Second, nil
}

func SecondOrDefault(key string, defaultValue time.Duration) time.Duration {
	value, err := Second(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustSecond(key string) time.Duration {
	env, err := Second(key)
	if err != nil {
		panic(err)
	}

	return env
}
