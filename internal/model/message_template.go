package model

import "errors"

type MessageTemplate struct {
	ID         int64
	LanguageID int64
	Code       MessageTemplateCode
	Text       string
}

var ErrMessageTemplateNotFound = errors.New("message template not found")

type MessageTemplateCode string

const (
	TranslationPromptText     MessageTemplateCode = "TRANSLATION_PROMPT_TEXT"
	TranslationPlaceholder    MessageTemplateCode = "TRANSLATION_PLACEHOLDER"
	NotRegisteredMsg          MessageTemplateCode = "NOT_REGISTERED_MSG"
	SuccessfullyRegisteredMsg MessageTemplateCode = "SUCCESSFULLY_REGISTERED_MSG"
	AlreadyRegisteredMsg      MessageTemplateCode = "ALREADY_REGISTERED_MSG"
	LanguageSetMsg            MessageTemplateCode = "LANGUAGE_SET_MSG"
	UnsupportedLanguageMsg    MessageTemplateCode = "UNSUPPORTED_LANGUAGE_MSG"
	PreferredLanguageMsg      MessageTemplateCode = "PREFERRED_LANGUAGE_MSG"
	TranslationErrorMsg       MessageTemplateCode = "TRANSLATION_ERROR_MSG"
	UnexpectedErrMsg          MessageTemplateCode = "UNEXPECTED_ERR_MSG"
	NotGeorgianTextErrorMsg   MessageTemplateCode = "NOT_GEORGIAN_TEXT_ERROR_MSG"
	GeorgianTranslitMsg       MessageTemplateCode = "GEORGIAN_TRANSLIT_MSG"
)
