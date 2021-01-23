package internal

import (
	"time"
)

type Configuration struct {
	BlocksSince		time.Time
	FromAccounts	[]string
	ToAccounts		[]string
}
