package i18nerr

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	languageStore = sync.Map{} // map[uint64]I18NService
	i18nService   I18NService
)

// getGoroutineID 获取当前 goroutine ID
func getGoroutineID() uint64 {
	buf := make([]byte, 64)
	buf = buf[:runtime.Stack(buf, false)]
	idStr := strings.TrimPrefix(string(buf), "goroutine ")
	idStr = strings.Fields(idStr)[0]
	id, _ := strconv.ParseUint(idStr, 10, 64)
	return id
}

func SetI18NLanguage(lang string) {
	id := getGoroutineID()
	languageStore.Store(id, lang)
}

func GetI18NLanguage() string {
	id := getGoroutineID()
	if val, ok := languageStore.Load(id); ok {
		if lang, ok := val.(string); ok {
			return lang
		}
	}
	return ""
}

func ClearI18nLanguage() {
	id := getGoroutineID()
	languageStore.Delete(id)
}

func SetI18nService(service I18NService) {
	i18nService = service
}

func GetI18nService() I18NService {
	return i18nService
}
