package coreumservicelib

import (
	"context"
	"coreumservice/go/stably_io/config"
	"coreumservice/go/stably_io/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"coreumservicemsg"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const prefix = "coreumservice"

func RunHttpServer() {
	r := mux.NewRouter()

	setupHealthCheckHandler(r)

	// Register endpoints
	validateIssuanceParams(r)
	getTreasuryAddress(r)
	getLatestBlockStatus(r)
	getBlockTransactions(r)
	getBlockTransactionsInRange(r)

	// Method to broadcast the issuance
	transferStablyToken(r)

	if utils.GetStage() != "prod" {
		// Transfer token given the sender mnemonic (just for integration test)
		transferTokenWithMnemonic(r)
	}

	// Get the information of the address from the blockchain
	getAccountInfoByAddress(r)

	// Method to return the used gas and gas price for the transfer transaction
	getGasForTransferStablyToken(r)

	// Method to get calculate the hash by the parameters
	calculateHashOfTransfer(r)

	// Return the transaction detail by the transaction hash
	getTransactionByHashRequest(r)

	// Method to query the balance of the address for the denom
	getBalanceOfAddressForDenom(r)

	port := config.GetConfigDefault().Blockchain.Coreum.PRCConfig.HTTPServerPort
	fmt.Printf("Start http server at port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func setupHealthCheckHandler(r *mux.Router) *mux.Route {
	return r.HandleFunc(fmt.Sprintf("/%s/healthcheck", prefix), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})
}

// Method to validate the issuance params
func validateIssuanceParams(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		"validate-transfer-params",
		func(input *coreumservicemsg.ValidateTransferParamsRequest) (*coreumservicemsg.ValidateTransferParamsReply, error) {
			return &coreumservicemsg.ValidateTransferParamsReply{}, nil
		})
}

func getAccountInfoByAddress(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		"get-account-info",
		func(input *coreumservicemsg.GetAccountInfoRequest) (*coreumservicemsg.GetAccountInfoReply, error) {
			ctx := context.Background()
			acc, err := GetAccountInfo(ctx, input.Address)
			if err != nil {
				return nil, errors.Errorf("GetAccountInfo: %v", err)
			}
			return &coreumservicemsg.GetAccountInfoReply{
				AccountNumber: acc.AccountNumber,
				Sequence:      acc.Sequence,
			}, nil
		})
}

func getGasForTransferStablyToken(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		"get-gas-for-transfer-stably-token",
		func(input *coreumservicemsg.GetGasForTransferStablyTokenRequest) (*coreumservicemsg.GetGasForTransferStablyTokenReply, error) {
			ctx := context.Background()

			treasuryMnemonic := GetTreasuryMnemonicFromSecretID(ctx, input.SenderSecretID)

			gasUsed, gasPrice, err := CalculateGasForTransfer(ctx,
				treasuryMnemonic,
				input.RecipientAddress,
				input.TokenDenom,
				input.TokenAmount,
				input.Memo,
				input.SequenceNumber,
			)
			if err != nil {
				return nil, errors.Errorf("CalculateGasForTransfer: %v", err)
			}

			return &coreumservicemsg.GetGasForTransferStablyTokenReply{
				GasUsed:  gasUsed,
				GasPrice: gasPrice,
			}, nil
		})
}

func transferStablyToken(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"transfer-stably-token",
		// The processing function
		func(input *coreumservicemsg.TransferStablyTokenRequest) (*coreumservicemsg.TransferStablyTokenReply, error) {
			ctx := context.Background()

			txResponse, err := TransferStablyToken(ctx,
				input.SenderSecretID,
				input.RecipientAddress,
				input.TokenDenom,
				input.TokenAmount,
				input.Memo,
				input.SequenceNumber,
				input.GasPrice,
				input.GasUsed,
			)
			if err != nil {
				return nil, errors.Errorf("TransferStablyToken: %v", err)
			}
			return &coreumservicemsg.TransferStablyTokenReply{
				TxHash: txResponse.TxHash,
			}, nil
		},
	)
}

