package telegram

import (
	"errors"
	"strings"

	"gopkg.in/telebot.v3"

	"tgbot/internal/core/domain"
	"tgbot/internal/core/port"
)

func showLanguageChooser(
	c telebot.Context,
	translator port.Translator,
	repository port.PersonRepository,
	supported []string,
	defaultLang string,
) error {
	lang := currentLanguage(c, repository, defaultLang)
	return c.Send(translator.T("select_language", lang), ChoosingLangInline(supported))
}

func currentLanguage(c telebot.Context, repository port.PersonRepository, defaultLang string) string {
	lang, err := repository.GetPersonLang(c.Sender().ID)
	if err != nil || lang == "" {
		return defaultLang
	}
	return lang
}

func onChooseLang(
	c telebot.Context,
	languages port.LanguageChooser,
	translator port.Translator,
	repository port.PersonRepository,
	defaultLang string,
) error {
	userID := c.Sender().ID
	current := currentLanguage(c, repository, defaultLang)
	parts := strings.Split(c.Callback().Data, ":")
	if len(parts) != 2 {
		return c.Respond(&telebot.CallbackResponse{Text: translator.T("callback_invalid_selection", current)})
	}
	result, err := languages.ChooseLanguage(userID, parts[1])
	if err != nil {
		key := "callback_unsupported_language"
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			key = "callback_update_language_failed"
		}
		return c.Respond(&telebot.CallbackResponse{Text: translator.T(key, current)})
	}
	if err := c.Send(translator.T(result.MessageKey, result.Language), UserMenu(translator, result.Language)); err != nil {
		return err
	}
	return c.Respond(&telebot.CallbackResponse{})
}
