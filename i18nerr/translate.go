package i18nerr

import (
	"context"

	"go.uber.org/zap"
)

func TranslateI18nKey(key string, logger *zap.Logger) string {
	service := GetI18nService()
	lang := GetI18NLanguage()
	translatedValue, err := service.Translate(context.TODO(), lang, key)
	if err != nil {
		logger.Warn("Failed to translate:", zap.String("key", key), zap.Error(err))
		return key
	}
	return translatedValue
}
