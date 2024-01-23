package env

import "os"

// GetWithDefault returns the value of an environment variable,
// or the provided default if the environment was not set.
func GetWithDefault(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
