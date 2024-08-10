package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type DBConfig struct {
	DB               string `yaml:"db"`
	Couponcollection string `yaml:"couponcollection"`
}

type Config struct {
	Database DBConfig `yaml:"database"`
}

func LoadConfig() Config {
	viper.AddConfigPath("./configs/")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Unable to load config", err)
		return Config{}
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Println("Unable to load config", err)
		return Config{}
	}

	fmt.Printf("Configuration: [%+v]\n", cfg)
	return cfg
}
