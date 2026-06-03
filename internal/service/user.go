package service

import (
	"translator/internal/model"
	"translator/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{repo: repository.NewUserRepository(db)}
}

func (u *UserService) GetUser(chatID int64) (model.User, error) {
	return u.repo.Get(chatID)
}

func (u *UserService) SetLang(user model.User, lang model.Language) error {
	user.Language = lang
	return u.repo.Save(user)
}

func (u *UserService) SaveUser(user model.User) error {
	return u.repo.Save(user)
}
