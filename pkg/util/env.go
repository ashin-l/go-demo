package util

import (
	"fmt"
	"os"
)

// Env reads specified environment variable. If no value has been found,
// fallback is returned.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		fmt.Println(v)
		return v
	}

	return fallback
}
