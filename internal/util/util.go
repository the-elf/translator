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

func RequireEnvs(keys ...string) (map[string]string, map[string]error) {
	result := make(map[string]string)
	errs := make(map[string]error)

	for _, key := range keys {
		val := os.Getenv(key)
		if val == "" {
			errs[key] = fmt.Errorf("required env variable %s is not set", key)
		}
		result[key] = val
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return result, nil
}
