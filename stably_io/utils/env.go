package utils

import (
	"os"

	"github.com/pkg/errors"
)

// Get the environment variable.
// If panicEmpty=true, it will panic if the environment variable is not presented or empty.
func GetEnvVar(key string, panicEmpty bool) string {
	value := os.Getenv(key)
	if panicEmpty && value == "" {
		panic(errors.Errorf("Missing %v environment variable", key))
	}
	return value
}
