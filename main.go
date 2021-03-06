package main

import (
	"fmt"
	"github.com/co-in/dash-dapi/db"
	"github.com/co-in/dash-dapi/db/jsonFile"
	"github.com/co-in/dash-dapi/evo"
	"github.com/co-in/dash-dapi/evo/interfaces"
	"github.com/co-in/dash-dapi/evo/structures"
	"log"
	"math"
	"os"
	"regexp"
	"sync"
	"time"
)

func discoveryNewEvoNodes(dAPI interfaces.IClient, logger *log.Logger, dbProvider db.IDatabase, evoNodes []string) {
	node, err := dAPI.SelectRandomNode()

	if err != nil {
		logger.Fatalln(err)
	}

	lastBlockHash, err := node.GetBestBlockHash()

	if err != nil {
		logger.Fatalln(err)
	}

	dbLastBlockHash := dbProvider.GetCurrentBlockHash()

	//Sync MasterNode list
	if string(*lastBlockHash) != dbLastBlockHash {
		//First Run
		if dbLastBlockHash == "" {
			baseBlockHash, err := node.GetBlockHash(0)

			if err != nil {
				logger.Fatalln(err)
			}

			dbProvider.SetCurrentBlockHash(string(*baseBlockHash))
		}

		mnList, err := node.GetMnListDiff(dbProvider.GetCurrentBlockHash(), string(*lastBlockHash))

		if err != nil {
			logger.Fatalln(err)
		}

		//Discovery new EVO nodes
		wg := new(sync.WaitGroup)
		wg.Add(len(mnList.MnList))
		re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
		m := new(sync.Mutex)

		for _, v := range mnList.MnList {
			//Async connection to EvoNodes
			go func(v structures.NodeInfoResponse) {
				var tempNode interfaces.IConnection

				defer func() {
					if tempNode != nil {
						tempNode.Remove()
					}

					wg.Done()
				}()

				ip := re.Find([]byte(v.Service))
				ipStr := string(ip)
				err = dAPI.AddNode(ipStr, 0)

				if err != nil {
					return
				}

				tempNode, err = dAPI.SelectNode(ipStr)

				if err != nil {
					return
				}

				if tempNode.CheckAvailability() {
					m.Lock()
					evoNodes = append(evoNodes, ipStr)
					m.Unlock()
				}
			}(v)
		}

		wg.Wait()

		dbProvider.SetEvoNodes(evoNodes)
		dbProvider.SetCurrentBlockHash(string(*lastBlockHash))

		err = dbProvider.Save()

		if err != nil {
			logger.Fatalln(err)
		}
	}
}

func getStatus(logger *log.Logger, node interfaces.IConnection) {
	result, err := node.GetStatus()

	if err != nil {
		logger.Fatalln(err)
	}

	fmt.Printf("Status of node:\t\t%s\nNetwork:\t\t%s\nCoreVersion:\t\t%d\nBlocks:\t\t\t%d\nRelayFee:\t\t%f\nConnections:\t\t%d\nNetworkDifficulty:\t%f\n\n",
		node.GetNodeName(),
		result.Network,
		result.CoreVersion,
		result.Blocks,
		result.RelayFee,
		result.Connections,
		result.Difficulty,
	)
}

func getTransactionStream(logger *log.Logger, node interfaces.IConnection) {
	g, err := node.SubscribeToTransactionsWithProofs(structures.SubscribeToTransactionsWithProofsRequest{
		BloomFilter: structures.BloomFilterRequest{
			HashFunc: 11,
			Data:     []byte{0xB5, 0x0F},
		},
	})

	if err != nil {
		logger.Fatalln(err)
	}

	for {
		r, err := g.Recv()

		if err != nil {
			logger.Fatalln(err)
		}

		transactions := r.GetRawTransactions()

		if transactions != nil {
			fmt.Println("Transactions:")

			for index, transaction := range transactions.GetTransactions() {
				fmt.Printf("%d:\t%0X\n", index, transaction)
			}
		}

		merkleBlock := r.GetRawMerkleBlock()

		if merkleBlock != nil {
			fmt.Printf("MerkleBlock: %0X\n", merkleBlock)
		}

		instantSendLock := r.GetInstantSendLockMessages()

		if instantSendLock != nil {
			fmt.Println("InstantSendLock:")
			fmt.Println(instantSendLock)
		}

		time.Sleep(1 * time.Second)
	}
}

func getIdentity(logger *log.Logger, node interfaces.IConnection) {
	r, err := node.GetIdentity("At44pvrZXLwjbJp415E2kjav49goGosRF3SB1WW1QJoG")
	//r, err := node.GetIdentity("A6AJAfRJyKuNoNvt33ygYfYh6OIYA8tF1s2BQcRA9RNg")

	if err != nil {
		logger.Fatalln(err)
	}

	i := r.GetIdentity()

	//TODO Parse identity
	fmt.Printf("Identity: %s\n", i)
}

func getDataContract(logger *log.Logger, node interfaces.IConnection) {
	r, err := node.GetDataContract("77w8Xqn25HwJhjodrHW133aXhjuTsTv9ozQaYpSHACE3")

	if err != nil {
		logger.Fatalln(err)
	}

	c := r.GetDataContract()

	//TODO Parse contract
	fmt.Printf("DataContract: %s\n", c)
}

func getDocuments(logger *log.Logger, node interfaces.IConnection) {
	limit := 1

	r, err := node.GetDocuments(
		"77w8Xqn25HwJhjodrHW133aXhjuTsTv9ozQaYpSHACE3",
		"domain", structures.GetDocumentsRequest{
			Limit: &limit,
		},
	)

	if err != nil {
		logger.Fatalln(err)
	}

	documents := r.GetDocuments()

	fmt.Println("Documents:")

	for i, d := range documents {
		//TODO Parse documents
		fmt.Printf("%d\t:%s\n", i, d)
	}
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	dbProvider := jsonFile.NewDB("db.json")

	err := dbProvider.Load()

	if err != nil {
		logger.Fatalln(err)
	}

	jsonRpcPort := dbProvider.GetEvoJsonRpcPort()

	if jsonRpcPort < 1 || jsonRpcPort > math.MaxUint16 {
		logger.Fatalln("Invalid or missing evo_json_rpc_port in db")
	}

	gRpcPort := dbProvider.GetEvoGRpcPort()

	if gRpcPort < 1 || gRpcPort > math.MaxUint16 {
		logger.Fatalln("Invalid or missing evo_grpc_port in db")
	}

	evoNodes := dbProvider.GetEvoNodes()

	if len(evoNodes) < 1 {
		logger.Fatalln("Empty evo_nodes in db")
	}

	dAPI, err := evo.NewClient(logger, 1, evoNodes[0], jsonRpcPort, gRpcPort)

	if err != nil {
		logger.Fatalln(err)
	}

	//Apply evoNodes
	for _, v := range dbProvider.GetEvoNodes() {
		err = dAPI.AddNode(v, 0)
	}

	//At first Run Discovery other nodes
	//if len(evoNodes) == 1 {
	//	discoveryNewEvoNodes(dAPI, logger, dbProvider, evoNodes)
	//}

	node, err := dAPI.SelectRandomNode()

	if err != nil {
		logger.Fatalln(err)
	}

	getStatus(logger, node)
	getIdentity(logger, node)
	getDataContract(logger, node)
	getDocuments(logger, node)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go getTransactionStream(logger, node)
	wg.Wait() //Wait forever
}
