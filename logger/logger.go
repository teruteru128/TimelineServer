package logger

import (
	"github.com/TinyKitten/TimelineServer/config"
	"go.uber.org/zap"
)

// NewLogger Zapロガーのインスタンスを取得
func NewLogger() *zap.Logger {
	debug := config.GetAPIConfig().Debug
	if debug {
		logger, _ := zap.NewDevelopment()
		return logger
	}
	logger, _ := zap.NewProduction()
	return logger
}
