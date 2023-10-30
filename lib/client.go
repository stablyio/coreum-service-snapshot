package coreumservicelib

import (
	"coreumservice/go/stably_io/config"
	"crypto/tls"

	"github.com/CoreumFoundation/coreum/pkg/client"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Generate the base client context
func GetClientContext() client.Context {
	// List required modules.
	// If you need types from any other module import them and add here.
	modules := module.NewBasicManager(
		auth.AppModuleBasic{},
	)

	gprcClient := GetGRPCClient()

	cosmosClientCtx := client.NewContext(client.DefaultContextConfig(), modules).
		WithChainID(GetChainIDByStage()).
		WithGRPCClient(gprcClient).
		WithKeyring(keyring.NewInMemory()).
		WithBroadcastMode(flags.BroadcastBlock)

	return cosmosClientCtx // Config.Coreum.CosmosClientCtx
}

func GetGRPCClient() *grpc.ClientConn {
	coreumConfig := config.GetConfigDefault().Blockchain.Coreum
	// Configure client context and tx factory
	// If you don't use TLS then replace `grpc.WithTransportCredentials()` with `grpc.WithInsecure()`
	gprcClient, err := grpc.Dial(coreumConfig.PRCConfig.GRPCNodeURL, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		panic(err)
	}
	return gprcClient
}

// Transaction Factory is generated from the client Context.
// It is used for:
// - Sign the transaction.
// - Fetch the sequence number
func CoreumTxFactory(clientCtx client.Context) tx.Factory {
	txFactory := client.Factory{}.
		WithKeybase(clientCtx.Keyring()).
		WithChainID(clientCtx.ChainID()).
		WithTxConfig(clientCtx.TxConfig()).
		// Not try to simulating and executing the transaction before committing to the blockchain
		WithSimulateAndExecute(false)
	return txFactory
}

func GetTendermintRPCClient() (*http.HTTP, error) {
	coreumConfig := config.GetConfigDefault().Blockchain.Coreum
	// Reference from Coreum team: https://pastebin.com/MgJS98Jz
	rpcClient, err := cosmosClient.NewClientFromNode(coreumConfig.PRCConfig.TendermintRPCNodeURL)
	if err != nil {
		return nil, errors.Errorf("cosmosClient.NewClientFromNode: %v", err)
	}
	return rpcClient, nil
}
