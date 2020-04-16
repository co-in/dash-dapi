package evo

import (
	"github.com/co-in/dash-dapi/evo/structures"
	"google.golang.org/grpc"
	"sort"
)

type connection struct {
	*client
	name  string
	fraud int
	conn  *grpc.ClientConn
}

func (c *connection) CheckAvailability() bool {
	body := make(map[string][]string)
	response := new(structures.BestBlockHashResponse)

	err := c.requestJSON(false, false, jsonEndpointGetBestBlockHash, body, response)

	//Check JSON RPC
	if err != nil || *response == "" {
		return false
	}

	status, err := c.GetStatus()

	//Check gRPC
	if err != nil || status == nil || status.Connections < 2 {
		return false
	}

	return true
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

func (c *connection) GetNodeName() string {
	return c.name
}

func (c *connection) Remove() {
	c.client.Lock()
	defer c.client.Unlock()

	if c.client.connections[c.name].conn != nil {
		_ = c.client.connections[c.name].conn.Close()
	}

	i := sort.SearchStrings(c.client.connectionKeys, c.name)
	connectionNewLen := len(c.client.connectionKeys) - 1

	if connectionNewLen+1 == i {
		return
	}

	c.client.connectionKeys[i] = c.client.connectionKeys[connectionNewLen]
	c.client.connectionKeys[connectionNewLen] = ""
	c.client.connectionKeys = c.client.connectionKeys[:connectionNewLen]
	delete(c.client.connections, c.name)
}
