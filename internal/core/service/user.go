package service

import (
	"fmt"

	"tgbot/internal/core/port"
)

type UserService struct {
	repository  port.PersonRepository
	defaultLang string
}

func NewUserService(repository port.PersonRepository, defaultLang string) *UserService {
	return &UserService{repository: repository, defaultLang: defaultLang}
}

func (s *UserService) Start(userID int64, username, fullname, langTG string) (*port.StartResult, error) {
	language, err := s.repository.GetPersonLang(userID)
	if err != nil || language == "" {
		language = s.defaultLang
	}
	person, err := s.repository.GetPerson(userID)
	if err == nil {
		return &port.StartResult{Person: person, Language: language, MessageKey: "main_message"}, nil
	}
	if addErr := s.repository.AddNewPerson(userID, username, fullname, s.defaultLang, &langTG); addErr != nil {
		return nil, fmt.Errorf("add user: %w", addErr)
	}
	person, err = s.repository.GetPerson(userID)
	if err != nil {
		return nil, fmt.Errorf("get user after create: %w", err)
	}
	return &port.StartResult{Person: person, Language: s.defaultLang, MessageKey: "main_message"}, nil
}

var _ port.UserStarter = (*UserService)(nil)

// StartResult is retained as an alias for callers migrating to the inbound port.
type StartResult = port.StartResult
