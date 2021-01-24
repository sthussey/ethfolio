package internal

import (
	"context"
	"log"
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

func ProcessBlocks(cfg Configuration, metrics MetricDefinitions){
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

	fc := FilterChain{or: false, filters: make([]TransactionFilter, 0)}
	for _, a := range cfg.FromAccounts {
		log.Printf("Adding a Source selector for txn from %s", a)
		sf, err := NewSourceFilter(1, *(big.NewInt(1)), a)
		if err != nil {
			log.Fatalf("Error creating SourceFilter: %v", err)
		}
		fc.filters = append(fc.filters, sf)
	}
	for _, a := range cfg.ToAccounts {
		log.Printf("Adding a Destination selector for txn to %s", a)
		df, err := NewDestinationFilter(1, *(big.NewInt(1)), a)
		if err != nil {
			log.Fatalf("Error creating DestinationFilter: %v", err)
		}
		fc.filters = append(fc.filters, df)
	}
	fc.filters = append(fc.filters, DateFilter{start: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), end: time.Date(2021, time.December, 31, 23, 59, 59, 0, time.UTC)})

	for blockTime := time.Unix((int64)(latestBlock.Time()), 0); blockTime.Equal(cfg.BlocksSince) || blockTime.After(cfg.BlocksSince); {
		blockTxns := latestBlock.Transactions()
		for _, t := range blockTxns {
			qual, err := fc.QualifyTransaction(latestBlock, t)
			metrics.txns_processed.Inc()
			if err != nil {
				log.Printf("Error qualifying txn %h: %v", t.Hash(), err)
			} else if qual {
				log.Printf("Transaction %h qualifies for lot.", t.Hash())
			}
		}
		metrics.blocks_processed.Inc()
		metrics.last_block_completed.Set(float64(latestBlock.Number().Int64()))
		latestBlock, err = client.BlockByNumber(ctx, big.NewInt(latestBlock.Number().Int64()-1))
		if err != nil {
			log.Fatalf("Error finding parent block: %v\n", err)
		}
		blockTime = time.Unix((int64)(latestBlock.Time()), 0)
	}
}
