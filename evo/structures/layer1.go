package structures

type AddressSummaryResponse struct {
	AddrStr                 []string
	Balance                 float64
	BalanceSat              int64
	TotalReceived           float64
	TotalReceivedSat        int64
	TotalSent               float64
	TotalSentSat            int64
	UnconfirmedBalance      float64
	UnconfirmedBalanceSat   int64
	UnconfirmedTxApperances int
	UnconfirmedAppearances  int
	TxApperances            int
	TxAppearances           int
	Transactions            []string
}

type BlockHashResponse string

type BestBlockHashResponse string

type UTXOItemResponse struct {
	Address     string
	Txid        string
	OutputIndex int
	Script      string
	Satoshis    int64
	Height      int
}

type UTXOResponse struct {
	TotalItems int
	From       int
	To         int
	FromHeight int
	ToHeight   int
	Items      []UTXOItemResponse
}

type NodeInfoResponse struct {
	ProRegTxHash   string
	ConfirmedHash  string
	Service        string
	PubKeyOperator string
	VotingAddress  string
	IsValid        bool
}

type QuorumInfoResponse struct {
	Version           int
	LLMQType          int
	QuorumHash        string
	SignersCount      int
	ValidMembersCount int
	QuorumPublicKey   string
}

type MnListDiffResponse struct {
	BaseBlockHash     string
	BlockHash         string
	CbTxMerkleTree    string
	CbTx              string
	DeletedMNs        []NodeInfoResponse
	MnList            []NodeInfoResponse
	DeletedQuorums    []QuorumInfoResponse
	NewQuorums        []QuorumInfoResponse
	MerkleRootMNList  string
	MerkleRootQuorums string
}

type UTXORequest struct {
	From       *int     `json:"from,omitempty"`
	To         *int     `json:"to,omitempty"`
	FromHeight *int     `json:"fromHeight,omitempty"`
	ToHeight   *int     `json:"toHeight,omitempty"`
	Addresses  []string `json:"address"`
}

type BlockRequest struct {
	Hash   *string `json:"hash,omitempty"`
	Height *int    `json:"height,omitempty"`
}

type BloomFilterRequest struct {
	Data     []byte `json:"v_data"`
	HashFunc uint32 `json:"n_hash_funcs"`
	Tweak    uint32 `json:"n_tweak"`
	Flags    uint32 `json:"n_flags"`
}

type SubscribeToTransactionsWithProofsRequest struct {
	BloomFilter           BloomFilterRequest `json:"bloom_filter"`
	Count                 *int               `json:"count,omitempty"`
	FromBlockHash         *[]byte            `json:"from_block_hash,omitempty"`
	FromBlockHeight       *int               `json:"from_block_height,omitempty"`
	SendTransactionHashes *bool              `json:"send_transaction_hashes,omitempty"`
}
