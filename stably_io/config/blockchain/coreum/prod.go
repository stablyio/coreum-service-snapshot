package coreumconfig

func prod() *Coreum {
	clientServiceURL := "http://internal-stably-internal-lb-87560538.us-west-2.elb.amazonaws.com/coreumservice"

	return &Coreum{
		RequiredNumberOfConfirmations: MainnetRequiredNumberOfConfirmations,
		ClientServiceURL:              clientServiceURL,
		DepositEnabled:                true,
		USDS: CoreumAssetConfig{
			TokenDecimal:       TokenDecimal,
			TokenDenom:         "microusds-core17z02cx2xxz2rehq6qay3rc06g5ksa9nxjwh5uv",
			IssuanceEnabled:    true,
			RedemptionEnabled:  true,
			SupplyAdjustment:   0.0,
			InitialTokenSupply: 10000000000000, //nolint:gomnd // 10M USDS, with 6 decimals
			TreasurySecretID:   "usds_treasury_wallet_mnemonic",
		},
		PRCConfig: CoreumRPCConfig{
			GRPCNodeURL:          "full-node.mainnet-1.coreum.dev:9090",
			TendermintRPCNodeURL: "https://full-node.mainnet-1.coreum.dev:26657",
			HTTPServerPort:       HTTPServerPort,
		},
	}
}
