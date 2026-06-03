package repository

import "translator/internal/model"

type UserRepository interface {
	Get(chatID int64) (model.User, error)
	Save(user model.User) error
}
