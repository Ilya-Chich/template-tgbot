package domain

import "time"

type Person struct {
	ID        int64
	TgID      int64
	Username  *string
	Fullname  *string
	Lang      *string
	LangTG    *string
	CreatedAt time.Time
}
