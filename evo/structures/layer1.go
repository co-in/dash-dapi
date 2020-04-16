package structures

import (
	"encoding/json"
	"errors"
)

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

type UTXORequest struct {
	From       *int     `json:"from,omitempty"`
	To         *int     `json:"to,omitempty"`
	FromHeight *int     `json:"fromHeight,omitempty"`
	ToHeight   *int     `json:"toHeight,omitempty"`
	Addresses  []string `json:"address"`
}

func (s UTXORequest) MarshalJSON() ([]byte, error) {
	if len(s.Addresses) == 0 {
		return nil, errors.New("empty field Addresses")
	}

	return json.Marshal(s)
}

type BlockRequest struct {
	Hash   *string `json:"hash,omitempty"`
	Height *int    `json:"height,omitempty"`
}

func (s BlockRequest) MarshalJSON() ([]byte, error) {
	if s.Hash == nil && s.Height == nil {
		return nil, errors.New("required one of fields (Hash, Height)")
	}

	return json.Marshal(s)
}

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
