package repository

import "translator/internal/model"

type LanguageRepository interface {
	GetByCode(code model.LanguageCode) (model.Language, error)
}
