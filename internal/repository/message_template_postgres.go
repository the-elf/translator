package repository

import (
	"context"
	"errors"
	"translator/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageTemplateRepository struct {
	db *pgxpool.Pool
}

func NewMessageTemplateRepository(db *pgxpool.Pool) MessageTemplateRepository {
	return &messageTemplateRepository{db: db}
}

func (m *messageTemplateRepository) GetByCodeAndLang(code model.MessageTemplateCode, lang model.Language) (model.MessageTemplate, error) {
	query := `select id, language_id, code, text from message_template where code = $1 and language_id = $2`

	var template model.MessageTemplate
	err := m.db.QueryRow(context.Background(), query, code, lang.ID).Scan(&template.ID, &template.LanguageId, &template.Code, &template.Text)
	if errors.Is(err, pgx.ErrNoRows) {
		return template, model.ErrMessageTemplateNotFound
	}

	return template, err
}
