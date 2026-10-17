# Telegram bot template

Minimal telebot template: `/start` plus language selection.
Go + telebot.v3 + Postgres. Users table stores Telegram id and language.
## Setup

1. Copy `.env.example` to `.env` (production) and `.env.dev` (local), then set `TG_TOKEN`, `NAME`, `LANGUAGES` and `POSTGRES_*`.
2. Install gettext, which provides `msgfmt`: `apt install gettext` or `apk add gettext`.
3. Compile translations on the host: `./compile_translations.sh`. The `.mo` files are required at runtime.
4. `make init` (docker compose + postgres, dev).
5. `make migrate-up` reads `POSTGRES_*` from `.env.dev` and connects to `127.0.0.1:5666`. Do not set `POSTGRES_URL` for it, since `migrate-up` ignores it. `POSTGRES_URL` is an optional override for `migrate-down`, `migrate-status` and `migrate-version`.

## Production

- `docker-compose.prod.yml` bind-mounts `./locale` into the image, so run `./compile_translations.sh` on the host before `make init-prod` or `make up-prod`.
- `make up-prod` starts the stack. `make down-prod` stops it and keeps the Postgres volume.
- `make clean-volumes` deletes the Postgres volume. Use it only when you want to wipe the database.

Build: `go build ./cmd/api/main.go`

## Architecture

`cmd/api/main.go` is the composition root. It wires configuration and the database pool into the Postgres and gotext adapters, then injects those adapters into the core services and Telegram handlers. Core services depend only on interfaces in `internal/core/port`.

## Handlers

- `/start` — register user, welcome, language reply button
- language button — inline language list
- one `OnCallback` — `language.OnChooseLang`

Reply keyboard buttons are defined in `keyboard/`.
All `c.Send()` text must go through the injected translator.
