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

type UTXOItem struct {
	Address     string
	Txid        string
	OutputIndex int
	Script      string
	Satoshis    int64
	Height      int
}

type UTXORequest struct {
	*LimitRange
	Addresses []string `json:"address"`
}

type UTXOResponse struct {
	TotalItems int
	From       int
	To         int
	FromHeight int
	ToHeight   int
	Items      []UTXOItem
}

type LimitRange struct {
	From       int
	To         int
	FromHeight int
	ToHeight   int
}

type BlockRequest struct {
	Hash   string
	Height int
}

type NodeInfo struct {
	ProRegTxHash   string
	ConfirmedHash  string
	Service        string
	PubKeyOperator string
	VotingAddress  string
	IsValid        bool
}

type QuorumInfo struct {
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
	DeletedMNs        []NodeInfo
	MnList            []NodeInfo
	DeletedQuorums    []QuorumInfo
	NewQuorums        []QuorumInfo
	MerkleRootMNList  string
	MerkleRootQuorums string
}
