package main

import (
	"os"
	"time"
	"github.com/sthussey/ethfolio/internal"
)

func main() {
	cfg := internal.Configuration{BlocksSince: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), ToAccounts: []string{"0xa36e426b6754cdc5d6574c7b10a539e89d011e4f"}}

	internal.ProcessBlocks(cfg)

	os.Exit(0)
}
