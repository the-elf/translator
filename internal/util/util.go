package util

import (
	"fmt"
	"os"
)

func Ternary[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}

	return ifFalse
}

func RequireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("required env variable %s is not set", key)
	}

	return val, nil
}
