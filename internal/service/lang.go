package service

import "translator/internal/model"

type LangTitle struct {
	Lang  model.Language
	Title string
}

func getLangTitles(lang model.Language) []LangTitle {
	return languageTitles[lang]
}

var languageTitles = map[model.Language][]LangTitle{
	model.EnglishLanguage: {
		{Lang: model.EnglishLanguage, Title: "English"},
		{Lang: model.RussianLanguage, Title: "Russian"},
	},
	model.RussianLanguage: {
		{Lang: model.EnglishLanguage, Title: "Английский"},
		{Lang: model.RussianLanguage, Title: "Русский"},
	},
}

var translationPlaceholders = map[model.Language]string{
	model.EnglishLanguage: "⏳Translating...",
	model.RussianLanguage: "⏳Перевожу...",
}
