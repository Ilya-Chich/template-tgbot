package port

import "tgbot/internal/core/domain"

type PersonRepository interface {
	AddNewPerson(tgID int64, username, fullname, lang string, langTG *string) error
	GetPerson(tgID int64) (*domain.Person, error)
	GetPersonLang(tgID int64) (string, error)
	UpdateLang(tgID int64, lang string) error
}

type Translator interface {
	T(key, lang string, args ...map[string]any) string
}

type Logger interface {
	Error(message string, err error, component string)
	Fatal(message string, err error, component string)
}

type UserStarter interface {
	Start(userID int64, username, fullname, langTG string) (*StartResult, error)
}

type StartResult struct {
	Person     *domain.Person
	Language   string
	MessageKey string
}

type LanguageChooser interface {
	ChooseLanguage(userID int64, selectedLang string) (*LanguageResult, error)
}

type LanguageResult struct {
	Language   string
	MessageKey string
}
