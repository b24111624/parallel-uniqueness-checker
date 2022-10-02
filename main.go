package main

import (
	"os"

	logger "github.com/sirupsen/logrus"
)

func main() {
	logger.SetLevel(logger.InfoLevel)
	logger.Info("Start parallel-uniqueness-checker")

	dir, err := os.Open("samples")
	if err != nil {
		logger.Fatal("err:", err)
	}
	defer dir.Close()

	files, err := dir.ReadDir(0)
	if err != nil {
		logger.Fatal("err:", err)
	}

	for _, f := range files {
		logger.Info("fileName:", f.Name())
	}

}
