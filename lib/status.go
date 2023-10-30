package coreumservicelib

import (
	"context"
	"coreumservicemsg"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	app "github.com/CoreumFoundation/coreum/app"
	"github.com/pkg/errors"

	"github.com/CoreumFoundation/coreum/pkg/client"
	"github.com/CoreumFoundation/coreum/pkg/config"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func GetLatestBlockStatus(ctx context.Context) (*coreumservicemsg.BlockStatus, error) {
	rpcClient, err := GetTendermintRPCClient()
	if err != nil {
		return nil, errors.Errorf("GetTendermintRPCClient: %v", err)
	}

	// Reference from Coreum team: https://pastebin.com/MgJS98Jz
	status, err := rpcClient.Status(ctx)
	if err != nil {
		return nil, errors.Errorf("rpcClient.Status: %v", err)
	}
	if status == nil {
		return nil, errors.Errorf("got nil status")
	}
	return &coreumservicemsg.BlockStatus{
		LatestBlockHash:     status.SyncInfo.LatestBlockHash.String(),
		LatestAppHash:       status.SyncInfo.LatestAppHash.String(),
		LatestBlockHeight:   status.SyncInfo.LatestBlockHeight,
		LatestBlockTime:     status.SyncInfo.LatestBlockTime.Unix(),
		EarliestBlockHash:   status.SyncInfo.EarliestBlockHash.String(),
		EarliestAppHash:     status.SyncInfo.EarliestAppHash.String(),
		EarliestBlockHeight: status.SyncInfo.EarliestBlockHeight,
		EarliestBlockTime:   status.SyncInfo.EarliestBlockTime.Unix(),
		CatchingUp:          status.SyncInfo.CatchingUp,
	}, nil
}

func GetBlockTransactions(ctx context.Context, blockNumber int64) ([]*coreumservicemsg.Transaction, error) {
	grpcClient := GetGRPCClient()
	modules := module.NewBasicManager(
		bank.AppModuleBasic{},
	)
	clientCtx := client.NewContext(client.DefaultContextConfig(), modules).WithGRPCClient(grpcClient)
	txClient := txtypes.NewServiceClient(clientCtx)

	blockResults, err := txClient.GetBlockWithTxs(ctx, &txtypes.GetBlockWithTxsRequest{
		Height: int64(blockNumber),
	})
	if err != nil {
		/*
			The txClient.GetBlockWithTxs() will return an error when the block has 0 transaction
			- Thus, to detect if the call was failed due to 0 txns, we double check it with the tendermintRPCClient.BlockResults
			- If the block has 0 transaction, we can return an empty list without returning an error
		*/
		tendermintRPCClient, err := GetTendermintRPCClient()
		if err != nil {
			return nil, errors.Errorf("GetTendermintRPCClient(%v): %v", blockNumber, err)
		}
		blockRes, err := tendermintRPCClient.BlockResults(ctx, &blockNumber)
		if err != nil {
			return nil, errors.Errorf("tendermintRPCClient.BlockResults(%v): %v", blockNumber, err)
		}
		if len(blockRes.TxsResults) == 0 {
			// Return empty list of detected transactions
			return []*coreumservicemsg.Transaction{}, nil
		}
		return nil, errors.Errorf("txClient.GetBlockWithTxs(%v): %v", blockNumber, err)
	}

	res := []*coreumservicemsg.Transaction{}
	for _, txBytes := range blockResults.Block.Data.Txs {
		tx, err := GetTransactionFromTxBytes(txBytes)
		if err != nil {
			return nil, errors.Errorf("GetTransactionFromTxBytes failed with txBytes [%v]: %v", txBytes, err)
		}
		if tx != nil {
			// Assign the block number to Coreum transaction
			// because there's no block number embedded in the transaction's body
			tx.BlockNumber = uint64(blockNumber)
			res = append(res, tx)
		}
	}

	return res, nil
}

type getBlockTransactionsChannelOutput struct {
	Res         []*coreumservicemsg.Transaction
	BlockNumber int64
	Err         error
}

func GetBlockTransactionsInRange(ctx context.Context, startblockNumber int64, endblockNumber int64) ([]*coreumservicemsg.Transaction, error) {
	// Speed up the queries with multiple goroutines
	channels := []chan getBlockTransactionsChannelOutput{}
	for i := startblockNumber; i <= endblockNumber; i++ {
		ch := make(chan getBlockTransactionsChannelOutput)
		go fetchTransactionInBlock(ctx, ch, i)
		channels = append(channels, ch)
	}

	// Join threads
	res := []*coreumservicemsg.Transaction{}
	for _, ch := range channels {
		out := <-ch
		if out.Err != nil {
			return nil, errors.Errorf("error from goroutine at block %v: %v", out.BlockNumber, out.Err)
		}
		res = append(res, out.Res...)
	}

	return res, nil
}

func fetchTransactionInBlock(ctx context.Context, ich chan getBlockTransactionsChannelOutput, blockNumber int64) {
	maxTrials := 10
	for {
		maxTrials -= 1

		txs, err := GetBlockTransactions(ctx, blockNumber)
		if err != nil {
			if maxTrials < 1 {
				ich <- getBlockTransactionsChannelOutput{
					Res:         nil,
					BlockNumber: blockNumber,
					Err:         errors.Errorf("GetBlockTransactions(%v): %v", blockNumber, err),
				}
				break
			}
			time.Sleep(5 * time.Second)
			continue
		}
		ich <- getBlockTransactionsChannelOutput{
			Res:         txs,
			BlockNumber: blockNumber,
			Err:         nil,
		}
		break
	}
}

// Try to find a MsgSend message from the txBytes, return nil if there is no such kinda of message
func GetTransactionFromTxBytes(txBytes []byte) (*coreumservicemsg.Transaction, error) {
	modules := app.ModuleBasics
	encodingConfig := config.NewEncodingConfig(modules)

	tx := &txtypes.Tx{}
	err := encodingConfig.Codec.Unmarshal(txBytes, tx)
	if err != nil {
		return nil, errors.Errorf("encodingConfig.Codec.Unmarshal: %v", err)
	}

	for _, msg := range tx.GetMsgs() {
		bankSend, ok := msg.(*banktypes.MsgSend)
		if !ok {
			continue
		}

		coins := []*coreumservicemsg.Coin{}
		for _, amount := range bankSend.Amount {
			coins = append(coins, &coreumservicemsg.Coin{
				Amount: amount.Amount.String(),
				Denom:  amount.Denom,
			})
		}

		txHash := strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256(txBytes)))

		return &coreumservicemsg.Transaction{
			TxHash:      txHash,
			FromAddress: bankSend.FromAddress,
			ToAddress:   bankSend.ToAddress,
			Memo:        tx.Body.Memo,
			Coins:       coins,
		}, nil
	}
	return nil, nil
}

func GetAccountInfo(ctx context.Context, address string) (*coreumservicemsg.GetAccountInfoReply, error) {
	cosmosAddress, err := cosmossdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, errors.Errorf("cosmossdk.AccAddressFromBech32: %v", err)
	}

	acc, err := client.GetAccountInfo(ctx, GetClientContext(), cosmosAddress)
	if err != nil {
		return nil, errors.Errorf("client.GetAccountInfo: %v", err)
	}

	return &coreumservicemsg.GetAccountInfoReply{
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}, nil
}

func GetAccountInfoFromMnemonic(ctx context.Context, mnemonic string) (*coreumservicemsg.GetAccountInfoReply, error) {
	keyringInfo, _, err := GetKeyringInfoFromMnemonic(mnemonic)
	if err != nil {
		return nil, errors.Errorf("cosmossdk.AccAddressFromBech32: %v", err)
	}

	address := keyringInfo.GetAddress().String()

	res, err := GetAccountInfo(ctx, address)
	if err != nil {
		return nil, errors.Errorf("GetAccountInfo: %v", err)
	}

	return res, nil
}
