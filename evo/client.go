package evo

import (
	"context"
	"crypto/rand"
	"errors"
	"github.com/co-in/dash-dapi/evo/interfaces"
	"log"
	"math/big"
	"strconv"
	"sync"
)

type client struct {
	*log.Logger
	evoJsonRpcPort string
	evoGRpcPort    string
	connections    map[string]*connection
	connectionKeys []string
	ctx            context.Context
	sync.Mutex
	id           int
	verboseLevel int
}

func NewClient(log *log.Logger, verboseLevel int, nodeAddress string, jsonRpcPort uint16, gRpcPort uint16) (*client, error) {
	c := &client{
		Logger:         log,
		verboseLevel:   verboseLevel,
		ctx:            context.Background(),
		evoJsonRpcPort: strconv.Itoa(int(jsonRpcPort)),
		evoGRpcPort:    strconv.Itoa(int(gRpcPort)),
		connections:    make(map[string]*connection),
	}

	err := c.AddNode(nodeAddress, 0)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) AddNode(nodeAddress string, fraud int) error {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.connections[nodeAddress]; ok {
		return nil
	}

	c.connectionKeys = append(c.connectionKeys, nodeAddress)
	c.connections[nodeAddress] = &connection{
		client: c,
		name:   nodeAddress,
		fraud:  fraud,
	}

	return nil
}

func (c *client) SelectRandomNode() (interfaces.IConnection, error) {
	c.Lock()
	defer c.Unlock()

	//TODO Order By Fraud
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(c.connectionKeys))))

	if err != nil {
		return nil, err
	}

	return c.connections[c.connectionKeys[int(randIndex.Int64())]], nil
}

func (c *client) SelectNode(address string) (interfaces.IConnection, error) {
	var conn *connection
	var ok bool

	if conn, ok = c.connections[address]; !ok {
		return nil, errors.New("the node does not exist")
	}

	return conn, nil
}
