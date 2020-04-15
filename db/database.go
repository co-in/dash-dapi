package db

type baseDatabase struct {
	EvoJsonRpcPort   uint16   `json:"evo_json_rpc_port"`
	EvoGRpcPort      uint16   `json:"evo_grpc_port"`
	EvoNodes         []string `json:"evo_nodes"`
	CurrentBlockHash string   `json:"current_block_hash"`
}

type IBaseDatabase interface {
	GetEvoJsonRpcPort() uint16
	GetEvoGRpcPort() uint16
	GetEvoNodes() []string
	SetEvoNodes(nodes []string)
	GetCurrentBlockHash() string
	SetCurrentBlockHash(hash string)
}

type IDatabase interface {
	IBaseDatabase
	Load() error
	Save() error
}

func NewBaseDatabase() *baseDatabase {
	return &baseDatabase{}
}

func (d *baseDatabase) GetEvoNodes() []string {
	return d.EvoNodes
}

func (d *baseDatabase) GetEvoJsonRpcPort() uint16 {
	return d.EvoJsonRpcPort
}

func (d *baseDatabase) GetEvoGRpcPort() uint16 {
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
