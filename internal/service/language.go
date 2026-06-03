package service

import (
	"translator/internal/model"
	"translator/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LanguageService struct {
	repo repository.LanguageRepository
}

func NewLanguageService(db *pgxpool.Pool) *LanguageService {
	return &LanguageService{repo: repository.NewLanguageRepository(db)}
}

func (l *LanguageService) GetLangByCode(code model.LanguageCode) (model.Language, error) {
	return l.repo.GetByCode(code)
}

func (l *LanguageService) GetDefaultLang() (model.Language, error) {
	return l.repo.GetByCode(model.DefaultLangCode)
}
