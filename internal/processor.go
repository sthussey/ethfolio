package internal

import (
	"context"
	"log"
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

func ProcessBlocks(cfg Configuration){
	ctx, _ := context.WithCancel(context.Background())

	client, err := ethclient.DialContext(ctx, "http://127.0.0.1:8545")

	if err != nil {
		log.Fatalf("Error connecting: %v\n", err)
	}

	defer client.Close()

	var latestBlock *types.Block
	latestBlock, err = client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("Error getting latest block: %v\n", err)
	}

	fc := FilterChain{filters: make([]TransactionFilter, 0)}
	for _, a := range cfg.FromAccounts {
		log.Printf("Adding a selector for txn from %s", a)
		fc.filters = append(fc.filters, SourceFilter{chainId: *(big.NewInt(1)), from: common.HexToAddress(a)})
	}
	for _, a := range cfg.ToAccounts {
		log.Printf("Adding a selector for txn to %s", a)
		fc.filters = append(fc.filters, DestinationFilter{chainId: *(big.NewInt(1)), to: common.HexToAddress(a)})
	}
	fc.filters = append(fc.filters, DateFilter{start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), end: time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC)})

	for blockTime := time.Unix((int64)(latestBlock.Time()), 0); blockTime.Equal(cfg.BlocksSince) || blockTime.After(cfg.BlocksSince); {
		blockTxns := latestBlock.Transactions()
		for _, t := range blockTxns {
			qual, err := fc.QualifyTransaction(latestBlock, t)
			if err != nil {
				log.Printf("Error qualifying txn %h: %v", t.Hash(), err)
			} else if qual {
				log.Printf("Transaction %h qualifies for lot.", t.Hash())
			}
		}
		latestBlock, err = client.BlockByNumber(ctx, big.NewInt(latestBlock.Number().Int64()-1))
		if err != nil {
			log.Fatalf("Error finding parent block: %v\n", err)
		}
		blockTime = time.Unix((int64)(latestBlock.Time()), 0)
	}
}
