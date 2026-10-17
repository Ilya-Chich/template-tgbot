package main

import (
	"fmt"
	"os"

	"tgbot/config"
)

func main() {
	cfg := &config.Config{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
	dsn, err := cfg.PostgresURL()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err = os.Stdout.WriteString(dsn + "\n"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
