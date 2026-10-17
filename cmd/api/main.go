package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"tgbot/config"
	"tgbot/database/db"
	"tgbot/internal/adapter/gotext"
	"tgbot/internal/adapter/logger"
	"tgbot/internal/adapter/postgres"
	"tgbot/internal/adapter/telegram"
	"tgbot/internal/core/service"

	"gopkg.in/telebot.v3"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New()
	pool, err := db.Init(cfg)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", err, "database")
	}
	defer pool.Close()
	languages := configuredLanguages(cfg.Languages)
	repository := postgres.NewPersonRepository(pool)
	translator := gotext.NewTranslator(cfg, languages)
	users := service.NewUserService(repository, cfg.Languages)
	languageService := service.NewLanguageService(repository, languages)
	bot, err := telebot.NewBot(
		telebot.Settings{
			Token:  cfg.TgToken,
			Poller: &telebot.LongPoller{Timeout: 10 * 60, AllowedUpdates: []string{"message", "callback_query"}},
		},
	)
	if err != nil {
		appLogger.Fatal("Failed to create bot", err, "bot")
	}
	telegram.Init(bot, users, languageService, translator, repository, appLogger, languages, cfg.Languages)
	log.Println("Starting Telegram bot...")
	go bot.Start()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
	bot.Stop()
	log.Println("Stopped")
}

func configuredLanguages(value string) []string {
	var result []string
	for _, language := range strings.Split(value, ",") {
		if language = strings.TrimSpace(language); language != "" {
			result = append(result, language)
		}
	}
	return result
}
