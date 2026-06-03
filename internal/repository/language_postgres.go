package repository

import (
	"context"
	"errors"
	"translator/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type languageRepository struct {
	db *pgxpool.Pool
}

func NewLanguageRepository(db *pgxpool.Pool) LanguageRepository {
	return &languageRepository{db: db}
}

func (l *languageRepository) GetByCode(code model.LanguageCode) (model.Language, error) {
	query := `select id, code from language where code = $1`

	var lang model.Language
	err := l.db.QueryRow(context.Background(), query, code).Scan(&lang.ID, &lang.Code)
	if errors.Is(err, pgx.ErrNoRows) {
		return lang, model.ErrUnsupportedLanguage
	}

	return lang, err
}
