package config

import (
	"net/url"
	"testing"
)

func TestPostgresURLEscapesCredentials(t *testing.T) {
	cfg := &Config{
		PostgresUser:     "user",
		PostgresPassword: "p@ss/#",
		PostgresHost:     "db.example",
		PostgresPort:     "5432",
		PostgresDB:       "bot",
		PostgresSSLMode:  "disable",
	}

	dsn, err := cfg.PostgresURL()
	if err != nil {
		t.Fatalf("PostgresURL() error = %v", err)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	user := parsed.User.Username()
	password, ok := parsed.User.Password()
	if user != "user" || !ok || password != "p@ss/#" {
		t.Fatalf("parsed credentials = %q/%q", user, password)
	}
	if parsed.Host != "db.example:5432" || parsed.Path != "/bot" {
		t.Fatalf("parsed address = %s%s", parsed.Host, parsed.Path)
	}
}
