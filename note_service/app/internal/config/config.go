package config

import (
	"fmt"
	"note_service/app/pkg/logging"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	UserService struct {
		URL string `yaml:"url" env-required:"true"`
	} `yaml:"userservice" env-required:"true"`
	PostgreSQL struct {
		Host     string `yaml:"host" env-required:"true"`
		Port     string `yaml:"port" env-required:"true"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database" env-required:"true"`
	} `yaml:"postgresql" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	logger := logging.GetLogger()
	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	parentTop := filepath.Dir(filepath.Dir(wd))
	pathToConfig := fmt.Sprintf("%s/config.yml", parentTop)

	once.Do(func() {
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig(pathToConfig, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
