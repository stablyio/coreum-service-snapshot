package config

import (
	awsconfig "coreumservice/go/stably_io/config/aws"
	blockchainconfig "coreumservice/go/stably_io/config/blockchain"

	"coreumservice/go/stably_io/utils"
)

type RunConfig struct {
	AWS        *awsconfig.AwsConfig
	Blockchain *blockchainconfig.Blockchain
}

// Get the run config based on the value of STAGE environment variable.
func GetConfigDefault() *RunConfig {
	return GetConfig(utils.GetStage())
}

// Get the run config with stage parameter can be `prod, beta, local, test,...`.
func GetConfig(stage utils.Stage) *RunConfig {
	return &RunConfig{
		AWS:        awsconfig.GetConfig(stage),
		Blockchain: blockchainconfig.GetConfig(stage),
	}
}
