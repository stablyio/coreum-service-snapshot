package coreumservicelib

import (
	"context"

	"github.com/CoreumFoundation/coreum/pkg/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
)

type TransferTokenParams struct {
	GasPrice       string
	GasUsed        uint64
	SequenceNumber uint64
}

// Propose the parameters that is used to generate the idempotent transaction
func proposeTransferTokenParams(ctx context.Context,
	senderMnemonic string,
	recipientAddress string,
	denom string,
	toAmount int64,
	memo string,
) (*TransferTokenParams, error) {
	senderInfo, senderKeyring, err := GetKeyringInfoFromMnemonic(senderMnemonic)
	if err != nil {
		return nil, errors.Errorf("GetKeyringInfoFromMnemonic: %v", err)
	}

	// Retrieve the sender address
	senderAddress := senderInfo.GetAddress().String()

	acc, err := GetAccountInfo(ctx, senderAddress)
	if err != nil {
		return nil, errors.Errorf("GetAccountInfo: %v", err)
	}
	currentSequenceNumber := acc.Sequence

	clientCtx, txFactory, msg, err := PrepareTransferTransaction(ctx,
		senderKeyring,
		senderAddress,
		recipientAddress,
		denom,
		toAmount,
		memo,
		currentSequenceNumber,
		"",
		0,
	)
	if err != nil {
		return nil, errors.Errorf("PrepareTransferTransaction: %v", err)
	}

	// calculate gas
	gasUsedForTransaction, gasPrice, err := CalculateGas(ctx, clientCtx, txFactory, msg)
	if err != nil {
		return nil, errors.Errorf("CalculateGas: %v", err)
	}

	transferParams := &TransferTokenParams{
		GasPrice:       gasPrice,
		GasUsed:        gasUsedForTransaction,
		SequenceNumber: currentSequenceNumber,
	}

	return transferParams, nil
}

func ProposeTransferStablyTokenParams(ctx context.Context,
	senderSecretID string,
	recipientAddress string,
	assetDenom string,
	toAmount int64,
	memo string,
) (*TransferTokenParams, error) {
	treasuryMnemonic := GetTreasuryMnemonicFromSecretID(ctx, senderSecretID)
	return proposeTransferTokenParams(ctx, treasuryMnemonic, recipientAddress, assetDenom, toAmount, memo)
}

func PrepareTransferTransaction(ctx context.Context,
	signingKeyRing keyring.Keyring,
	fromAddressStr string,
	recipientAddress string,
	denom string,
	toAmount int64,
	memo string,
	sequenceNumber uint64,
	gasPrice string,
	gasUsed uint64,
) (client.Context, client.Factory, cosmossdk.Msg, error) {
	fromAddress, err := cosmossdk.AccAddressFromBech32(fromAddressStr)
	if err != nil {
		return client.Context{}, client.Factory{}, nil, errors.Errorf("cosmossdk.AccAddressFromBech32: %v", err)
	}

	acc, err := GetAccountInfo(ctx, fromAddressStr)
	if err != nil {
		return client.Context{}, client.Factory{}, nil, errors.Errorf("GetAccountInfo: %v", err)
	}

	clientCtx := GetClientContext().
		// Assign the keyring that has the private key to sign the transaction
		WithKeyring(signingKeyRing).
		// From the specific address
		WithFromAddress(fromAddress)

	// Tx Factory contains the parameters that is used to generate the idempotent transaction
	txFactory := CoreumTxFactory(clientCtx).
		WithAccountNumber(acc.AccountNumber).
		WithSequence(sequenceNumber).
		WithGasPrices(gasPrice).
		WithGas(gasUsed).
		WithMemo(memo)

	// We send "toAmount" tokens from your wallet to recipientAddress
	msg := &banktypes.MsgSend{
		FromAddress: fromAddressStr,
		ToAddress:   recipientAddress,
		Amount:      cosmossdk.NewCoins(cosmossdk.NewInt64Coin(denom, toAmount)),
	}

	return clientCtx, txFactory, msg, nil
}

// The purpose of this function:
// - Determine the valid treasury wallet and asset denom by asset ticker
// - Try to submit the transfer transaction to the blockchain network
func TransferStablyToken(ctx context.Context,
	senderSecretID string,
	recipientAddress string,
	assetDenom string,
	toAmount int64,
	memo string,
	sequenceNumber uint64,
	gasPrice string,
	gasUsed uint64,
) (*cosmossdk.TxResponse, error) {
	treasuryMnemonic := GetTreasuryMnemonicFromSecretID(ctx, senderSecretID)
	cosmosTxResult, err := TransferTokenWithMnemonic(ctx,
		treasuryMnemonic,
		recipientAddress,
		assetDenom,
		toAmount,
		memo,
		sequenceNumber,
		gasPrice,
		gasUsed,
	)
	if err != nil {
		return nil, errors.Errorf("TransferStablyToken: %v", err)
	}
	return cosmosTxResult, nil
}

func CalculateHashForTransfer(ctx context.Context,
	senderSecretID string,
	recipientAddress string,
	assetDenom string,
	toAmount int64,
	memo string,
	sequenceNumber uint64,
	gasPrice string,
	gasUsed uint64,
) (string, error) {
	treasuryMnemonic := GetTreasuryMnemonicFromSecretID(ctx, senderSecretID)

	senderInfo, keyring, err := GetKeyringInfoFromMnemonic(treasuryMnemonic)
	if err != nil {
		return "", errors.Errorf("GetKeyringInfoFromMnemonic: %v", err)
	}

	// Retrieve the sender address
	senderAddressString := senderInfo.GetAddress().String()

	clientCtx, txFactory, msg, err := PrepareTransferTransaction(ctx,
		keyring,
		senderAddressString,
		recipientAddress,
		assetDenom,
		toAmount,
		memo,
		sequenceNumber,
		gasPrice,
		gasUsed,
	)
	if err != nil {
		return "", errors.Errorf("PrepareTransferTransaction: %v", err)
	}

	_, signedTransactionInBytes, err := CreateSignedTx(ctx, clientCtx, txFactory, msg)
	if err != nil {
		return "", errors.Errorf("CreateSignedTx: %v", err)
	}

	// Calculate the hash value from the signed transaction
	_, calculatedSignedTransactionInStr := CalculateHashOfTransaction(signedTransactionInBytes)
	return calculatedSignedTransactionInStr, nil
}
