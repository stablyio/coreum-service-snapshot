//go:build integration
// +build integration

package coreumservicelib_test

import (
	"context"
	lib "coreumservice/go/lib"
	"fmt"
	"testing"

	"coreumservice/go/stably_io/config"

	"github.com/stretchr/testify/require"
)

func TestGetAddress(t *testing.T) {
	validCases := []struct {
		mnenomic string
		index    int64
		address  string
		path     string
	}{
		{
			mnenomic: "hazard misery record advice ceiling clean manage ten approve render abstract horse door federal congress stadium job tribe begin shaft digital aerobic upset record",
			index:    0,
			address:  "testcore1un00l6nzdg58htj6e9fmx24433srcxpgdft57e",
			path:     "m/44'/990'/0'/0/0",
		},
		{
			mnenomic: "hazard misery record advice ceiling clean manage ten approve render abstract horse door federal congress stadium job tribe begin shaft digital aerobic upset record",
			index:    1,
			address:  "testcore1kjhpcsee5ptet33xe8fcwmzvtzj3xgy9tl52a3",
			path:     "m/44'/990'/1'/0/0",
		},
		{
			mnenomic: "hazard misery record advice ceiling clean manage ten approve render abstract horse door federal congress stadium job tribe begin shaft digital aerobic upset record",
			index:    2,
			address:  "testcore13akgpj55l77vsh442rqssqc32znudk8k44kkxs",
			path:     "m/44'/990'/2'/0/0",
		},
	}

	for _, testCase := range validCases {
		t.Run(fmt.Sprintf("Case [%v] with index %v", testCase.address, testCase.index), func(it *testing.T) {
			addressInfo, err := lib.GetAddress(testCase.mnenomic, testCase.index)
			require.NoError(it, err)

			it.Log("Address", addressInfo.Address)
			it.Log("DerivationPath", addressInfo.DerivationPath)
			require.Equal(it, testCase.address, addressInfo.Address)
			require.Equal(it, testCase.path, addressInfo.DerivationPath)
		})
	}
}

func TestGetTreasuryAddress(t *testing.T) {
	addressInfo, err := lib.GetTreasuryAddress(context.Background(), config.GetConfigDefault().Blockchain.Coreum.USDS.TreasurySecretID)
	require.NoError(t, err)
	require.Equal(t, "testcore1av2q6yuaeqw5rqy958842fu6u9xzw62qjy8j3u", addressInfo.Address)
}
