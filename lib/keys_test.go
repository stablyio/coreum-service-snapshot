//go:build integration
// +build integration

package coreumservicelib_test

import (
	"context"
	lib "coreumservice/go/lib"
	"testing"

	"coreumservice/go/stably_io/config"

	"github.com/stretchr/testify/require"
)

func TestGetTreasuryMnemonic(t *testing.T) {
	ctx := context.Background()
	mnemonic := lib.GetTreasuryMnemonicFromSecretID(ctx, config.GetConfigDefault().Blockchain.Coreum.USDS.TreasurySecretID)

	expectedTestMnemonic := "enemy liberty cotton cost wrist abuse swear staff very bar critic genre elbow heart unaware deliver witness target relax genre chaos visa risk dutch"

	t.Log("mnemonic", mnemonic)
	require.Equal(t, expectedTestMnemonic, mnemonic)
}
