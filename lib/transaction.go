package coreumservicelib

import (
	"context"
	"coreumservicemsg"
	"fmt"
	"strings"
	"time"

	"github.com/CoreumFoundation/coreum/pkg/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/pkg/errors"
	tmtypes "github.com/tendermint/tendermint/types"
)

func CalculateGasForTransfer(ctx context.Context,
	senderMnemonic string,
	recipientAddress string,
	assetDenom string,
	toAmount int64,
	memo string,
	sequenceNumber uint64,
) (uint64, string, error) {
	senderInfo, signingKeyRing, err := GetKeyringInfoFromMnemonic(senderMnemonic)
	if err != nil {
		return 0, "", errors.Errorf("GetKeyringInfoFromMnemonic: %v", err)
	}

	clientCtx, txFactory, msg, err := PrepareTransferTransaction(ctx,
		signingKeyRing,
		senderInfo.GetAddress().String(),
		recipientAddress,
		assetDenom,
		toAmount,
		memo,
		sequenceNumber,
		"", // assume empty
		0,  // assume empty
	)
	if err != nil {
		return 0, "", errors.Errorf("PrepareTransferTransaction: %v", err)
	}

	// calculate gas
	gasUsedForTransaction, gasPrice, err := CalculateGas(ctx, clientCtx, txFactory, msg)
	if err != nil {
		return 0, "", errors.Errorf("CalculateGas: %v", err)
	}

	return gasUsedForTransaction, gasPrice, nil
}

/*
client.CalculateGas returns error if the sequence number has been taken by the submitted blockchain transaction.

So we need to calculate the gas value seperately, and save to DB before submitting the transaction to the blockchain
*/
func CalculateGas(ctx context.Context, clientCtx client.Context, txf client.Factory, msgs ...sdk.Msg) (uint64, string, error) {
	gasPrice, err := client.GetGasPrice(ctx, clientCtx)
	if err != nil {
		return 0, "", errors.Errorf("client.GetGasPrice: %v", err)
	}
	gasPrice.Amount = gasPrice.Amount.Mul(clientCtx.GasPriceAdjustment())

	gasPriceStr := gasPrice.String()

	// Save the previous value of gas price
	oldGasPrice := txf.GasPrices()
	txf = txf.WithGasPrices(gasPriceStr)

	_, adjusted, err := client.CalculateGas(ctx, clientCtx, txf, msgs...)

	// Revert the old gas price
	txf.WithGasPrices(oldGasPrice.String())
	if err != nil {
		return 0, "", errors.Errorf("client.CalculateGas: %v", err)
	}
	return adjusted, gasPriceStr, nil
}

func CreateSignedTx(ctx context.Context, clientCtx client.Context, txf client.Factory, msgs ...sdk.Msg) (signing.Tx, []byte, error) {
	unsignedTx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, nil, errors.Errorf("txf.BuildUnsignedTx: %v", err)
	}

	unsignedTx.SetFeeGranter(clientCtx.FeeGranterAddress())

	// in case the name is not provided by that address, take the name by the address
	fromName := clientCtx.FromName()
	if fromName == "" && len(clientCtx.FromAddress()) > 0 {
		key, err := clientCtx.Keyring().KeyByAddress(clientCtx.FromAddress())
		if err != nil {
			return nil, nil, errors.Errorf("failed to get key by the address %q from the keyring", clientCtx.FromAddress().String())
		}
		fromName = key.GetName()
	}

	err = tx.Sign(txf, fromName, unsignedTx, true)
	if err != nil {
		return nil, nil, errors.Errorf("tx.Sign: %v", err)
	}

	signedTx := unsignedTx.GetTx()

	signedTxBytes, err := clientCtx.TxConfig().TxEncoder()(signedTx)
	if err != nil {
		return nil, nil, errors.Errorf("clientCtx.TxConfig().TxEncoder()(signedTx): %v", err)
	}

	return signedTx, signedTxBytes, nil
}

// Convert the transaction type to byte array
func GetBytesOfTx(clientCtx client.Context, tx signing.Tx) ([]byte, error) {
	return clientCtx.TxConfig().TxEncoder()(tx)
}

// Returns the hash value in byte form and string form
func CalculateHashOfTransaction(transactionInBytes []byte) ([]byte, string) {
	hashByte := tmtypes.Tx(transactionInBytes).Hash()
	return hashByte, fmt.Sprintf("%X", hashByte)
}

// Return nil if there's no transaction with the corresponding hash.
// There's error returned only if there's error when communicating with the blockchain server.
func GetTransactionByHash(ctx context.Context, clientCtx client.Context, transactionHash string) (*coreumservicemsg.CoreumTransactionDetail, error) {
	txSvcClient := sdktx.NewServiceClient(clientCtx)

	requestCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := txSvcClient.GetTx(requestCtx, &sdktx.GetTxRequest{
		Hash: transactionHash,
	})
	if err != nil {
		// there's no transaction with the provided hash on the blockchain
		if strings.Contains(err.Error(), "tx not found") {
			return nil, nil
		}
		return nil, err
	}
	txResponse := res.TxResponse

	if txResponse.Code != 0 {
		// there's no transaction with the provided hash on the blockchain
		return nil, nil
	}

	transactionDetail := &coreumservicemsg.CoreumTransactionDetail{
		Body: coreumservicemsg.CoreumTransactionBody{
			Memo: "",
		},
	}

	return transactionDetail, nil
}
