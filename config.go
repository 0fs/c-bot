package main

import (
	"github.com/spf13/viper"
	"log"
)

var config = viper.New()

func initConfig() {

	config.SetConfigFile("config.yml")

	err := config.ReadInConfig()

	if err != nil {
		log.Fatal("Fatal error config file")
	}
}
