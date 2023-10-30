package coreumservicelib

import (
	"context"

	"github.com/CoreumFoundation/coreum/pkg/client"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

/*
Transfer the smart token to the recipientAddress, given the asset ID (denom) and the sender mnemonic
*/
func TransferTokenWithMnemonic(ctx context.Context,
	senderMnemonic string,
	recipientAddress string,
	assetDenom string,
	toAmount int64,
	memo string,
	sequenceNumber uint64,
	gasPrice string,
	gasUsed uint64,
) (*cosmossdk.TxResponse, error) {
	senderInfo, signingKeyRing, err := GetKeyringInfoFromMnemonic(senderMnemonic)
	if err != nil {
		return nil, errors.Errorf("GetKeyringInfoFromMnemonic: %v", err)
	}

	// Retrieve the sender address
	fromAddressStr := senderInfo.GetAddress().String()

	clientCtx, txFactory, msg, err := PrepareTransferTransaction(ctx,
		signingKeyRing,
		fromAddressStr,
		recipientAddress,
		assetDenom,
		toAmount,
		memo,
		sequenceNumber,
		gasPrice,
		gasUsed,
	)
	if err != nil {
		return nil, errors.Errorf("PrepareTransferTransaction: %v", err)
	}

	cosmosTxResult, err := client.BroadcastTx(ctx, clientCtx, txFactory, msg)
	if err != nil {
		return nil, errors.Errorf("client.BroadcastTx: %v", err)
	}
	return cosmosTxResult, nil
}
