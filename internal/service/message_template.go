package service

import (
	"translator/internal/model"
	"translator/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageTemplateService struct {
	repo repository.MessageTemplateRepository
}

func NewMessageTemplateService(db *pgxpool.Pool) *MessageTemplateService {
	return &MessageTemplateService{repo: repository.NewMessageTemplateRepository(db)}
}

func (m *MessageTemplateService) GetByCodeAndLang(code model.MessageTemplateCode, lang model.Language) (model.MessageTemplate, error) {
	return m.repo.GetByCodeAndLang(code, lang)
}
