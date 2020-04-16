package evo

import (
	"bytes"
	"encoding/json"
	"errors"
	proto "github.com/co-in/dash-dapi/evo/protobuf"
	"github.com/co-in/dash-dapi/evo/structures"

	"fmt"
	"io/ioutil"
	"net/http"
)

type jsonRPCRequest struct {
	messageId int
	method    string
	params    interface{}
}

type jsonRPCResponse struct {
	MessageId int             `json:"id"`
	Result    json.RawMessage `json:"result"`
	Error     *struct {
		Code    int
		Message string
	}
}

const (
	jsonEndpointGetAddressSummary = "getAddressSummary"
	jsonEndpointGetBestBlockHash  = "getBestBlockHash"
	jsonEndpointGetBlockHash      = "getBlockHash"
	jsonEndpointGetMnListDiff     = "getMnListDiff"
	jsonEndpointGetUTXO           = "getUTXO"
)

func (r jsonRPCRequest) Marshal() ([]byte, error) {
	payload := make(map[string]interface{})
	payload["jsonrpc"] = "2.0"
	payload["method"] = r.method
	payload["id"] = r.messageId
	payload["params"] = r.params

	b, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *connection) GetBestBlockHash() (*structures.BestBlockHashResponse, error) {
	body := make(map[string][]string)

	response := new(structures.BestBlockHashResponse)
	err := c.requestJSON(true, true, jsonEndpointGetBestBlockHash, body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetBlockHash(height int) (*structures.BlockHashResponse, error) {
	body := make(map[string]int)
	body["height"] = height

	response := new(structures.BlockHashResponse)
	err := c.requestJSON(true, true, jsonEndpointGetBlockHash, body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetMnListDiff(baseBlockHash string, blockHash string) (*structures.MnListDiffResponse, error) {
	body := make(map[string]string)
	body["baseBlockHash"] = baseBlockHash
	body["blockHash"] = blockHash

	response := new(structures.MnListDiffResponse)
	err := c.requestJSON(true, true, jsonEndpointGetMnListDiff, body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetUTXO(request structures.UTXORequest) (*structures.UTXOResponse, error) {
	response := new(structures.UTXOResponse)
	err := c.requestJSON(true, true, jsonEndpointGetUTXO, request, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetAddressSummary(addresses []string) (*structures.AddressSummaryResponse, error) {
	body := make(map[string][]string)
	body["address"] = addresses

	response := new(structures.AddressSummaryResponse)
	err := c.requestJSON(true, true, jsonEndpointGetAddressSummary, body, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetBlock(block structures.BlockRequest) (*proto.GetBlockResponse, error) {
	if block.Hash == nil && block.Height == nil {
		return nil, errors.New("required one of fields (Hash, Height)")
	}

	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewCoreClient(c.conn)
	request := new(proto.GetBlockRequest)

	if block.Hash == nil {
		r := new(proto.GetBlockRequest_Hash)
		r.Hash = *block.Hash
		request.Block = r
	} else {
		r := new(proto.GetBlockRequest_Height)
		r.Height = uint32(*block.Height)
		request.Block = r
	}

	response, err := layer1.GetBlock(c.ctx, request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetStatus() (*proto.GetStatusResponse, error) {
	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewCoreClient(c.conn)
	request := new(proto.GetStatusRequest)
	response, err := layer1.GetStatus(c.ctx, request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetTransaction(id string) (*proto.GetTransactionResponse, error) {
	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewCoreClient(c.conn)
	request := new(proto.GetTransactionRequest)
	request.Id = id

	response, err := layer1.GetTransaction(c.ctx, request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) SendTransaction(data []byte, allowHighFees bool, bypassLimits bool) (*proto.SendTransactionResponse, error) {
	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewCoreClient(c.conn)
	request := &proto.SendTransactionRequest{
		Transaction:   data,
		AllowHighFees: allowHighFees,
		BypassLimits:  bypassLimits,
	}

	response, err := layer1.SendTransaction(c.ctx, request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) requestJSON(withVerbose bool, withFraud bool, method string, params interface{}, response interface{}) error {
	c.Lock()
	c.id++
	id := c.id
	c.Unlock()

	nodeAddressURL := "http://" + c.name + ":" + c.evoJsonRpcPort

	request, err := (jsonRPCRequest{
		messageId: id,
		method:    method,
		params:    params,
	}).Marshal()

	if err != nil {
		return err
	}

	r := bytes.NewReader(request)

	if withVerbose && c.verboseLevel > 2 {
		c.Printf("Send JSON-RPC packet #%d to %s: %s\n", c.id, nodeAddressURL, request)
	}

	resp, err := http.Post(nodeAddressURL, "application/json", r)

	if err != nil {
		if withFraud {
			c.fraud++
		}

		if withVerbose && c.verboseLevel > 1 {
			c.Printf("Increase Fraud score: %s to %d. (Reason:%s)\n",
				nodeAddressURL,
				c.fraud,
				err,
			)
		}

		return err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = resp.Body.Close()

	if err != nil {
		if withVerbose && c.verboseLevel > 1 {
			c.Println(err)
		}
	}

	jsonRPCResponse := new(jsonRPCResponse)
	err = json.Unmarshal(data, jsonRPCResponse)

	if err != nil {
		return err
	}

	if jsonRPCResponse.Error != nil {
		return fmt.Errorf("[ERROR] Recv JSON-RPC packet #%d from %s: [%d] %s\n",
			jsonRPCResponse.MessageId,
			c.name,
			jsonRPCResponse.Error.Code,
			jsonRPCResponse.Error.Message,
		)
	}

	if withVerbose && c.verboseLevel > 2 {
		c.Printf("Recv JSON-RPC packet #%d from %s: %s\n", jsonRPCResponse.MessageId, c.name, jsonRPCResponse.Result)
	}

	err = json.Unmarshal(jsonRPCResponse.Result, response)

	if err != nil {
		return err
	}

	return nil
}
