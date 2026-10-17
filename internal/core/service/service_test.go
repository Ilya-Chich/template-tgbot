package service

import (
	"errors"
	"testing"

	"tgbot/internal/core/domain"
)

type fakeRepository struct {
	person                    *domain.Person
	getErr, addErr, updateErr error
	added                     bool
	addedLang                 string
	addedLangTG               *string
	updatedLang               string
}

func (f *fakeRepository) AddNewPerson(_ int64, _, _, lang string, langTG *string) error {
	f.added = true
	f.addedLang = lang
	f.addedLangTG = langTG
	if f.addErr != nil {
		return f.addErr
	}
	f.person = &domain.Person{TgID: 7, Lang: &lang}
	f.getErr = nil
	return nil
}
func (f *fakeRepository) GetPerson(int64) (*domain.Person, error) { return f.person, f.getErr }
func (f *fakeRepository) GetPersonLang(int64) (string, error) {
	if f.person == nil {
		return "", f.getErr
	}
	if f.person.Lang == nil {
		return "", nil
	}
	return *f.person.Lang, nil
}

func (f *fakeRepository) UpdateLang(_ int64, lang string) error {
	f.updatedLang = lang
	return f.updateErr
}

func TestLanguageServiceChooseLanguage(t *testing.T) {
	r := &fakeRepository{}
	s := NewLanguageService(r, []string{"en", "ru"})
	result, err := s.ChooseLanguage(42, "ru")
	if err != nil || result.Language != "ru" || result.MessageKey != "inform_language" || r.updatedLang != "ru" {
		t.Fatalf("unexpected result: %+v, %v, repo=%+v", result, err, r)
	}
}

func TestLanguageServiceRejectsUnsupportedLanguage(t *testing.T) {
	r := &fakeRepository{}
	s := NewLanguageService(r, []string{"en"})
	_, err := s.ChooseLanguage(42, "rs")
	if !errors.Is(err, ErrInvalidLanguage) {
		t.Fatalf("error=%v", err)
	}
	if r.updatedLang != "" {
		t.Fatal("updated unsupported language")
	}
}

func TestLanguageServiceReturnsPersistenceFailure(t *testing.T) {
	r := &fakeRepository{updateErr: errors.New("database unavailable")}
	s := NewLanguageService(r, []string{"en"})
	_, err := s.ChooseLanguage(42, "en")
	if !errors.Is(err, ErrLanguagePersistence) {
		t.Fatalf("error=%v", err)
	}
}

func TestUserServiceStartsExistingUser(t *testing.T) {
	lang := "en"
	r := &fakeRepository{person: &domain.Person{TgID: 7, Lang: &lang}}
	users := NewUserService(r, "en")
	result, err := users.Start(7, "@user", "User", "en")
	if err != nil || r.added || result.Person.TgID != 7 || result.Language != "en" ||
		result.MessageKey != "main_message" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
}

func TestUserServiceCreatesMissingUser(t *testing.T) {
	r := &fakeRepository{getErr: errors.New("not found")}
	users := NewUserService(r, "ru")
	result, err := users.Start(7, "@user", "User", "ru-RU")
	if err != nil || !r.added || r.addedLang != "ru" || r.addedLangTG == nil || *r.addedLangTG != "ru-RU" ||
		result.Person.TgID != 7 ||
		result.Language != "ru" {
		t.Fatalf("unexpected result: %+v repo=%+v err=%v", result, r, err)
	}
}

func TestUserServiceReturnsCreateError(t *testing.T) {
	createErr := errors.New("insert failed")
	r := &fakeRepository{getErr: errors.New("not found"), addErr: createErr}
	users := NewUserService(r, "en")
	_, err := users.Start(7, "@user", "User", "en")
	if !errors.Is(err, createErr) {
		t.Fatalf("error=%v", err)
	}
}
