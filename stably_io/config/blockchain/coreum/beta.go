package coreumconfig

func beta() *Coreum {
	clientServiceURL := "http://internal-stably-internal-lb-285640036.us-west-2.elb.amazonaws.com/coreumservice"
	return &Coreum{
		RequiredNumberOfConfirmations: TestnetRequiredNumberOfConfirmations,
		ClientServiceURL:              clientServiceURL,
		DepositEnabled:                true,
		USDS: CoreumAssetConfig{
			TokenDecimal:       TokenDecimal,
			TokenDenom:         TestUsdsTokenDenom,
			IssuanceEnabled:    true,
			RedemptionEnabled:  true,
			SupplyAdjustment:   0.0,
			InitialTokenSupply: TestnetInitialTokenSupply,
			TreasurySecretID:   "usds_treasury_wallet_mnemonic",
		},
		PRCConfig: CoreumRPCConfig{
			GRPCNodeURL:          "full-node.testnet-1.coreum.dev:9090",
			TendermintRPCNodeURL: "https://full-node.testnet-1.coreum.dev:26657",
			HTTPServerPort:       HTTPServerPort,
		},
	}
}
