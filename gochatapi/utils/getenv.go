package utils

import (
	"os"
)

// GetEnv ...
func GetEnv(key, def string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return def
	}
	return value
}
