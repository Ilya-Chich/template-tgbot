package telegram

import (
	"tgbot/internal/core/port"

	"gopkg.in/telebot.v3"
)

func onStart(c telebot.Context, users port.UserStarter, translator port.Translator, logger port.Logger) error {
	user := c.Sender()
	if user.IsBot {
		return nil
	}
	username := "@" + user.Username
	if user.Username == "" {
		username = "@None"
	}
	result, err := users.Start(user.ID, username, user.FirstName, user.LanguageCode)
	if err != nil {
		logger.Error("User not found", err, "mainmenu")
		return c.Send(translator.T("user_not_found_try_again", "en"))
	}
	return c.Send(translator.T(result.MessageKey, result.Language), UserMenu(translator, result.Language))
}
