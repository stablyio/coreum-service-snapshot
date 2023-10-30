package awssecretmanager

const (
	KeyUsdsTreasuryWalletMnemonic = "usds_treasury_wallet_mnemonic"
)

type TokenizationSecrets struct {
	Coreum struct {
		USDsTreasuryWalletMnemonic string `json:"usds_treasury_wallet_mnemonic"`
	} `json:"coreum"`
}

func GetTokenizationSecrets() TokenizationSecrets {
	var tokenizationSecrets TokenizationSecrets
	unmarshalJSONSecretOrPanic("Tokenization", &tokenizationSecrets)
	return tokenizationSecrets
}
