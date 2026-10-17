package service

import (
	"fmt"

	"tgbot/internal/core/domain"
	"tgbot/internal/core/port"
)

var (
	ErrInvalidLanguage     = domain.ErrInvalidLanguage
	ErrLanguagePersistence = domain.ErrLanguagePersistence
)

type LanguageService struct {
	repository port.PersonRepository
	languages  map[string]struct{}
}

func NewLanguageService(repository port.PersonRepository, languages []string) *LanguageService {
	supported := make(map[string]struct{}, len(languages))
	for _, language := range languages {
		supported[language] = struct{}{}
	}
	return &LanguageService{repository: repository, languages: supported}
}

func (s *LanguageService) ChooseLanguage(userID int64, selectedLang string) (*port.LanguageResult, error) {
	if _, ok := s.languages[selectedLang]; !ok {
		return nil, ErrInvalidLanguage
	}
	if err := s.repository.UpdateLang(userID, selectedLang); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLanguagePersistence, err)
	}
	return &port.LanguageResult{Language: selectedLang, MessageKey: "inform_language"}, nil
}

var _ port.LanguageChooser = (*LanguageService)(nil)
