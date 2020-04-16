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

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	dbProvider := jsonFile.NewDB("_db.json")

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

	//At first Run Discovery other nodes
	if len(evoNodes) == 1 {
		discoveryNewEvoNodes(dAPI, logger, dbProvider, evoNodes)
	}

	//Apply evoNodes
	for _, v := range dbProvider.GetEvoNodes() {
		err = dAPI.AddNode(v, 0)
	}

	node, err := dAPI.SelectRandomNode()

	if err != nil {
		logger.Fatalln(err)
	}

	result, err := node.GetStatus()

	if err != nil {
		logger.Fatalln(err)
	}

	fmt.Printf("Status of node:\t\t%s\nNetwork:\t\t%s\nCoreVersion:\t\t%d\nBlocks:\t\t\t%d\nRelayFee:\t\t%f\nConnections:\t\t%d\nNetworkDifficulty:\t%f\n",
		node.GetNodeName(),
		result.Network,
		result.CoreVersion,
		result.Blocks,
		result.RelayFee,
		result.Connections,
		result.Difficulty,
	)
}
