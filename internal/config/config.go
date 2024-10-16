package config

import (
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port              string `yaml:"port" env-default:"9001"`
	DBName            string `yaml:"db_name" env-default:"gotth.db"`
	SessionCookieName string `yaml:"session_cookie_name" env-default:"session"`
}

func MustLoadConfig() *Config {
	configPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	configPath = strings.Join([]string{configPath, `\config\config.yaml`}, "")

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
