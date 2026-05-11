package i18nerr

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func InjectLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		languages := parseAcceptLanguage(lang)
		SetI18NLanguage(languages[0])

		// 处理请求
		c.Next()

		// 请求结束后清理 context（使用 defer 确保即使发生 panic 也会清理）
		defer ClearI18nLanguage()
	}
}

// parseAcceptLanguage 解析Accept-Language头部
func parseAcceptLanguage(header string) []string {
	if header == "" {
		return []string{"zh-CN"}
	}

	// 按逗号分割语言项
	items := strings.Split(header, ",")
	languages := make([]string, 0, len(items))

	for _, item := range items {
		// 移除权重部分(q=0.x)
		lang := strings.Split(strings.TrimSpace(item), ";")[0]
		if lang != "" {
			languages = append(languages, lang)
		}
	}
	if len(languages) == 0 {
		return []string{"zh-CN"}
	}

	return languages
}
