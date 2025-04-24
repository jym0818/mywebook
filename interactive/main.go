package main

import (
	"github.com/spf13/viper"
)

func main() {
	initViper()
	app := InitGRPCServer()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}

	err := app.server.Serve()
	if err != nil {
		panic(err)
	}

}

func initViper() {
	viper.SetConfigName("dev")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
