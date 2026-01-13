package config

import (
	"fmt"
	"os"
)

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

// GetPort ambil port dari env, default 8080
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf(":%s", port)
}