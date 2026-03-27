package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port string
	Env  string
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("error reading config:", err)
	}
	return Config{
		Port: viper.GetString("port"),
		Env:  viper.GetString("env"),
	}
}
