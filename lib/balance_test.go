//go:build integration
// +build integration

package coreumservicelib_test

import (
	"context"
	lib "coreumservice/go/lib"
	"strconv"
	"testing"

	"coreumservice/go/stably_io/config"

	"github.com/stretchr/testify/require"
)

func TestGetBalanceOfAddress(t *testing.T) {
	ctx := context.Background()

	usdsConfig := config.GetConfigDefault().Blockchain.Coreum.USDS
	senderSecretID := usdsConfig.TreasurySecretID
	senderMnemonic := lib.GetTreasuryMnemonicFromSecretID(ctx, senderSecretID)
	t.Log("senderMnemonic", senderMnemonic)

	keyringInfo, _, err := lib.GetKeyringInfoFromMnemonic(senderMnemonic)
	require.NoError(t, err)

	senderAddress := keyringInfo.GetAddress().String()
	t.Log("senderAddress", senderAddress)

	tokenDenom := usdsConfig.TokenDenom
	balance, err := lib.GetBalanceOfAddress(ctx, senderAddress, tokenDenom)
	require.NoError(t, err)

	balanceValue, err := strconv.ParseFloat(balance, 64)
	require.NoError(t, err)

	t.Log("balance", balance)
	require.Greater(t, balanceValue, float64(1000000)) // at least 1 USDS (1e6 microUSDS)
}
