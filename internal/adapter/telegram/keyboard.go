package telegram

import (
	"tgbot/internal/core/port"

	"gopkg.in/telebot.v3"
)

var languageNames = map[string]string{"en": "English 🇬🇧🇺🇸", "ru": "Русский 🇷🇺", "rs": "Srpski 🇷🇸🇭🇷🇧🇦🇲🇪"}

func UserMenu(translator port.Translator, lang string) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{ResizeKeyboard: true}
	kb.Reply(kb.Row(kb.Text(translator.T("language_btn", lang))))
	return kb
}

func ChoosingLangInline(languages []string) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{}
	for _, code := range languages {
		text := languageNames[code]
		if text == "" {
			text = code
		}
		kb.InlineKeyboard = append(kb.InlineKeyboard, []telebot.InlineButton{{Text: text, Data: "choose_lang:" + code}})
	}
	return kb
}
