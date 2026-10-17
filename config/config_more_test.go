package config

import "testing"

func TestParseBoolAndSSLMode(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want bool
	}{{"true", true}, {" TRUE ", true}, {"1", true}, {"yes", false}, {"", false}} {
		if got := parseBool(tc.in); got != tc.want {
			t.Errorf("parseBool(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
	if sslMode("  ") != "disable" || sslMode("require") != "require" {
		t.Fatal("sslMode normalization failed")
	}
}

func TestPostgresURLValidationAndDefaults(t *testing.T) {
	if _, err := (&Config{}).PostgresURL(); err == nil {
		t.Fatal("expected incomplete config error")
	}
	cfg := &Config{PostgresUser: "u", PostgresPassword: "p", PostgresHost: "localhost", PostgresPort: "5432", PostgresDB: "db"}
	dsn, err := cfg.PostgresURL()
	if err != nil {
		t.Fatal(err)
	}
	if dsn != "postgres://u:p@localhost:5432/db?sslmode=disable" {
		t.Errorf("dsn = %q", dsn)
	}
}
