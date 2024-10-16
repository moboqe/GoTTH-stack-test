package config

import (
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigTrinity struct {
	MySQL       `yaml:"mysql"`
	WorldServer `yaml:"world_server"`
}

type MySQL struct {
	User   string `yaml:"user" env-required:"true"`
	Passwd string `yaml:"password" env-required:"true"`
	Net    string `yaml:"net" env-required:"true"`
	Addr   string `yaml:"addr" env-required:"true"`
	DBName string `yaml:"db_name" env-required:"true"`
}

type WorldServer struct {
	WorldServerURL string `yaml:"world_serverurl" env-required:"true"`
	RootUsername   string `yaml:"root_username" env-required:"true"`
	RootPassword   string `yaml:"root_password" env-required:"true"`
}

func MustLoadTrinity() *ConfigTrinity {
	configPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	configPath = strings.Join([]string{configPath, `\config\config_trinity.yaml`}, "")

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg ConfigTrinity

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
