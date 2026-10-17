-- +goose Up

CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	tgid BIGINT UNIQUE,
	username TEXT,
	fullname TEXT,
	lang TEXT DEFAULT 'en',
	lang_tg TEXT,
	created_at TIMESTAMPTZ DEFAULT now()
);

-- +goose Down

DROP TABLE IF EXISTS users;
