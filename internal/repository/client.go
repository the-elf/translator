package repository

import (
	"errors"
	"translator/internal/model"
)

type ClientRepository interface {
	SaveIfAbsent(id int64) bool
	GetLang(id int64) (model.Language, error)
	SetLang(id int64, lang model.Language) error
}

var ErrClientNotFound = errors.New("client not found")
