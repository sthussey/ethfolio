package main

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/ethereum/go-ethereum/common"
	ulid "github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

type Lot struct {
	LotId		ulid.ULID
	Date		time.Time
	Rate		decimal.Decimal
	TxnList		[]LotTxn
}

type LotTxn struct {
	TxnFrom		common.Address
	TxnNonce	uint64
	EthAmount	decimal.Decimal
}

func NewLot(date time.Time) Lot {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	fmt.Println()
	return Lot{ LotId: ulid.MustNew(ulid.Timestamp(t), entropy),
				Date: date}
}

func NewLotTxn(from common.Address, nonce uint64, amt decimal.Decimal) LotTxn {
	return LotTxn { TxnFrom: from,
					TxnNonce: nonce,
					EthAmount: amt}
}

func (l *Lot) AddTxn(from common.Address, nonce uint64, amt decimal.Decimal) error {
	txn := NewLotTxn(from, nonce, amt);
	l.TxnList = append(l.TxnList, txn)
	return nil
}
