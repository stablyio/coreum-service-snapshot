package coreumservicelib

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func GetBalanceOfAddress(ctx context.Context, recipientAddress string, denom string) (string, error) {
	clientCtx := GetClientContext()
	bankClient := banktypes.NewQueryClient(clientCtx)
	balance, err := bankClient.Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: recipientAddress,
		Denom:   denom,
	})
	if err != nil {
		return "", err
	}
	return balance.Balance.Amount.String(), nil
}
