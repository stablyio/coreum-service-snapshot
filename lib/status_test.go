//go:build integration
// +build integration

package coreumservicelib_test

import (
	lib "coreumservice/go/lib"
	"coreumservicemsg"
	"encoding/hex"
	"fmt"

	"github.com/stretchr/testify/require"

	"context"
	"testing"
)

func TestGetLatestBlockHeight(t *testing.T) {
	ctx := context.Background()
	blockStatus, err := lib.GetLatestBlockStatus(ctx)
	require.NoError(t, err)

	t.Log("blockHeight", lib.ToJSONPretty(blockStatus))
	require.Greater(t, blockStatus.LatestBlockHeight, int64(4076985))
	require.Greater(t, blockStatus.LatestBlockTime, int64(1682168444))
}

func TestGetBlockTransactions(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		blockNumber int64
		expected    []*coreumservicemsg.Transaction
	}{
		{
			blockNumber: 4169066,
			expected: []*coreumservicemsg.Transaction{
				{
					TxHash:      "DD4814669E9BFEEBFFDEEB4264BAB30E85FDDC5DA0158CD1194228E6569E4CAB",
					FromAddress: "testcore1av2q6yuaeqw5rqy958842fu6u9xzw62qjy8j3u",
					ToAddress:   "testcore1un00l6nzdg58htj6e9fmx24433srcxpgdft57e",
					Memo:        "testing",
					Coins: []*coreumservicemsg.Coin{
						{
							Amount: "123",
							Denom:  "microusds-testcore162rs3klx73exmyupxlqjju0u7aggcp0fswetn2",
						},
					},
					BlockNumber: 4169066,
				},
			},
		},
		{
			blockNumber: 4294073,
			expected:    []*coreumservicemsg.Transaction{},
		},
		{
			// This case covers the transactions having no MsgSend message
			blockNumber: 4313356,
			expected:    []*coreumservicemsg.Transaction{},
		},
		{
			// For prod config only, this case covers the problem of
			// no concrete type registered for type URL /cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward
			// against interface *types.Msg
			blockNumber: 2513081,
			expected:    []*coreumservicemsg.Transaction{},
		},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("Case %v", testCase.blockNumber), func(it *testing.T) {
			transactions, err := lib.GetBlockTransactions(ctx, testCase.blockNumber)
			require.NoError(it, err)
			require.Equal(it, &testCase.expected, &transactions)
		})
	}
}

func TestGetBlockTransactionsInRange(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		startBlockNumber int64
		endBlockNumber   int64
		expected         []*coreumservicemsg.Transaction
	}{
		{
			startBlockNumber: 4169066,
			endBlockNumber:   4169076,
			expected: []*coreumservicemsg.Transaction{
				{
					TxHash:      "DD4814669E9BFEEBFFDEEB4264BAB30E85FDDC5DA0158CD1194228E6569E4CAB",
					FromAddress: "testcore1av2q6yuaeqw5rqy958842fu6u9xzw62qjy8j3u",
					ToAddress:   "testcore1un00l6nzdg58htj6e9fmx24433srcxpgdft57e",
					Memo:        "testing",
					Coins: []*coreumservicemsg.Coin{
						{
							Amount: "123",
							Denom:  "microusds-testcore162rs3klx73exmyupxlqjju0u7aggcp0fswetn2",
						},
					},
					BlockNumber: 4169066,
				},
			},
		},
		{
			// This case covers the blocks having no MsgSend-based transactions
			startBlockNumber: 4313356,
			endBlockNumber:   4313367,
			expected:         []*coreumservicemsg.Transaction{},
		},
	}

	for _, testCase := range cases {
		t.Run(
			fmt.Sprintf(
				"Case [%v, %v] (%v blocks)",
				testCase.startBlockNumber,
				testCase.endBlockNumber,
				testCase.endBlockNumber-testCase.startBlockNumber+1,
			),
			func(it *testing.T) {
				transactions, err := lib.GetBlockTransactionsInRange(ctx, testCase.startBlockNumber, testCase.endBlockNumber)
				require.NoError(it, err)
				require.NoError(it, err)
				require.Equal(it, &testCase.expected, &transactions)
			})
	}
}

