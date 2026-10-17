package db

import (
	"context"

	"tgbot/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn, err := cfg.PostgresURL()
	if err != nil {
		return nil, err
	}
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(context.Background(), poolCfg)
}
