package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestGetConfig() {
	type test struct {
		desc           string
		envs           map[string]string
		expectedConfig *Config
	}
	tests := []test{
		{
			desc: "base case",
			envs: map[string]string{
				"FILE_DIR": "FILE_DIR",
			},
			expectedConfig: &Config{
				FileDir: "FILE_DIR",
			},
		},
		{
			desc: "set to default value",
			envs: map[string]string{},
			expectedConfig: &Config{
				FileDir: defaultFileDir,
			},
		},
	}

	for _, t := range tests {
		s.Run(t.desc, func() {
			// Set up envs
			for envName, envValue := range t.envs {
				s.T().Setenv(envName, envValue)
			}

			// Call GetConfig
			config := GetConfig()
			s.Equal(t.expectedConfig, config)
		})
	}
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
