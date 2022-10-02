package main

import (
	"context"
	"encoding/csv"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	ctx, cancel := context.WithCancel(context.Background())
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

	// Set event listener
	eventChan := make(chan os.Signal, 1)
	signal.Notify(eventChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Started parallel-uniqueness-checker")

	// Wait til notified
	select {
	case <-eventChan:
	case <-ctx.Done():
	}

	logger.Info("Stopping parallel-uniqueness-checker ...")

	cancel()

	// Wait for group to finish
	if err := group.Wait(); err != nil {
		logger.WithField("err", err).Fatal("group.Wait failed")
	}

}
