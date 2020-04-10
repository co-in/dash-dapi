package main

import (
	"dapi/db"
	"dapi/db/jsonFile"
	"dapi/evo"
	"dapi/interfaces"
	layer1 "dapi/protobuf/jsonRPC"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

func discoveryNewEvoNodes(dAPI interfaces.IClient, logger *log.Logger, dbProvider interfaces.IDatabase, evoNodes []string) {
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
	if lastBlockHash.String() != dbLastBlockHash {
		//First Run
		if dbLastBlockHash == "" {
			baseBlockHash, err := node.GetBlockHash(0)

			if err != nil {
				logger.Fatalln(err)
			}

			dbProvider.SetCurrentBlockHash(baseBlockHash.String())
		}

		mnList, err := node.GetMnListDiff(dbProvider.GetCurrentBlockHash(), lastBlockHash.String())

		if err != nil {
			logger.Fatalln(err)
		}

		//Discovery new EVO nodes
		wg := new(sync.WaitGroup)
		wg.Add(len(mnList.MnList))
		nodes := make([]db.ExtNodeInfo, 0)
		re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
		m := new(sync.Mutex)

		for _, v := range mnList.MnList {
			//Async connection to EvoNodes
			go func(v *layer1.NodeInfo) {
				n := db.ExtNodeInfo{
					NodeInfo: v,
				}

				ip := re.Find([]byte(v.Service))
				ipStr := string(ip)
				node, err = dAPI.SelectNode(ipStr)

				if err != nil {
					wg.Done()

					return
				}

				if node.CheckAvailability() {
					n.IsEvoNode = true
					m.Lock()
					evoNodes = append(evoNodes, ipStr)
					m.Unlock()
				}

				m.Lock()
				nodes = append(nodes, n)
				m.Unlock()

				wg.Done()
			}(v)
		}

		wg.Wait()

		keys := make(map[string]bool)
		var list []string

		for _, entry := range evoNodes {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				list = append(list, entry)
			}
		}

		//dbProvider.SetNodes(nodes)
		dbProvider.SetEvoNodes(evoNodes)
		dbProvider.SetCurrentBlockHash(lastBlockHash.String())

		err = dbProvider.Save()

		if err != nil {
			logger.Fatalln(err)
		}
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

	if jsonRpcPort < 1 || jsonRpcPort > 65535 {
		logger.Fatalln("Invalid or missing evo_json_rpc_port in db")
	}

	gRpcPort := dbProvider.GetEvoGRpcPort()

	if gRpcPort < 1 || gRpcPort > 65535 {
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

	if len(evoNodes) == 1 {
		discoveryNewEvoNodes(dAPI, logger, dbProvider, evoNodes)
	}

	//Apply evoNodes
	for _, v := range dbProvider.GetEvoNodes() {
		err = dAPI.AddNode(v)
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
