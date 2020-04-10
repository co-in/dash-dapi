package interfaces

import (
	proto "dapi/protobuf"
	jsonRPC "dapi/protobuf/jsonRPC"
)

type IClient interface {
	AddNode(hostname string) error
	SelectRandomNode() (IConnection, error)
	SelectNode(hostname string) (IConnection, error)
}

type IConnection interface {
	ILayer1
	ILayer2
	GetNodeName() string
	CheckAvailability() bool
}

type ILayer1 interface {
	ILayer1JSON
	ILayer1GRPC
}

type ILayer1JSON interface {
	GetBestBlockHash() (*jsonRPC.BestBlockHashResponse, error)
	GetBlockHash(height int) (*jsonRPC.BlockHashResponse, error)
	GetMnListDiff(baseBlockHash string, blockHash string) (*jsonRPC.MnListDiffResponse, error)
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
