package model

import (
	"errors"
)

type Language string

const (
	EnglishLanguage Language = "EN"
	RussianLanguage Language = "RU"
)

func ParseLang(lang string) (Language, error) {
	language := Language(lang)
	switch language {
	case EnglishLanguage, RussianLanguage:
		return language, nil
	default:
		return "", ErrUnsupportedLanguage
	}
}

var ErrUnsupportedLanguage = errors.New("unsupported language")
