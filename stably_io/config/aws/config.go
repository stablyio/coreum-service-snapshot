package awsconfig

import "coreumservice/go/stably_io/utils"

type AwsConfig struct {
	SecretManager *SecretManager
}

type SecretManager struct {
	Region string
}

func GetConfig(_ utils.Stage) *AwsConfig {
	return &AwsConfig{
		SecretManager: &SecretManager{
			Region: "us-west-2",
		},
	}
}

func GetConfigDefault() *AwsConfig {
	return GetConfig(utils.GetStage())
}
