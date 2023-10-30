package coreumconfig

func test() *Coreum {
	return &Coreum{
		RequiredNumberOfConfirmations: TestnetRequiredNumberOfConfirmations,
		ClientServiceURL:              "http://localhost:5011/coreumservice",
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
