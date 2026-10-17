package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strings"
	"tgbot/internal/callback"
	"tgbot/internal/core/port"
)

func RegisterCallbackHandlers(bot *telebot.Bot, languages port.LanguageChooser, translator port.Translator, repository port.PersonRepository, logger port.Logger, supported []string, defaultLang string) {
	bot.Handle(telebot.OnCallback, HTMLParseMiddleware(LoggingMiddleware(func(c telebot.Context) error {
		data := cleanCallbackData(c.Callback().Data)
		if strings.HasPrefix(data, callback.ChooseLang) {
			return onChooseLang(c, languages, translator, repository, supported, defaultLang)
		}
		logger.Error("Unknown callback action", fmt.Errorf("action: %s", data), "callback_router")
		return c.Respond(&telebot.CallbackResponse{Text: translator.T("callback_unknown_action", currentLanguage(c, repository, defaultLang))})
	})))
}
func cleanCallbackData(data string) string { return strings.TrimPrefix(data, "\f") }
