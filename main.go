package main

import (
	"encoding/csv"
	"os"

	logger "github.com/sirupsen/logrus"

	config "github.com/b24111624/parallel-uniqueness-checker/config"
)

func main() {
	// Init log
	logger.SetLevel(logger.InfoLevel)
	logger.Info("Start parallel-uniqueness-checker")

	// Get config
	config := config.GetConfig()

	// Get sample data folder
	dir, err := os.Open(config.FileDir)
	if err != nil {
		logger.WithField("err", err).Fatal("os.Open failed")
	}
	defer dir.Close()

	// Get files under sample data folder
	files, err := dir.ReadDir(0)
	if err != nil {
		logger.WithField("err", err).Fatal("dir.ReadDir failed")
	}

	// Set work dir to sample data folder
	err = os.Chdir(dir.Name())
	if err != nil {
		logger.WithField("err", err).Fatal("os.Chdir failed")
	}

	// Parse smaple data files
	for _, f := range files {
		csvFile, err := os.Open(f.Name())
		if err != nil {
			logger.WithField("err", err).Fatal("os.Open failed")
		}

		csvReader := csv.NewReader(csvFile)

		data, err := csvReader.ReadAll()
		if err != nil {
			logger.WithField("err", err).Fatal("csvReader.ReadAll failed")
		}

		logger.Info(data)
	}

}
