package internal

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"time"
)

type TransactionFilter interface {
	QualifyTransaction(EthBlock, EthTransaction) (bool, error)
}

type EthBlock interface {
	Time()	uint64
}

type EthTransaction interface {
	ChainId()	*big.Int
	Hash()		common.Hash
    AsMessage(s types.Signer) (types.Message, error)
}

type EthMessage interface {
	To()	*common.Address
	From()	common.Address
}

type FilterChain struct {
	or      bool
	filters []TransactionFilter
}

func (fc FilterChain) QualifyTransaction(blk EthBlock, txn EthTransaction) (bool, error) {
	result := false
	result_set := false
	for _, f := range fc.filters {
		r, err := f.QualifyTransaction(blk, txn)
		if !fc.or && err != nil {
			return false, fmt.Errorf("Error running filter: %v", err)
		} else if err != nil {
			fmt.Printf("Error running filter: %v", err)
		}
		if !result_set {
			result = r
			result_set = true
			continue
		}
		if fc.or {
			result = result || r
		} else {
			result = result && r
			// Fast exit for and-ed chain and a false subfilter
			if !result {
				return result, nil
			}
		}
	}
	if result {
		log.Printf("Txn %h qualifies for FilterChain", txn.Hash())
	}
	return result, nil
}

type SourceFilter struct {
	chainVersion	int
	chainId			big.Int
	from			common.Address
}

func NewSourceFilter(ver int, id big.Int, addr string) (*SourceFilter, error) {
	if ver != 1 {
		return nil, fmt.Errorf("NewSourceFilter: Chain version %d not supported.", ver)
	}
	from := common.HexToAddress(addr)
	return &SourceFilter{chainVersion: ver, chainId: id, from: from}, nil
}

func compMessageSource(msg EthMessage, addr common.Address)	bool {
	return msg.From().Hex() == addr.Hex()
}

func (f SourceFilter) QualifyTransaction(blk EthBlock, txn EthTransaction) (bool, error) {
	msg, err := txn.AsMessage(types.NewEIP155Signer(txn.ChainId()))
	if err != nil {
		return false, fmt.Errorf("SourceFilter - Error getting txn as message: %v", err)
	}
	if compMessageSource(msg, f.from) {
		log.Printf("Txn %h qualifies for SourceFilter", txn.Hash())
		return true, nil
	} else {
		return false, nil
	}
}
type DestinationFilter struct {
	chainVersion	int
	chainId			big.Int
	to				common.Address
}

func NewDestinationFilter(ver int, id big.Int, addr string) (*DestinationFilter, error) {
	if ver != 1 {
		return nil, fmt.Errorf("NewDestinationFilter: Chain version %d not supported.", ver)
	}
	to := common.HexToAddress(addr)
	return &DestinationFilter{chainVersion: ver, chainId: id, to: to}, nil
}

func compMessageDest(msg EthMessage, addr common.Address) bool {
	if msg.To() == nil {
		return false
	}
	return msg.To().Hex() == addr.Hex()
}


func (f DestinationFilter) QualifyTransaction(blk EthBlock, txn EthTransaction) (bool, error) {
	msg, err := txn.AsMessage(types.NewEIP155Signer(txn.ChainId()))
	if err != nil {
		return false, fmt.Errorf("DestinationFilter - Error getting txn as message: %v", err)
	}
	if compMessageDest(msg, f.to) {
		log.Printf("Txn %h qualifies for DestinationFilter", txn.Hash())
		return true, nil
	} else {
		return false, nil
	}
}


type DateFilter struct {
	start time.Time
	end   time.Time
}

func NewDateFilter(start time.Time, end time.Time) (*DateFilter, error) {
	if start.After(end) {
		return nil, fmt.Errorf("Error creating DateFilter: start must be before end")
	}
	return &DateFilter{start: start, end: end}, nil
}

func (f DateFilter) QualifyTransaction(blk EthBlock, txn EthTransaction) (bool, error) {
	t := time.Unix((int64)(blk.Time()), 0)
	if t.Equal(f.start) || t.Equal(f.end) || (t.After(f.start) && t.Before(f.end)) {
		return true, nil
	} else {
		return false, nil
	}
}
