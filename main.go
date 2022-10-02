package main

import (
	"context"
	"encoding/csv"
	"os"
	"sync"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/b24111624/parallel-uniqueness-checker/checker"
	config "github.com/b24111624/parallel-uniqueness-checker/config"
)

func main() {
	// Init log
	logger.SetLevel(logger.DebugLevel)
	logger.Info("Starting parallel-uniqueness-checker ...")

	// Get config
	config := config.GetConfig()

	// Init Context
	ctx := context.Background()
	group, ctx := errgroup.WithContext(ctx)

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

	// Process smaple data files
	recordMap := sync.Map{}
	for _, f := range files {
		csvFile, err := os.Open(f.Name())
		if err != nil {
			logger.WithField("err", err).Fatal("os.Open failed")
		}
		defer csvFile.Close()

		// Get all data from file
		csvReader := csv.NewReader(csvFile)
		records, err := csvReader.ReadAll()
		if err != nil {
			logger.WithField("err", err).Fatal("csvReader.ReadAll failed")
		}

		// Start Checker Store
		checkerStore := checker.NewStore(records, &recordMap)
		group.Go(checkerStore.Run(ctx))
	}

	logger.Info("Started parallel-uniqueness-checker")

	// Wait for group to finish
	if err := group.Wait(); err != nil {
		logger.WithField("err", err).Fatal("group.Wait failed")
	}

	logger.Info("No duplicate code")
}
