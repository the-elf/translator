package service

import (
	"translator/internal/model"
	"translator/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ButtonService struct {
	repo repository.ButtonRepository
}

func NewButtonService(db *pgxpool.Pool) *ButtonService {
	return &ButtonService{repo: repository.NewButtonRepository(db)}
}

func (b *ButtonService) GetGroupByNameAndLang(name model.ButtonGroupName, lang model.Language) ([]model.Button, error) {
	return b.repo.GetGroupByNameAndLang(name, lang)
}
