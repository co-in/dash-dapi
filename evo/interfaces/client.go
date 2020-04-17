package interfaces

import (
	proto "github.com/co-in/dash-dapi/evo/protobuf"
	"github.com/co-in/dash-dapi/evo/structures"
)

type IClient interface {
	AddNode(hostname string, fraud int) error
	SelectRandomNode() (IConnection, error)
	SelectNode(hostname string) (IConnection, error)
}

type IConnection interface {
	ILayer1
	ILayer2
	GetNodeName() string
	CheckAvailability() bool
	Remove()
}

type ILayer1 interface {
	ILayer1JSON
	ILayer1GRPC
}

type ILayer1JSON interface {
	GetBestBlockHash() (*structures.BestBlockHashResponse, error)
	GetBlockHash(height int) (*structures.BlockHashResponse, error)
	GetMnListDiff(baseBlockHash string, blockHash string) (*structures.MnListDiffResponse, error)
	GetAddressSummary(addresses []string) (*structures.AddressSummaryResponse, error)
	GetUTXO(request structures.UTXORequest) (*structures.UTXOResponse, error)
}

type ILayer1GRPC interface {
	GetStatus() (*proto.GetStatusResponse, error)
	GetBlock(block structures.BlockRequest) (*proto.GetBlockResponse, error)
	GetTransaction(id string) (*proto.GetTransactionResponse, error)
	SendTransaction(data []byte, allowHighFees bool, bypassLimits bool) (*proto.SendTransactionResponse, error)
	SubscribeToTransactionsWithProofs(params structures.SubscribeToTransactionsWithProofsRequest) (proto.TransactionsFilterStream_SubscribeToTransactionsWithProofsClient, error)
}

type ILayer2 interface {
	ApplyStateTransition(stateTransition []byte) (*proto.ApplyStateTransitionResponse, error)
	GetIdentity(id string) (*proto.GetIdentityResponse, error)
	GetDataContract(id string) (*proto.GetDataContractResponse, error)
	GetDocuments(dataContractId string, documentType string, filter structures.GetDocumentsRequest) (*proto.GetDocumentsResponse, error)
}
