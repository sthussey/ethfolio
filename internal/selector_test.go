package internal

import (
	"math/big"
	"testing"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TestBlock struct {
	testTime	time.Time
}

func (b TestBlock) Time() uint64 {
	return (uint64)(b.testTime.Unix())
}

type TestTransaction struct {
	hash	common.Hash
	msg		types.Message
}

func (t TestTransaction) ChainId()	*big.Int {
	return big.NewInt(1)
}

func (t TestTransaction) Hash()	common.Hash {
	return t.hash
}

func (t TestTransaction) AsMessage(s types.Signer) (types.Message, error) {
	return t.msg, nil
}

func TestNewSourceFilter(t *testing.T){
	addr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	sf, err := NewSourceFilter(1, *(big.NewInt(1)), addr)
	if err != nil {
		t.Errorf("Error creating SourceFilter: %v", err)
	}

	if sf.from.Hex() != addr {
		t.Errorf("Error: SourceFilter from has bad address - expected: %s, actual: %v", addr, sf.from.Hex())
	}
}

func TestSourceFilterQualify(t *testing.T){
	addr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	encoded_addr := common.HexToAddress(addr)
	sf, _ := NewSourceFilter(1, *(big.NewInt(1)), addr)
	tb := TestBlock{testTime: time.Now()}
	tt := TestTransaction{hash: common.HexToHash("0xDEADBEEF"),
			msg: types.NewMessage(encoded_addr, &encoded_addr, 1,
									big.NewInt(1), 1, big.NewInt(1), nil, false)}
	qual, err := sf.QualifyTransaction(tb, tt)
	if err != nil {
		t.Errorf("Error from SourceFilter: %v", err)
	}
	if !qual {
		t.Errorf("Error: SourceFilter failed to qualify transaction")
	}

}

func TestNewDateFilter(t *testing.T){
	_, err := NewDateFilter(time.Now(), time.Now())
	if err != nil {
		t.Errorf("Error creating DateFilter: %v", err)
	}
}

func TestDateFilterQualify(t *testing.T){
	addr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	encoded_addr := common.HexToAddress(addr)
	start := time.Now().AddDate(0,0,-1)
	blockTime := time.Now()
	end := time.Now().AddDate(0,0,1)
	df, _ := NewDateFilter(start, end)
	tb := TestBlock{testTime: blockTime}
	tt := TestTransaction{hash: common.HexToHash("0xDEADBEEF"),
			msg: types.NewMessage(encoded_addr, &encoded_addr, 1,
									big.NewInt(1), 1, big.NewInt(1), nil, false)}
	qual, err := df.QualifyTransaction(tb, tt)
	if err != nil {
		t.Errorf("Error from DateFilter: %v", err)
	}
	if !qual {
		t.Errorf("Error: DateFilter failed to qualify transaction")
	}

}
func TestFilterChainQualifyAnd(t *testing.T){
	addr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	encoded_addr := common.HexToAddress(addr)
	sf, _ := NewSourceFilter(1, *(big.NewInt(1)), addr)
	start := time.Now().AddDate(0,0,-1)
	blockTime := time.Now()
	end := time.Now().AddDate(0,0,1)
	df, _ := NewDateFilter(start, end)
	fc := FilterChain{or: false, filters: []TransactionFilter{sf, df}}
	tb := TestBlock{testTime: blockTime}
	tt := TestTransaction{hash: common.HexToHash("0xDEADBEEF"),
			msg: types.NewMessage(encoded_addr, &encoded_addr, 1,
									big.NewInt(1), 1, big.NewInt(1), nil, false)}
	qual, err := fc.QualifyTransaction(tb, tt)
	if err != nil {
		t.Errorf("Error from FilterChain: %v", err)
	}
	if !qual {
		t.Errorf("Error: FilterChain failed to qualify transaction")
	}

}

func TestFilterChainQualifySingle(t *testing.T){
	addr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	encoded_addr := common.HexToAddress(addr)
	sf, _ := NewSourceFilter(1, *(big.NewInt(1)), addr)
	fc := FilterChain{or: false, filters: []TransactionFilter{sf}}
	tb := TestBlock{testTime: time.Now()}
	tt := TestTransaction{hash: common.HexToHash("0xDEADBEEF"),
			msg: types.NewMessage(encoded_addr, &encoded_addr, 1,
									big.NewInt(1), 1, big.NewInt(1), nil, false)}
	qual, err := fc.QualifyTransaction(tb, tt)
	if err != nil {
		t.Errorf("Error from FilterChain: %v", err)
	}
	if !qual {
		t.Errorf("Error: FilterChain failed to qualify transaction")
	}
}
