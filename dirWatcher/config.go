package main

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Directory       string
	TimeInterval    string
	MagicString     string
	DatabasePath    string
	APIPort         int
}

func loadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.SetDefault("Directory", "D:/folder")
	viper.SetDefault("TimeInterval", 5*time.Minute)
	viper.SetDefault("MagicString", "There you are!!")
	viper.SetDefault("DatabasePath", "D:/folder/dirwatcher.db")
	viper.SetDefault("APIPort", 8080)

	err := viper.ReadInConfig()
	if err != nil {
		panic("Error reading config file")
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic("Unable to unmarshal config")
	}

	return &config
}
