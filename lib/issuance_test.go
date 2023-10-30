//go:build integration
// +build integration

package coreumservicelib

import (
	"context"
	"testing"

	"coreumservice/go/stably_io/config"

	"github.com/stretchr/testify/require"
)

const (
	// senderMnemonic = "nut clog audit reward display era divide galaxy boil sport bless disorder total hidden pair range senior risk disorder affair dress barrel nuclear exhibit"
	// recipientMnemonic = "hazard misery record advice ceiling clean manage ten approve render abstract horse door federal congress stadium job tribe begin shaft digital aerobic upset record"
	recipientAddress = "testcore1un00l6nzdg58htj6e9fmx24433srcxpgdft57e"
)

func TestSuccessfulIssuanceUSDS(t *testing.T) {
	ctx := context.Background()

	usdsConfig := config.GetConfigDefault().Blockchain.Coreum.USDS
	senderSecretID := usdsConfig.TreasurySecretID
	assetDenom := usdsConfig.TokenDenom

	memo := "testing"
	var toAmount int64 = 1

	// Prepare the parameter for issuance
	transferParams, err := ProposeTransferStablyTokenParams(ctx,
		senderSecretID,
		recipientAddress,
		assetDenom,
		toAmount,
		memo,
	)
	require.NoError(t, err)
	require.NotEmpty(t, transferParams.GasPrice)
	require.NotZero(t, transferParams.GasUsed)
	require.NotZero(t, transferParams.SequenceNumber)

	// Submit the transaction on the blockchain
	cosmosTxResult, err := TransferStablyToken(ctx,
		senderSecretID,
		recipientAddress,
		assetDenom,
		toAmount,
		memo,
		transferParams.SequenceNumber,
		transferParams.GasPrice,
		transferParams.GasUsed,
	)
	require.NoError(t, err)
	require.NotEmpty(t, cosmosTxResult.TxHash)
}
