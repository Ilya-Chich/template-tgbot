package config

import (
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	TgToken          string
	Languages        string
	Name             string
	Debug            bool
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresSSLMode  string
}

func Load() *Config {
	getEnv := func(key string, required bool) string {
		val := os.Getenv(key)
		if required && val == "" {
			log.Fatalf("Missing required env: %s", key)
		}
		return val
	}

	return &Config{
		TgToken:          getEnv("TG_TOKEN", true),
		Languages:        getEnv("LANGUAGES", true),
		Name:             getEnv("NAME", true),
		Debug:            parseBool(getEnv("DEBUG", false)),
		PostgresDB:       getEnv("POSTGRES_DB", true),
		PostgresUser:     getEnv("POSTGRES_USER", true),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", true),
		PostgresHost:     getEnv("POSTGRES_HOST", true),
		PostgresPort:     getEnv("POSTGRES_PORT", true),
		PostgresSSLMode:  sslMode(getEnv("POSTGRES_SSLMODE", false)),
	}
}

func sslMode(v string) string {
	if strings.TrimSpace(v) == "" {
		return "disable"
	}
	return v
}

func (cfg *Config) PostgresURL() (string, error) {
	if cfg.PostgresUser == "" || cfg.PostgresHost == "" || cfg.PostgresPort == "" || cfg.PostgresDB == "" {
		return "", errors.New("incomplete postgres config")
	}
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.PostgresUser, cfg.PostgresPassword),
		Host:   net.JoinHostPort(cfg.PostgresHost, cfg.PostgresPort),
		Path:   "/" + cfg.PostgresDB,
	}
	q := u.Query()
	q.Set("sslmode", sslMode(cfg.PostgresSSLMode))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func parseBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1":
		return true
	default:
		return false
	}
}