func TestGetTransactionFromTxBytes(t *testing.T) {
	t.Run("Case with MsgSend transaction", func(it *testing.T) {
		expected := &coreumservicemsg.Transaction{
			TxHash:      "DD4814669E9BFEEBFFDEEB4264BAB30E85FDDC5DA0158CD1194228E6569E4CAB",
			FromAddress: "testcore1av2q6yuaeqw5rqy958842fu6u9xzw62qjy8j3u",
			ToAddress:   "testcore1un00l6nzdg58htj6e9fmx24433srcxpgdft57e",
			Memo:        "testing",
			Coins: []*coreumservicemsg.Coin{
				{
					Amount: "123",
					Denom:  "microusds-testcore162rs3klx73exmyupxlqjju0u7aggcp0fswetn2",
				},
			},
		}

		txBytesHex := "0ad1010ac5010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e6412a4010a2f74657374636f7265316176327136797561657177357271793935383834326675367539787a773632716a79386a3375122f74657374636f726531756e30306c366e7a6467353868746a366539666d7832343433337372637870676466743537651a400a396d6963726f757364732d74657374636f72653136327273336b6c78373365786d797570786c716a6a7530753761676763703066737765746e321203313233120774657374696e67126b0a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21035841d46c964f7356e38bf912ecf5e9834913420ef12bc01f1639e1fddbf331de12040a020801180812170a110a097574657374636f72651204333935321090c2041a40e714315cb9391b563394aa29bc79bb77a7d9d7fe84ce01f3128b3d20f25c330e35680132c5c41ee7f1ac9ff5d3616e4b8861abb8936d6340f931d4073e1b940b"
		txBytes, err := hex.DecodeString(txBytesHex)
		require.NoError(it, err)

		tx, err := lib.GetTransactionFromTxBytes(txBytes)
		require.NoError(it, err)
		require.Equal(it, expected, tx)
	})

	t.Run("Case without MsgSend transaction", func(it *testing.T) {
		txBytesHex := "0abf010abc010a212f636f736d6f732e62616e6b2e763162657461312e4d73674d756c746953656e641296010a490a2f74657374636f7265313334346a68376b6767347134667a7075616d7274716a657375786579656e3730306e6c76653612160a097574657374636f7265120931303030303030303012490a2f74657374636f726531786a65686d7479327a356a376d666d707a786538646772663530366337306e3337343763393512160a097574657374636f72651209313030303030303030126d0a520a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21022429fe19f84c6f389f6c4b172b5ade512d8aee1aebc5c3eb6f72d97720870bf712040a02080118ffe40912170a110a097574657374636f726512043332383610c0b2041a40cbd0527dd20a60ac602bb48abea538ea802cfb749f785e1aed7f5cf177b3c6fc5ef65f9abb4328b6ddacefe999eae1a7278cb5661a91f76c1ea89cbf93c8e34a"
		txBytes, err := hex.DecodeString(txBytesHex)
		require.NoError(it, err)

		tx, err := lib.GetTransactionFromTxBytes(txBytes)
		require.NoError(it, err)
		require.Nil(it, tx)
	})
}

func TestGetAccountInfo(t *testing.T) {
	ctx := context.Background()

	validCases := []struct {
		address  string
		expected *coreumservicemsg.GetAccountInfoReply
	}{
		{
			address:  "testcore1av2q6yuaeqw5rqy958842fu6u9xzw62qjy8j3u",
			expected: &coreumservicemsg.GetAccountInfoReply{},
		},
	}

	for _, testCase := range validCases {
		t.Run(fmt.Sprintf("Case %v", testCase.address), func(it *testing.T) {
			accountInfo, err := lib.GetAccountInfo(ctx, testCase.address)
			require.NoError(it, err)

			it.Log("Sequence", accountInfo.Sequence)
			it.Log("AccountNumber", accountInfo.AccountNumber)

			require.GreaterOrEqual(it, accountInfo.Sequence, uint64(0))
			require.GreaterOrEqual(it, accountInfo.AccountNumber, uint64(0))
		})
	}
}

func TestGetAccountInfoFromMnemonic(t *testing.T) {
	ctx := context.Background()

	mnemonic := "nut clog audit reward display era divide galaxy boil sport bless disorder total hidden pair range senior risk disorder affair dress barrel nuclear exhibit"

	accountInfo, err := lib.GetAccountInfoFromMnemonic(ctx, mnemonic)
	require.NoError(t, err)

	t.Log("accountInfo.Sequence", accountInfo.Sequence)
	t.Log("accountInfo.AccountNumber", accountInfo.AccountNumber)

	require.GreaterOrEqual(t, accountInfo.Sequence, uint64(0))
	require.GreaterOrEqual(t, accountInfo.AccountNumber, uint64(0))

	// Now, double check the result with GetAccountInfo()

	keyringInfo, _, err := lib.GetKeyringInfoFromMnemonic(mnemonic)
	require.NoError(t, err)

	address := keyringInfo.GetAddress().String()
	accountInfo2, err := lib.GetAccountInfo(ctx, address)
	require.NoError(t, err)

	t.Log("accountInfo2.Sequence", accountInfo2.Sequence)
	t.Log("accountInfo2.AccountNumber", accountInfo2.AccountNumber)

	require.Equal(t, accountInfo, accountInfo2)
}
