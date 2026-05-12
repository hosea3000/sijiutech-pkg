package log

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ginLoggerWriter 实现 io.Writer 接口，将 Gin 的日志重定向到 zap logger
type ginLoggerWriter struct {
	logger *zap.Logger
	level  zapcore.Level
}

func (w *ginLoggerWriter) Write(p []byte) (n int, err error) {
	// 移除末尾的换行符
	msg := strings.TrimRight(string(p), "\n")
	if w.logger.Core().Enabled(w.level) {
		w.logger.Log(w.level, msg)
	}
	return len(p), nil
}

// SetupGinLogger 配置 Gin 使用项目的 logger 输出日志
func (l *Logger) SetupGinLogger() {
	// 创建一个 info 级别的 writer 用于常规日志
	gin.DefaultWriter = &ginLoggerWriter{
		logger: l.Logger,
		level:  zapcore.InfoLevel,
	}
	// 创建一个 error 级别的 writer 用于错误日志
	gin.DefaultErrorWriter = &ginLoggerWriter{
		logger: l.Logger,
		level:  zapcore.ErrorLevel,
	}
}
