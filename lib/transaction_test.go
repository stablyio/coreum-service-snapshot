//go:build integration
// +build integration

package coreumservicelib

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"coreumservice/go/stably_io/config"

	coreumconstant "github.com/CoreumFoundation/coreum/pkg/config/constant"
	"github.com/stretchr/testify/require"
)

func TestCalculateGasForTransfer(t *testing.T) {
	ctx := context.Background()
	coreumConfig := config.GetConfigDefault().Blockchain.Coreum
	usdsConfig := coreumConfig.USDS

	senderMnemonic := GetTreasuryMnemonicFromSecretID(ctx, usdsConfig.TreasurySecretID)

	keyringInfo, _, err := GetKeyringInfoFromMnemonic(senderMnemonic)
	require.NoError(t, err)

	t.Log("senderAddress", keyringInfo.GetAddress().String())

	accountInfo, err := GetAccountInfoFromMnemonic(ctx, senderMnemonic)
	require.NoError(t, err)

	gasUsed, gasPrice, err := CalculateGasForTransfer(ctx,
		senderMnemonic,
		"testcore1un00l6nzdg58htj6e9fmx24433srcxpgdft57e",
		usdsConfig.TokenDenom,
		2000000,
		"testing",
		accountInfo.Sequence,
	)
	require.NoError(t, err)

	t.Log("gasPrice", gasPrice)
	t.Log("gasUsed", gasUsed)

	// Remove the native token denom
	normalizeGasPriceString := strings.ReplaceAll(gasPrice, coreumconstant.DenomTest, "")

	gasPriceValue, err := strconv.ParseFloat(normalizeGasPriceString, 64)
	require.NoError(t, err)

	require.Greater(t, gasPriceValue, float64(0))
	require.Greater(t, gasUsed, uint64(0))
}
