package main

import (
	"fmt"
	"log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

type TransactionFilter interface {
	QualifyTransaction(types.Block, types.Transaction) (bool, error)
}

type FilterChain struct {
	filters []TransactionFilter
}

func (fc FilterChain) QualifyTransaction(blk types.Block, txn types.Transaction) (bool, error) {
	for _, f := range fc.filters {
		r, err := f.QualifyTransaction(blk, txn)
		if err != nil {
			return false, fmt.Errorf("Error running filter: %v", err)
		}
		if !r {
			return false, nil
		}
	}
	log.Printf("Txn %h qualifies for FilterChain", txn.Hash())
	return true, nil
}

type SourceFilter struct {
	chainId big.Int
	from    common.Address
}

type DestinationFilter struct {
	chainId big.Int
	to      common.Address
}

type DateFilter struct {
	start   time.Time
	end     time.Time
}

func (f SourceFilter) QualifyTransaction(blk types.Block, txn types.Transaction) (bool, error) {
	msg, err := txn.AsMessage(types.NewEIP155Signer(txn.ChainId()))
	if err != nil {
		return false, fmt.Errorf("SourceFilter - Error getting txn as message: %v", err)
	}
	log.Printf("Checking txn from %s against filter from %s", msg.From().Hex(), f.from.Hex())
	if msg.From().Hex() == f.from.Hex() {
		log.Printf("Txn %h qualifies for SourceFilter", txn.Hash())
		return true, nil
	} else {
		return false, nil
	}
}

func (f DestinationFilter) QualifyTransaction(blk types.Block, txn types.Transaction) (bool, error) {
	msg, err := txn.AsMessage(types.NewEIP155Signer(txn.ChainId()))
	if err != nil {
		return false, fmt.Errorf("Error getting txn as message: %v", err)
	}
	if *(msg.To()) == f.to {
		return true, nil
	} else {
		return false, nil
	}
}

func (f DateFilter) QualifyTransaction(blk types.Block, txn types.Transaction) (bool, error) {
	t := time.Unix((int64)(blk.Time()), 0)
	if t.Equal(f.start) || t.Equal(f.end) || (t.After(f.start) && t.Before(f.end)) {
		return true, nil
	} else {
		return false, nil
	}
}
