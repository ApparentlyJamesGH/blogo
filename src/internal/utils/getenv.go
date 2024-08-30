package utils

import "os"

func Getenv(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return def
}
