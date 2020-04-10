package evo

import (
	"bytes"
	proto "dapi/protobuf"
	jsonRPC "dapi/protobuf/jsonRPC"
	"encoding/json"
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

type AddressSummaryResponse struct {
	AddrStr                 []string
	Balance                 float64
	BalanceSat              int64
	TotalReceived           float64
	TotalReceivedSat        int64
	TotalSent               float64
	TotalSentSat            int64
	UnconfirmedBalance      float64
	UnconfirmedBalanceSat   int64
	UnconfirmedTxApperances int
	UnconfirmedAppearances  int
	TxApperances            int
	TxAppearances           int
	Transactions            []string
}

type BlockHashResponse string

type BestBlockHashResponse string

type UTXOItem struct {
	Address     string
	Txid        string
	OutputIndex int
	Script      string
	Satoshis    int64
	Height      int
}

type UTXORequest struct {
	*LimitRange
	Addresses []string `json:"address"`
}

type UTXOResponse struct {
	TotalItems int
	From       int
	To         int
	FromHeight int
	ToHeight   int
	Items      []UTXOItem
}

type LimitRange struct {
	from       int
	to         int
	fromHeight int
	toHeight   int
}

type Block struct {
	hash   string
	height int
}

type NodeInfo struct {
	ProRegTxHash   string
	ConfirmedHash  string
	Service        string
	PubKeyOperator string
	VotingAddress  string
	IsValid        bool
}

type QuorumInfo struct {
	Version           int
	LLMQType          int
	QuorumHash        string
	SignersCount      int
	ValidMembersCount int
	QuorumPublicKey   string
}

type MnListDiffResponse struct {
	BaseBlockHash     string
	BlockHash         string
	CbTxMerkleTree    string
	CbTx              string
	DeletedMNs        []NodeInfo
	MnList            []NodeInfo
	DeletedQuorums    []QuorumInfo
	NewQuorums        []QuorumInfo
	MerkleRootMNList  string
	MerkleRootQuorums string
}

func (c *connection) GetBestBlockHash() (*jsonRPC.BestBlockHashResponse, error) {
	body := make(map[string][]string)

	response := new(jsonRPC.BestBlockHashResponse)
	err := c.requestJSON(true, true, jsonEndpointGetBestBlockHash, body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetBlockHash(height int) (*jsonRPC.BlockHashResponse, error) {
	body := make(map[string]int)
	body["height"] = height

	response := new(jsonRPC.BlockHashResponse)
	err := c.requestJSON(true, true, jsonEndpointGetBlockHash, body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetMnListDiff(baseBlockHash string, blockHash string) (*jsonRPC.MnListDiffResponse, error) {
	body := make(map[string]string)
	body["baseBlockHash"] = baseBlockHash
	body["blockHash"] = blockHash

	response := new(jsonRPC.MnListDiffResponse)
	err := c.requestJSON(true, true, jsonEndpointGetMnListDiff, body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetUTXO(addresses []string, limitRange *LimitRange) (*UTXOResponse, error) {
	body := new(UTXORequest)
	body.Addresses = addresses

	response := new(UTXOResponse)
	response.Items = make([]UTXOItem, 0)
	err := c.requestJSON(true, true, jsonEndpointGetUTXO, body, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetAddressSummary(addresses []string) (*AddressSummaryResponse, error) {
	body := make(map[string][]string)
	body["address"] = addresses

	response := new(AddressSummaryResponse)
	err := c.requestJSON(true, true, jsonEndpointGetAddressSummary, body, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetBlock(block *Block) (*proto.GetBlockResponse, error) {
	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewCoreClient(c.conn)
	request := new(proto.GetBlockRequest)

	if block.hash != "" {
		r := new(proto.GetBlockRequest_Hash)
		r.Hash = block.hash
		request.Block = r
	} else {
		r := new(proto.GetBlockRequest_Height)
		r.Height = uint32(block.height)
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

func (c *connection) GetNodeName() string {
	return c.name
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
