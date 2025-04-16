package main

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	initLogger()
	initViper()
	r := InitWeb()
	r.Run(":8080")
}

func initViper() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("dev")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
