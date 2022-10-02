package config

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

const (
	defaultFileDir = "sameples"
)

type Config struct {
	FileDir string
}

func GetConfig() *Config {
	config := &Config{}

	fileDir, ok := os.LookupEnv("FILE_DIR")
	if !ok {
		logger.Warn("FILE_DIR not found, use default value")
		fileDir = defaultFileDir
	}

	config.FileDir = fileDir

	return config
}
