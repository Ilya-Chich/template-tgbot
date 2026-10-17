package telegram

import (
	"gopkg.in/telebot.v3"
	"tgbot/internal/core/port"
)

func Register(bot *telebot.Bot, users port.UserStarter, languages port.LanguageChooser, translator port.Translator, repository port.PersonRepository, logger port.Logger, supported []string, defaultLang string) {
	bot.Handle("/start", HTMLParseMiddleware(LoggingMiddleware(func(c telebot.Context) error { return onStart(c, users, translator, logger) })))
	bot.Handle(telebot.OnText, HTMLParseMiddleware(LoggingMiddleware(func(c telebot.Context) error { return handleOnText(c, translator, repository, supported, defaultLang) })))
	RegisterCallbackHandlers(bot, languages, translator, repository, logger, supported, defaultLang)
}

func handleOnText(c telebot.Context, translator port.Translator, repository port.PersonRepository, supported []string, defaultLang string) error {
	for _, lang := range supported {
		if c.Text() == translator.T("language_btn", lang) {
			return showLanguageChooser(c, translator, repository, supported, defaultLang)
		}
	}
	return nil
}
