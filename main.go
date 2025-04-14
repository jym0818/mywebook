package main

import "github.com/spf13/viper"

func main() {
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
