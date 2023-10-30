package coreumservicelib

import (
	"context"
	"fmt"

	awssecretmanager "coreumservice/go/stably_io/secretmanager/aws"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

var coreumTreasuryPrivateMnemonicCache = map[string]string{}

// Function to retrieve the mnemonic of the treasury wallet
func GetTreasuryMnemonicFromSecretID(ctx context.Context, secretID string) string {
	secretValue := coreumTreasuryPrivateMnemonicCache[secretID]
	if secretValue != "" {
		return secretValue
	}

	tokenizationSecrets := awssecretmanager.GetTokenizationSecrets()
	coreumConfig := tokenizationSecrets.Coreum
	treasuryMnemonicString := ""
	switch secretID {
	case awssecretmanager.KeyUsdsTreasuryWalletMnemonic:
		treasuryMnemonicString = coreumConfig.USDsTreasuryWalletMnemonic
	}

	if treasuryMnemonicString == "" {
		fmt.Printf("missing treasuryMnemonicString for secret ID %v", secretID)
		return ""
	}

	// Cache the fetched value
	coreumTreasuryPrivateMnemonicCache[secretID] = treasuryMnemonicString

	return treasuryMnemonicString
}

func GetKeyringInfoFromMnemonic(mnemonic string) (keyring.Info, keyring.Keyring, error) {
	keyringInMemory := keyring.NewInMemory()
	// Generate private key and add it to the keystore
	keyringInfo, err := keyringInMemory.NewAccount(
		"",
		mnemonic,
		"",
		cosmossdk.GetConfig().GetFullBIP44Path(),
		hd.Secp256k1,
	)
	if err != nil {
		return nil, nil, errors.Errorf("keyringInMemory.NewAccount: %v", err)
	}
	return keyringInfo, keyringInMemory, nil
}
