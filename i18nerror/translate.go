package i18nerr

import (
	"context"
	"sijiutech-billing-api/pkg/log"

	"go.uber.org/zap"
)

func TranslateI18nKey(key string, logger *log.Logger) string {
	service := GetI18nService()
	lang := GetI18NLanguage()
	translatedValue, err := service.Translate(context.TODO(), lang, key)
	if err != nil {
		logger.Warn("Failed to translate:", zap.String("key", key), zap.Error(err))
		return key
	}
	return translatedValue
}
