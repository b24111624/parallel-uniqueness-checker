package checker

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	mockCTX = context.Background()
)

type CheckerTestSuite struct {
	suite.Suite
}

func (s *CheckerTestSuite) SetupSuite() {
}

func (s *CheckerTestSuite) SetupTest() {
}

func (s *CheckerTestSuite) TestRun() {
	type test struct {
		desc     string
		data     [][]string
		hasError bool
		error    error
	}
	tests := []test{
		{
			desc: "base case",
			data: [][]string{
				{
					"barcode", "code", "YearWeek",
				},
				{
					"TestProcessEligibleChannel2_1", "XJCQ", "TestProcessEligibleChannel2_1",
				},
				{
					"TestProcessEligibleChannel2_1", "PD2P", "TestProcessEligibleChannel2_1",
				},
				{
					"TestProcessEligibleChannel2_1", "FG2F", "TestProcessEligibleChannel2_1",
				},
			},
			hasError: false,
			error:    nil,
		},
		{
			desc: "error duplicate",
			data: [][]string{
				{
					"barcode", "code", "YearWeek",
				},
				{
					"TestProcessEligibleChannel2_1", "XJCQ", "TestProcessEligibleChannel2_1",
				},
				{
					"TestProcessEligibleChannel2_1", "XJCQ", "TestProcessEligibleChannel2_1",
				},
				{
					"TestProcessEligibleChannel2_1", "FG2F", "TestProcessEligibleChannel2_1",
				},
			},
			hasError: true,
			error:    fmt.Errorf("code %s is dupliced", "XJCQ"),
		},
		{
			desc: "error invalid title",
			data: [][]string{
				{
					"barcode", "code", "YearWeek", "invalid", "title",
				},
				{
					"TestProcessEligibleChannel2_1", "XJCQ", "TestProcessEligibleChannel2_1",
				},
				{
					"TestProcessEligibleChannel2_1", "PD2P", "TestProcessEligibleChannel2_1",
				},
				{
					"TestProcessEligibleChannel2_1", "FG2F", "TestProcessEligibleChannel2_1",
				},
			},
			hasError: true,
			error:    fmt.Errorf("invlid title, got %d column(s)", 5),
		},
		{
			desc: "error invalid title",
			data: [][]string{
				{
					"barcode", "invalid", "title",
				},
				{
					"TestProcessEligibleChannel2_1", "XJCQ", "TestProcessEligibleChannel2_1",
				},
				{
					"TestProcessEligibleChannel2_1", "PD2P", "TestProcessEligibleChannel2_1",
				},
				{
					"TestProcessEligibleChannel2_1", "FG2F", "TestProcessEligibleChannel2_1",
				},
			},
			hasError: true,
			error:    fmt.Errorf("invalid csv format"),
		},
	}
	for _, t := range tests {
		// init store
		checkerStore := NewStore(t.data, &sync.Map{})

		err := checkerStore.Run(mockCTX)()
		if t.hasError {
			s.Error(err, t.desc)
			s.Equal(t.error, err)
		} else {
			s.NoError(err)
		}
	}
}

func TestCheckerTestSute(t *testing.T) {
	suite.Run(t, new(CheckerTestSuite))
}
