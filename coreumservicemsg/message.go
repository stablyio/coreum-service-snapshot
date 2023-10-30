package coreumservicemsg

type ValidateTransferParamsRequest struct {
	ToTokenDenom string `json:"to_token_denom"`
	ToAddress    string `json:"to_address"`
	ToAmount     string `json:"to_amount"`
}

type ValidateTransferParamsReply struct{}

type GetGasForTransferStablyTokenRequest struct {
	SenderSecretID   string `json:"sender_secret_id"`
	TokenDenom       string `json:"token_denom"`
	TokenAmount      int64  `json:"token_amount"`
	RecipientAddress string `json:"recipient_address"`
	SequenceNumber   uint64 `json:"sequence_number"`
	Memo             string `json:"memo,omitempty"` // optional
}

type GetGasForTransferStablyTokenReply struct {
	GasPrice string `json:"gas_price"`
	GasUsed  uint64 `json:"gas_used"`
}

type TransferStablyTokenRequest struct {
	SenderSecretID   string `json:"sender_secret_id"`
	TokenDenom       string `json:"token_denom"`
	SequenceNumber   uint64 `json:"sequence_number"`
	TokenAmount      int64  `json:"token_amount"`
	Memo             string `json:"memo,omitempty"` // optional
	RecipientAddress string `json:"recipient_address"`
	GasPrice         string `json:"gas_price"`
	GasUsed          uint64 `json:"gas_used"`
}

type TransferStablyTokenReply struct {
	// The transaction hash after submitting to the blockchain
	TxHash string `json:"tx_hash"`
}

type TransferTokenWithMnemonicRequest struct {
	SenderMnemonic   string `json:"sender_mnemonic"`
	RecipientAddress string `json:"recipient_address"`
	TokenDenom       string `json:"token_denom"`
	TokenAmount      int64  `json:"token_amount"`
	Memo             string `json:"memo,omitempty"` // optional
	SequenceNumber   uint64 `json:"sequence_number"`
	GasPrice         string `json:"gas_price,omitempty"` // auto-suggested if not filled
	GasUsed          uint64 `json:"gas_used,omitempty"`  // auto-suggested if not filled
}

type TransferTokenWithMnemonicReply struct {
	// The transaction hash after submitting to the blockchain
	TxHash string `json:"tx_hash"`
}

type CalculateHashOfTransactionRequest struct {
	SenderSecretID   string `json:"sender_secret_id"`
	TokenDenom       string `json:"token_denom"`
	TokenAmount      int64  `json:"token_amount"`
	Memo             string `json:"memo,omitempty"` // optional
	RecipientAddress string `json:"recipient_address"`
	SequenceNumber   uint64 `json:"sequence_number"`
	GasPrice         string `json:"gas_price"`
	GasUsed          uint64 `json:"gas_used"`
}

type CalculateHashOfTransactionReply struct {
	// The transaction hash calculated from the supplied value
	CalculatedTxHash string `json:"calculated_tx_hash"`
}

type GetTreasuryAddressRequest struct {
	TreasurySecretID string `json:"treasury_secret_id"`
}
type GetTreasuryAddressReply struct {
	Address string `json:"address"`
	Path    string `json:"path"`
}

type GetAccountInfoRequest struct {
	Address string `json:"address"`
}

type GetAccountInfoReply struct {
	AccountNumber uint64 `json:"account_number"`
	Sequence      uint64 `json:"sequence"`
}

type CoreumTransactionBody struct {
	Memo string
}

// Reference: https://docs.coreum.dev/api/api.html#cosmos.tx.v1beta1.Tx
type CoreumTransactionDetail struct {
	Body CoreumTransactionBody
}

type GetTransactionRequest struct {
	TransactionHash string `json:"transaction_hash"`
}

type GetTransactionReply struct {
	TransactionDetail *CoreumTransactionDetail
}

type BlockStatus struct {
	// Reference: the SyncInfo in github.com/informalsystems/tendermint@v0.34.26/rpc/core/types/responses.go
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestAppHash     string `json:"latest_app_hash"`
	LatestBlockHeight int64  `json:"latest_block_height"`
	LatestBlockTime   int64  `json:"latest_block_time"`

	EarliestBlockHash   string `json:"earliest_block_hash"`
	EarliestAppHash     string `json:"earliest_app_hash"`
	EarliestBlockHeight int64  `json:"earliest_block_height"`
	EarliestBlockTime   int64  `json:"earliest_block_time"`

	CatchingUp bool `json:"catching_up"`
}

type GetLatestBlockStatusReply struct {
	BlockStatus *BlockStatus `json:"block_status"`
}

type Transaction struct {
	TxHash      string  `json:"tx_hash"`
	FromAddress string  `json:"from_address"`
	ToAddress   string  `json:"to_address"`
	Memo        string  `json:"memo,omitempty"`
	Coins       []*Coin `json:"coins"`
	BlockNumber uint64  `json:"block_number"`
}

type Coin struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type GetBlockTransactionsRequest struct {
	BlockNumber uint64 `json:"block_number"`
}
type GetBlockTransactionsReply struct {
	Transactions []*Transaction `json:"transactions"`
}

type GetBlockTransactionsInRangeRequest struct {
	StartBlockNumber uint64 `json:"start_block_number"`
	EndBlockNumber   uint64 `json:"end_block_number"`
}
type GetBlockTransactionsInRangeReply struct {
	Transactions []*Transaction `json:"transactions"`
}

type GetBalanceOfAddressForDenomRequest struct {
	Address string `json:"address"`
	Denom   string `json:"denom"`
}

type GetBalanceOfAddressForDenomReply struct {
	Amount string `json:"amount"`
}

type Error struct {
	ErrorMessage string `json:"ErrorMessage"`
}