func transferTokenWithMnemonic(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"test-only/transfer-stably-token",
		// The processing function
		func(input *coreumservicemsg.TransferTokenWithMnemonicRequest) (*coreumservicemsg.TransferTokenWithMnemonicReply, error) {
			ctx := context.Background()

			gasPrice := input.GasPrice
			gasUsed := input.GasUsed

			if gasPrice == "" || gasUsed == uint64(0) {
				newGasUsed, newGasPrice, err := CalculateGasForTransfer(ctx,
					input.SenderMnemonic,
					input.RecipientAddress,
					input.TokenDenom,
					input.TokenAmount,
					input.Memo,
					input.SequenceNumber,
				)
				if err != nil {
					return nil, errors.Errorf("CalculateGasForTransfer: %v", err)
				}
				gasPrice = newGasPrice
				gasUsed = newGasUsed
			}

			txResponse, err := TransferTokenWithMnemonic(ctx,
				input.SenderMnemonic,
				input.RecipientAddress,
				input.TokenDenom,
				input.TokenAmount,
				input.Memo,
				input.SequenceNumber,
				gasPrice,
				gasUsed,
			)
			if err != nil {
				return nil, errors.Errorf("TransferTokenWithMnemonic: %v", err)
			}
			return &coreumservicemsg.TransferTokenWithMnemonicReply{
				TxHash: txResponse.TxHash,
			}, nil
		},
	)
}

func calculateHashOfTransfer(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"calculate-hash-of-transfer",
		// The processing function
		func(input *coreumservicemsg.CalculateHashOfTransactionRequest) (*coreumservicemsg.CalculateHashOfTransactionReply, error) {
			ctx := context.Background()

			calculatedHash, err := CalculateHashForTransfer(ctx,
				input.SenderSecretID,
				input.RecipientAddress,
				input.TokenDenom,
				input.TokenAmount,
				input.Memo,
				input.SequenceNumber,
				input.GasPrice,
				input.GasUsed,
			)
			if err != nil {
				return nil, errors.Errorf("CalculateHashForTransfer: %v", err)
			}
			return &coreumservicemsg.CalculateHashOfTransactionReply{
				CalculatedTxHash: calculatedHash,
			}, nil
		},
	)
}

func getTreasuryAddress(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"get-treasury-address",
		// The processing function
		func(input *coreumservicemsg.GetTreasuryAddressRequest) (*coreumservicemsg.GetTreasuryAddressReply, error) {
			ctx := context.Background()
			addressInfo, err := GetTreasuryAddress(ctx, input.TreasurySecretID)
			if err != nil {
				return nil, errors.Errorf("GetTreasuryAddress: %v", err)
			}
			return &coreumservicemsg.GetTreasuryAddressReply{
				Address: addressInfo.Address,
				Path:    addressInfo.DerivationPath,
			}, nil
		},
	)
}

func getLatestBlockStatus(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"get-latest-block-status",
		// The processing function
		func(_ *struct{}) (*coreumservicemsg.GetLatestBlockStatusReply, error) {
			ctx := context.Background()
			blockStatus, err := GetLatestBlockStatus(ctx)
			if err != nil {
				return nil, errors.Errorf("GetLatestBlockStatus: %v", err)
			}
			return &coreumservicemsg.GetLatestBlockStatusReply{
				BlockStatus: blockStatus,
			}, nil
		},
	)
}

func getBlockTransactions(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"get-block-transactions",
		// The processing function
		func(input *coreumservicemsg.GetBlockTransactionsRequest) (*coreumservicemsg.GetBlockTransactionsReply, error) {
			ctx := context.Background()
			transactions, err := GetBlockTransactions(ctx, int64(input.BlockNumber))
			if err != nil {
				return nil, errors.Errorf("GetBlockTransactions: %v", err)
			}
			return &coreumservicemsg.GetBlockTransactionsReply{
				Transactions: transactions,
			}, nil
		},
	)
}

