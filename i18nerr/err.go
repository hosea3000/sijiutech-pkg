package i18nerr

import "github.com/pkg/errors"

type I18nTranslationErrorKey string

type I18nTranslationError struct {
	Code       int
	MessageKey I18nTranslationErrorKey
}

func (e I18nTranslationError) Error() string {
	return string(e.MessageKey)
}

func NewError(code I18nCode) error {
	err := I18nTranslationError{
		Code:       code.Code,
		MessageKey: code.MessageKey,
	}
	return errors.WithStack(err)
}

type I18nCode struct {
	Code       int
	MessageKey I18nTranslationErrorKey
}

func NewCode(code int, messageKey I18nTranslationErrorKey) I18nCode {
	return I18nCode{
		Code:       code,
		MessageKey: messageKey,
	}
}
