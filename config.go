package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

var config = viper.New()

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02 15:04:05.000000") + " - " + string(bytes))
}

func initConfig() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	config.SetConfigFile("config.yml")

	err := config.ReadInConfig()

	if err != nil {
		log.Fatal("Fatal error config file")
	}
}
