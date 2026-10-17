package domain

import "errors"

var (
	ErrInvalidLanguage     = errors.New("invalid language selection")
	ErrLanguagePersistence = errors.New("failed to persist language selection")
)
