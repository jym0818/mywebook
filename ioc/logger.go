package ioc

import (
	"github.com/jym/mywebook/pkg/logger"
	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(log)
}
