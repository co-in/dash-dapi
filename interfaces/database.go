package interfaces

type IBaseDatabase interface {
	GetEvoJsonRpcPort() int
	GetEvoGRpcPort() int
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
