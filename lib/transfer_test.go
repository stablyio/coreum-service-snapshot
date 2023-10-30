//go:build integration
// +build integration

package coreumservicelib_test

import (
	"context"
	"testing"

	lib "coreumservice/go/lib"

	"coreumservice/go/stably_io/config"

	"github.com/stretchr/testify/require"
)

func TestTransferTokenWithMnemonic(t *testing.T) {
	ctx := context.Background()

	amount := int64(2000000)
	recipientAddress := "testcore1un00l6nzdg58htj6e9fmx24433srcxpgdft57e"
	senderMnemonic := lib.GetTreasuryMnemonicFromSecretID(ctx, config.GetConfigDefault().Blockchain.Coreum.USDS.TreasurySecretID)
	memo := "TestTransferTokenWithMnemonic"

	runTransferTokenTest(ctx, t,
		senderMnemonic,
		recipientAddress,
		config.GetConfigDefault().Blockchain.Coreum.USDS.TokenDenom,
		amount,
		memo,
	)
}

func runTransferTokenTest(ctx context.Context,
	t *testing.T,
	senderMnemonic string,
	recipientAddress string,
	tokenDenom string,
	amount int64,
	memo string,
) {
	accountInfo, err := lib.GetAccountInfoFromMnemonic(ctx, senderMnemonic)
	require.NoError(t, err)

	gasUsed, gasPrice, err := lib.CalculateGasForTransfer(ctx,
		senderMnemonic,
		recipientAddress,
		tokenDenom,
		amount,
		memo,
		accountInfo.Sequence,
	)
	require.NoError(t, err)

	t.Log("gasUsed", gasUsed)
	t.Log("gasPrice", gasPrice)

	tx, err := lib.TransferTokenWithMnemonic(ctx,
		senderMnemonic,
		recipientAddress,
		tokenDenom,
		amount,
		memo,
		accountInfo.Sequence,
		gasPrice,
		gasUsed,
	)
	require.NoError(t, err)

	t.Log("tx.TxHash", tx.TxHash)
}
