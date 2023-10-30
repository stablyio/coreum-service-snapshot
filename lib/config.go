package coreumservicelib

import (
	"sync"

	coreumconstant "github.com/CoreumFoundation/coreum/pkg/config/constant"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
)

// Config for this service
var (
	configInit sync.Once
)

func init() {
	configInit.Do(func() {
		// Configure Cosmos SDK
		config := cosmossdk.GetConfig()

		addressPrefix := GetAddressPrefixByStage()
		if addressPrefix == "" {
			panic("Failed to get the address prefix")
		}
		
		config.SetBech32PrefixForAccount(addressPrefix, addressPrefix + "pub")
		config.SetCoinType(coreumconstant.CoinType)
		config.Seal()
	})
}
