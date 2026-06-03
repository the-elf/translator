package repository

import "translator/internal/model"

type ButtonRepository interface {
	GetGroupByNameAndLang(name model.ButtonGroupName, lang model.Language) ([]model.Button, error)
}
