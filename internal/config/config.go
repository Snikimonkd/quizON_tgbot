package config

import (
	"os"

	"quizon_bot/internal/logger"

	"gopkg.in/yaml.v3"
)

const (
	configPath = "values/local.yaml"
)

// GlobalConfig - глобальный конфиг
var GlobalConfig config

// Config - структура хранящая конфиг сервиса
type config struct {
	// Database - структурка для конфига базы
	Database struct {
		// DSN - строка для подключения к базе
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	// Keys - всякие ключи
	Keys struct {
		// Token - токен, который выдает телега
		Token string `yaml:"token"`
	} `yaml:"keys"`
}

func init() {
	file, err := os.Open(configPath)
	if err != nil {
		logger.Fatalf("can't open config file: %v", err)
	}
	defer func() {
		deferErr := file.Close()
		if err != nil {
			logger.Fatalf("can't close file: %v", deferErr)
		}
	}()

	d := yaml.NewDecoder(file)

	err = d.Decode(&GlobalConfig)
	if err != nil {
		logger.Fatalf("can't parse config file: %v", err)
	}
}
