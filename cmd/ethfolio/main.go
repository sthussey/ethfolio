package main

import (
	"os"
	"time"
	"github.com/sthussey/ethfolio/internal"
)

func main() {
	cfg := internal.Configuration{BlocksSince: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), ToAccounts: []string{"0x6677429c9Fd93F15BC1679D916B83807BC4df5e2"}}

	metrics := internal.InitializeMetrics()
	internal.ProcessBlocks(cfg, metrics)

	os.Exit(0)
}