func getBlockTransactionsInRange(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"get-block-transactions-in-range",
		// The processing function
		func(input *coreumservicemsg.GetBlockTransactionsInRangeRequest) (*coreumservicemsg.GetBlockTransactionsInRangeReply, error) {
			ctx := context.Background()
			transactions, err := GetBlockTransactionsInRange(ctx, int64(input.StartBlockNumber), int64(input.EndBlockNumber))
			if err != nil {
				return nil, errors.Errorf("GetBlockTransactionsInRange: %v", err)
			}
			return &coreumservicemsg.GetBlockTransactionsInRangeReply{
				Transactions: transactions,
			}, nil
		},
	)
}

func getTransactionByHashRequest(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"get-transaction-by-hash",
		// The processing function
		func(input *coreumservicemsg.GetTransactionRequest) (*coreumservicemsg.GetTransactionReply, error) {
			ctx := context.Background()
			transactionDetail, err := GetTransactionByHash(ctx, GetClientContext(), input.TransactionHash)
			if err != nil {
				return nil, errors.Errorf("GetTransactionByHash: %v", err)
			}
			return &coreumservicemsg.GetTransactionReply{
				TransactionDetail: transactionDetail,
			}, nil
		},
	)
}

func getBalanceOfAddressForDenom(r *mux.Router) *mux.Route {
	return httpEndpointProcessing(r,
		// The endpoint
		"get-balance-of-address",
		// The processing function
		func(input *coreumservicemsg.GetBalanceOfAddressForDenomRequest) (*coreumservicemsg.GetBalanceOfAddressForDenomReply, error) {
			ctx := context.Background()
			balanceAmount, err := GetBalanceOfAddress(ctx, input.Address, input.Denom)
			if err != nil {
				return nil, errors.Errorf("GetBalanceOfAddress: address(%s), denom(%s), %+v", input.Address, input.Denom, err)
			}
			return &coreumservicemsg.GetBalanceOfAddressForDenomReply{
				Amount: balanceAmount,
			}, nil
		},
	)
}

func httpEndpointProcessing[T any, R any](
	r *mux.Router,
	endpoint string,
	processFunc func(input *T) (*R, error),
) *mux.Route {
	return r.HandleFunc(fmt.Sprintf("/%s/%s", prefix, endpoint), func(writer http.ResponseWriter, request *http.Request) {
		readBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			handleBadRequest(writer, err, endpoint)
			return
		}
		var requestParams T
		err = json.Unmarshal(readBytes, &requestParams)
		if err != nil {
			handleBadRequest(writer, err, endpoint)
			return
		}

		// Process the response
		response, err := processFunc(&requestParams)
		if err != nil {
			handleBadRequest(writer, err, endpoint)
			return
		}
		if response == nil {
			err = errors.Errorf("got nil response from processing")
			handleBadRequest(writer, err, endpoint)
			return
		}

		responseMessage, err := json.Marshal(response)
		if err != nil {
			handleBadRequest(writer, err, endpoint)
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(responseMessage)
		if err != nil {
			fmt.Printf("[%v] Error from writer.Write(responseMessage): %v\n", endpoint, err)
		}
	})
}

func handleBadRequest(writer http.ResponseWriter, err error, endpoint string) {
	writer.WriteHeader(http.StatusBadRequest)
	if err != nil {
		fmt.Printf("[%v] Bad request: %+v\n", endpoint, err)
	}
	_, err = writer.Write([]byte(
		fmt.Sprintf("{\"ErrorMessage\" : \"bad request: %v\"}", err),
	))
	if err != nil {
		fmt.Printf("[%v] Error from writer.Write(responseMessage): %+v\n", endpoint, err)
	}
}
