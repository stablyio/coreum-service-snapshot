package coreumconfig

import (
	configutils "coreumservice/go/stably_io/config/utils"
	"coreumservice/go/stably_io/utils"
)

const TestnetInitialTokenSupply = 9000000000000000
const TokenDecimal = 6
const TestnetRequiredNumberOfConfirmations = 1
const MainnetRequiredNumberOfConfirmations = 2
const HTTPServerPort = 5011

//nolint:gosec // This is the common value used in the test config
const TestUsdsTokenDenom = "microusds-testcore162rs3klx73exmyupxlqjju0u7aggcp0fswetn2"

type Coreum struct {
	RequiredNumberOfConfirmations int
	ClientServiceURL              string

	DepositEnabled bool
	USDS           CoreumAssetConfig
	PRCConfig      CoreumRPCConfig
}

type CoreumAssetConfig struct {
	TokenDenom         string
	TokenDecimal       int
	IssuanceEnabled    bool
	RedemptionEnabled  bool
	SupplyAdjustment   float64
	InitialTokenSupply uint64
	TreasurySecretID   string
}

type CoreumNetworkConfig struct {
	TokenTransferGasLimit         int
	ClientServiceURL              string
	RequiredNumberOfConfirmations int
	USDS                          Coreum
}

type CoreumRPCConfig struct {
	GRPCNodeURL          string
	TendermintRPCNodeURL string
	HTTPServerPort       int
}

func GetConfig(stage utils.Stage) *Coreum {
	return configutils.SwitchOnStage(stage,
		prod,
		beta,
		test,
		local,
	)
}

func GetConfigDefault() *Coreum {
	return GetConfig(utils.GetStage())
}
