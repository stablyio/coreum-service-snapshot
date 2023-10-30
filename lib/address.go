package coreumservicelib

import (
	"context"
	"coreumservice/go/stably_io/utils"
	"fmt"

	coreumconstant "github.com/CoreumFoundation/coreum/pkg/config/constant"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/pkg/errors"
)

type AddressInfo struct {
	Address        string
	DerivationPath string
}

func GetAddress(mnemonic string, index int64) (*AddressInfo, error) {
	derivationPath := GetDerivationPath(index)

	keyringInMemory := keyring.NewInMemory()
	accountInfo, err := keyringInMemory.NewAccount(
		"", // uuid
		mnemonic,
		"", // passphrase
		derivationPath,
		hd.Secp256k1,
	)
	if err != nil {
		return nil, errors.Errorf("clientCtx.Keyring().NewAccount: %v", err)
	}
	return &AddressInfo{
		Address:        accountInfo.GetAddress().String(),
		DerivationPath: derivationPath,
	}, nil
}

func GetTreasuryAddress(ctx context.Context, secretID string) (*AddressInfo, error) {
	mnemonic := GetTreasuryMnemonicFromSecretID(ctx, secretID)
	addressInfo, err := GetAddress(mnemonic, 0)
	if err != nil {
		return nil, errors.Errorf("GetAddress(mnemonic, 0): %v", err)
	}
	return addressInfo, nil
}

func GetDerivationPath(index int64) string {
	// Reference:
	// - Coreum section in https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	coinType := 990
	path := fmt.Sprintf("m/44'/%d'/%v'/0/0", coinType, index)
	return path
}

func GetAddressPrefixByStage() string {
	addressPrefix := ""
	switch utils.GetStage() {
	case utils.Prod:
		addressPrefix = coreumconstant.AddressPrefixMain
	case utils.Beta, utils.Local, utils.Test:
		addressPrefix = coreumconstant.AddressPrefixDev
	default:
		break
	}
	return addressPrefix
}
