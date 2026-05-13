package i18ntype

import (
	"context"
	"encoding/json"
	"log"

	i18nerr "github.com/hosea3000/sijiutech-pkg/i18nerr"
)

type I18nString string

func (i I18nString) MarshalJSON() ([]byte, error) {
	translatedValue := string(i) // 默认值就是原始值，但是要转换一下类型，不然会循环调用

	service := i18nerr.GetI18nService()
	lang := i18nerr.GetI18NLanguage()
	translatedValue, err := service.Translate(context.TODO(), lang, translatedValue)
	if err != nil {
		log.Print("Failed to translate:", err)
		return json.Marshal(string(i))
	}

	// 返回翻译后的值
	return json.Marshal(translatedValue)
}
