package model

type Button struct {
	ID         int64
	LanguageID int64
	GroupName  string
	Data       string
	Title      string
}

type ButtonGroupName string

const SetLangButtonGroup ButtonGroupName = "set_lang"
