package db

import (
	layer1 "dapi/protobuf/jsonRPC"
)

type ExtNodeInfo struct {
	*layer1.NodeInfo
	IsEvoNode bool
}

type baseDatabase struct {
	Nodes            []ExtNodeInfo `json:"nodes"`
	EvoJsonRpcPort   int           `json:"evo_json_rpc_port"`
	EvoGRpcPort      int           `json:"evo_grpc_port"`
	EvoNodes         []string      `json:"evo_nodes"`
	CurrentBlockHash string        `json:"current_block_hash"`
}

func NewBaseDatabase() *baseDatabase {
	return &baseDatabase{}
}

func (d *baseDatabase) GetEvoNodes() []string {
	return d.EvoNodes
}

func (d *baseDatabase) GetEvoJsonRpcPort() int {
	return d.EvoJsonRpcPort
}

func (d *baseDatabase) GetEvoGRpcPort() int {
	return d.EvoGRpcPort
}

func (d *baseDatabase) SetEvoNodes(nodes []string) {
	d.EvoNodes = nodes
}

func (d *baseDatabase) GetCurrentBlockHash() string {
	return d.CurrentBlockHash
}

func (d *baseDatabase) SetCurrentBlockHash(hash string) {
	d.CurrentBlockHash = hash
}
