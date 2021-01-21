package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

	client, err := ethclient.DialContext(ctx, "http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Error connecting: %v\n", err)
	}

	// filterAddr := "0x976813864377495515FBB0c2CdE1cbAC897fE52a"

	filterAddr := "0x3561e7113da3ec62b52c050d24f1ee000760f885"
	oldestBlock := time.Date(2021, time.January, 19, 0, 0, 0, 0, time.UTC)
	defer client.Close()

	var latestBlock *types.Block
	latestBlock, err = client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("Error getting latest block: %v\n", err)
	}

	fc := FilterChain{ filters: []TransactionFilter{SourceFilter{ chainId: *(big.NewInt(1)), from: common.HexToAddress(filterAddr) }}}

	for blockTime := time.Unix((int64)(latestBlock.Time()), 0); blockTime.Equal(oldestBlock) || blockTime.After(oldestBlock); {
		log.Printf("Qualifying block %s.", latestBlock.Number().String())
		blockTxns := latestBlock.Transactions()
		log.Printf("Checking %d transactions in block.", len(blockTxns))
		for _, t := range blockTxns {
			qual, err := fc.QualifyTransaction(*latestBlock, *t)
			if err != nil {
				log.Printf("Error qualifying txn %h: %v", t.Hash(), err)
			} else if qual {
				log.Printf("Transaction %h qualifies for lot.", t.Hash())
			}
		}
		latestBlock, err = client.BlockByNumber(ctx, big.NewInt(latestBlock.Number().Int64() - 1))
		if err != nil {
			log.Fatalf("Error finding parent block: %v\n", err)
		}
		blockTime = time.Unix((int64)(latestBlock.Time()), 0) 
	}

	os.Exit(0)
}
