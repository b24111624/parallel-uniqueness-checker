package checker

import (
	"context"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
)

type Checker struct {
	Data      [][]string
	RecordMap *sync.Map
}

func NewStore(data [][]string, recordMap *sync.Map) Store {
	return &Checker{
		Data:      data,
		RecordMap: recordMap,
	}
}

func (c *Checker) Run(ctx context.Context) func() error {
	return func() error {
		if err := checkTitle(c.Data[0]); err != nil {
			logger.WithField("err", err).Error("invalid title")
			return err
		}

		for i := 1; i < len(c.Data); i++ {
			line := c.Data[i]
			logger.WithField("line", line).Debug()

			if _, ok := c.RecordMap.Load(line[1]); ok {
				logger.WithFields(logger.Fields{
					"code": line[1],
				}).Info("found duplication")
				return fmt.Errorf("code is dupliced")
			}

			c.RecordMap.Store(line[1], true)
		}

		return nil
	}
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
