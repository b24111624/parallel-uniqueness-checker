package checker

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CheckerTestSuite struct {
	suite.Suite
}

func (s *CheckerTestSuite) TestRun() {
	// TODO:
}

func TestCheckerTestSute(t *testing.T) {
	suite.Run(t, new(CheckerTestSuite))
}
