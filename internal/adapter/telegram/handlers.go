package telegram

import (
	"tgbot/internal/core/port"

	"gopkg.in/telebot.v3"
)

func Init(
	bot *telebot.Bot,
	users port.UserStarter,
	languages port.LanguageChooser,
	translator port.Translator,
	repository port.PersonRepository,
	logger port.Logger,
	supported []string,
	defaultLang string,
) {
	Register(bot, users, languages, translator, repository, logger, supported, defaultLang)
}
