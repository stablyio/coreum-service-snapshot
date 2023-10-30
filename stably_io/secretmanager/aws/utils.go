package awssecretmanager

// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/

import (
	"context"
	awsconfig "coreumservice/go/stably_io/config/aws"
	"coreumservice/go/stably_io/utils"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"
)

//nolint:gochecknoglobals // For in-memory caching of secrets
var (
	initOnce = make(map[string]*sync.Once)
	secrets  = make(map[string]string)
)

func unmarshalJSONSecretOrPanic(secretName string, v interface{}) {
	secretJSON := getSecretStringOrPanic(secretName)
	err := json.Unmarshal([]byte(secretJSON), v)
	if err != nil {
		panic(errors.Wrapf(err, "failed to unmarshal JSON secret: %s", secretName))
	}
}

func getSecretStringOrPanic(secretName string) string {
	initSecretOnce, ok := initOnce[secretName]
	if !ok {
		initSecretOnce = &sync.Once{}
		initOnce[secretName] = initSecretOnce
	}
	initSecretOnce.Do(func() {
		ctx := context.Background()
		res, err := fetchSecretStringFromAWS(ctx, secretName)
		if err != nil {
			err = errors.Wrapf(err, "failed to fetch secret with name: %v", secretName)
			panic(utils.GetErrorWithStack(err))
		}
		secrets[secretName] = res
	})
	if secrets[secretName] == "" {
		panic(errors.New("secret is empty: " + secretName))
	}
	return secrets[secretName]
}

func fetchSecretStringFromAWS(ctx context.Context, secretName string) (string, error) {
	region := awsconfig.GetConfigDefault().SecretManager.Region
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", errors.Wrapf(err, "failed to load SDK configuration for region %s", region)
	}

	svc := secretsmanager.NewFromConfig(config)
	result, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	})
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return "", errors.Wrapf(err, "failed to get secret: %s", secretName)
	}
	if result == nil {
		return "", errors.New("result is nil for secret: " + secretName)
	}
	if result.SecretString == nil {
		return "", errors.New("secret string is nil for secret: " + secretName)
	}

	return *result.SecretString, nil
}
