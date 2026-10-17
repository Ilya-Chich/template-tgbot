package postgres

import (
	"context"

	"tgbot/internal/core/domain"
	"tgbot/internal/core/port"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonRepository struct{ pool *pgxpool.Pool }

func NewPersonRepository(pool *pgxpool.Pool) *PersonRepository { return &PersonRepository{pool: pool} }

func (r *PersonRepository) AddNewPerson(tgID int64, username, fullname, lang string, langTG *string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO users (tgid, username, fullname, lang, lang_tg) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (tgid) DO NOTHING`,
		tgID,
		username,
		fullname,
		lang,
		langTG,
	)
	return err
}

func (r *PersonRepository) GetPerson(tgID int64) (*domain.Person, error) {
	var person domain.Person
	err := r.pool.QueryRow(context.Background(), `SELECT id, tgid, username, fullname, lang, lang_tg, created_at FROM users WHERE tgid = $1 LIMIT 1`, tgID).
		Scan(&person.ID, &person.TgID, &person.Username, &person.Fullname, &person.Lang, &person.LangTG, &person.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *PersonRepository) GetPersonLang(tgID int64) (string, error) {
	person, err := r.GetPerson(tgID)
	if err != nil || person.Lang == nil {
		return "", err
	}
	return *person.Lang, nil
}

func (r *PersonRepository) UpdateLang(tgID int64, lang string) error {
	_, err := r.pool.Exec(context.Background(), `UPDATE users SET lang = $1 WHERE tgid = $2`, lang, tgID)
	return err
}

var _ port.PersonRepository = (*PersonRepository)(nil)
