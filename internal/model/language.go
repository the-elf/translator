package model

import (
	"errors"
)

type Language struct {
	ID   int64
	Code LanguageCode
}

func EmptyLang() Language {
	return Language{}
}

type LanguageCode string

const DefaultLangCode LanguageCode = "ru"

var ErrUnsupportedLanguage = errors.New("unsupported language")
