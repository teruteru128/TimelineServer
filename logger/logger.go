package logger

import (
	"github.com/TinyKitten/Timeline/config"
	"go.uber.org/zap"
)

// GetLogger Zapロガーのインスタンスを取得
func GetLogger() *zap.Logger {
	debug := config.GetAPIConfig().Debug
	if debug {
		logger, _ := zap.NewDevelopment()
		return logger
	}
	logger, _ := zap.NewProduction()
	return logger
}
