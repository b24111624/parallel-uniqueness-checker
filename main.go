package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	logger "github.com/sirupsen/logrus"

	config "github.com/b24111624/parallel-uniqueness-checker/config"
)

func main() {
	// Init log
	logger.SetLevel(logger.DebugLevel)
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
	recordMap := map[string]string{}
	for _, f := range files {
		csvFile, err := os.Open(f.Name())
		if err != nil {
			logger.WithField("err", err).Fatal("os.Open failed")
		}
		defer csvFile.Close()

		csvReader := csv.NewReader(csvFile)

		// Get header
		title, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.WithField("err", err).Fatal(err)
		}

		err = checkTitle(title)
		if err != nil {
			logger.WithField("err", err).Fatal(err)
		}

		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.WithField("err", err).Fatal(err)
			}

			logger.Debug(line)

			if recordFile, ok := recordMap[line[1]]; ok {
				logger.WithFields(logger.Fields{
					"file": line[0] + " and " + recordFile,
					"code": line[1],
				}).Info("found duplication")
				return
			}

			recordMap[line[1]] = line[0]
		}
	}

	logger.Info("no dplication found!")
}

func checkTitle(title []string) error {
	if len(title) != 3 {
		return fmt.Errorf("invlid title, got %d column(s)", len(title))
	}

	if title[0] != "barcode" || title[1] != "code" || title[2] != "YearWeek" {
		return fmt.Errorf("invalid csv format")
	}

	return nil
}
