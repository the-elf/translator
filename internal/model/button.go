package model

type Button struct {
	ID         int64
	LanguageID int64
	GroupName  ButtonGroupName
	Data       ButtonData
	Title      string
}

type ButtonGroupName string
type ButtonData string

const SetLangButtonGroup ButtonGroupName = "set_lang"
const GeoTranslitButtonGroup ButtonGroupName = "geo_translit"

const GeoTranslitButtonYes ButtonData = "yes"
const GeoTranslitButtonNo ButtonData = "no"
