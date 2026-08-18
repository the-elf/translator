package util

import (
	"errors"
	"fmt"
	"os"
)

func RequireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("required env variable %s is not set", key)
	}

	return val, nil
}

func RequireEnvs(keys ...string) (map[string]string, error) {
	result := make(map[string]string, len(keys))
	var errs []error

	for _, key := range keys {
		val := os.Getenv(key)
		if val == "" {
			errs = append(errs, fmt.Errorf("required env variable %s is not set", key))
		}
		result[key] = val
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return result, nil
}
