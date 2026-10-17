package telegram

import (
	"gopkg.in/telebot.v3"
)

// HTMLParseMiddleware automatically adds HTML parse mode to all text messages
func HTMLParseMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		wrapper := &htmlParseContext{Context: c}
		return next(wrapper)
	}
}

type htmlParseContext struct {
	telebot.Context
}

func (c *htmlParseContext) Send(what any, opts ...any) error {
	if text, ok := what.(string); ok {
		sendOpts := &telebot.SendOptions{ParseMode: telebot.ModeHTML}
		for _, opt := range opts {
			if optOpts, ok := opt.(*telebot.SendOptions); ok {
				if optOpts.ParseMode == "" {
					optOpts.ParseMode = telebot.ModeHTML
				}
				sendOpts = optOpts
			} else if replyMarkup, ok := opt.(*telebot.ReplyMarkup); ok {
				sendOpts.ReplyMarkup = replyMarkup
			}
		}

		return c.Context.Send(text, sendOpts)
	}
	if photo, ok := what.(*telebot.Photo); ok && photo.Caption != "" {
		sendOpts := &telebot.SendOptions{ParseMode: telebot.ModeHTML}
		for _, opt := range opts {
			if optOpts, ok := opt.(*telebot.SendOptions); ok {
				if optOpts.ParseMode == "" {
					optOpts.ParseMode = telebot.ModeHTML
				}
				sendOpts = optOpts
			} else if replyMarkup, ok := opt.(*telebot.ReplyMarkup); ok {
				sendOpts.ReplyMarkup = replyMarkup
			}
		}

		return c.Context.Send(photo, sendOpts)
	}
	return c.Context.Send(what, opts...)
}
