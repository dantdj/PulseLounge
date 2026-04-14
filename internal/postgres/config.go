package postgres

import (
	"fmt"
	"net/url"
	"os"
)

func BuildURLFromEnv() string {
	user := getenvOrDefault("POSTGRES_USER", "postgres")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := getenvOrDefault("POSTGRES_HOST", "localhost")
	port := getenvOrDefault("POSTGRES_PORT", "5432")
	name := getenvOrDefault("POSTGRES_DB", "pulselounge")
	sslMode := getenvOrDefault("POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		name,
		url.QueryEscape(sslMode),
	)
}

func getenvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
