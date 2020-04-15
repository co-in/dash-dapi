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
	//GetAddressSummary(addresses []string) (*evo.AddressSummaryResponse, error)
	//GetUTXO(addresses []string, limitRange *evo.LimitRange) (*evo.UTXOResponse, error)
}

type ILayer1GRPC interface {
	GetStatus() (*proto.GetStatusResponse, error)
	GetTransaction(id string) (*proto.GetTransactionResponse, error)
	SendTransaction(data []byte, allowHighFees bool, bypassLimits bool) (*proto.SendTransactionResponse, error)
}

type ILayer2 interface {
}
