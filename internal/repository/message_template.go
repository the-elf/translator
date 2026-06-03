package repository

import "translator/internal/model"

type MessageTemplateRepository interface {
	GetByCodeAndLang(code model.MessageTemplateCode, lang model.Language) (model.MessageTemplate, error)
}
