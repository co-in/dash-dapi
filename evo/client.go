package evo

import (
	"context"
	"crypto/rand"
	"dapi/interfaces"
	"errors"
	"google.golang.org/grpc"
	"log"
	"math/big"
	"strconv"
	"sync"
)

type connection struct {
	*client
	name  string
	fraud int
	conn  *grpc.ClientConn
}

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

func NewClient(log *log.Logger, verboseLevel int, nodeAddress string, jsonRpcPort int, gRpcPort int) (*client, error) {
	c := &client{
		Logger:         log,
		ctx:            context.Background(),
		evoJsonRpcPort: strconv.Itoa(jsonRpcPort),
		evoGRpcPort:    strconv.Itoa(gRpcPort),
		verboseLevel:   verboseLevel,
		connections:    make(map[string]*connection),
	}

	err := c.AddNode(nodeAddress)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) AddNode(nodeAddress string) error {
	if _, ok := c.connections[nodeAddress]; ok {
		return nil
	}

	c.connectionKeys = append(c.connectionKeys, nodeAddress)
	c.connections[nodeAddress] = &connection{
		client: c,
		name:   nodeAddress,
		//conn:   conn,
	}

	return nil
}

func (c *connection) LazyConnection() error {
	if c.conn != nil {
		return nil
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())

	//TODO How to get Cert of Node?
	if false {
		//serverHostOverride := "evonet.thephez.com"
		//caFile := "PEM file"
		//credential, err := credentials.NewClientTLSFromFile(caFile, serverHostOverride)
		//
		//if err != nil {
		//	return err
		//}
		//
		//opts = append(opts, grpc.WithTransportCredentials(credential))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(c.name+":"+c.evoGRpcPort, opts...)

	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *client) SelectRandomNode() (interfaces.IConnection, error) {
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
		return nil, errors.New("node dont exist")
	}

	return conn, nil
}

func (c *connection) CheckAvailability() bool {
	body := make(map[string][]string)
	response := new(BestBlockHashResponse)

	err := c.requestJSON(false, false, jsonEndpointGetBestBlockHash, body, response)

	if err != nil || *response == "" {
		return false
	}

	return true
}
