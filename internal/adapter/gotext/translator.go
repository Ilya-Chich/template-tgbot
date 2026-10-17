package gotext

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"tgbot/config"
	"tgbot/internal/core/port"

	"github.com/leonelquinteros/gotext"
)

type Translator struct {
	defaultLang string
	languages   []string
	once        sync.Once
	langs       map[string]*gotext.Locale
}

func NewTranslator(cfg *config.Config, languages []string) *Translator {
	return &Translator{defaultLang: cfg.Languages, languages: languages, langs: make(map[string]*gotext.Locale)}
}
func (t *Translator) T(key, lang string, args ...map[string]any) string {
	t.once.Do(t.loadLanguages)
	loc := t.langs[lang]
	if loc == nil {
		loc = t.langs[t.defaultLang]
	}
	if loc == nil {
		return key
	}
	translations, ok := loc.Domains[loc.GetDomain()]
	if !ok || translations == nil {
		return key
	}
	text := translations.Get(key)
	if len(args) > 0 {
		for name, value := range args[0] {
			text = replaceVar(text, name, fmt.Sprint(value))
		}
	}
	return text
}
func (t *Translator) loadLanguages() {
	path := filepath.Join(".", "locale")
	for _, code := range t.languages {
		locale := gotext.NewLocale(path, code)
		locale.AddDomain("bot")
		t.langs[code] = locale
	}
}
func replaceVar(s, name, value string) string {
	return strings.ReplaceAll(s, fmt.Sprintf("{%s}", name), value)
}

var _ port.Translator = (*Translator)(nil)
