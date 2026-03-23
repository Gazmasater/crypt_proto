package main

import (
	"log"
	"os"

	"crypt_proto/cmd/exchange/builder"
	"crypt_proto/cmd/exchange/common"
	"crypt_proto/cmd/exchange/kucoin"
	"crypt_proto/cmd/exchange/mexc"
	"crypt_proto/cmd/exchange/okx"
)

type exchangeJob struct {
	name     string
	loadFunc func() map[string]common.Market
	outFile  string
}

func main() {
	if err := os.MkdirAll("data", 0o755); err != nil {
		log.Fatalf("create data dir error: %v", err)
	}

	jobs := []exchangeJob{
		{
			name:     "kucoin",
			loadFunc: kucoin.LoadMarkets,
			outFile:  "data/kucoin_triangles_usdt.csv",
		},
		{
			name:     "mexc",
			loadFunc: mexc.LoadMarkets,
			outFile:  "data/mexc_triangles_usdt.csv",
		},
		{
			name:     "okx",
			loadFunc: okx.LoadMarkets,
			outFile:  "data/okx_triangles_usdt.csv",
		},
	}

	for _, job := range jobs {
		markets := job.loadFunc()
		log.Printf("[%s] markets loaded: %d", job.name, len(markets))

		triangles := builder.BuildTriangles(markets, "USDT")
		log.Printf("[%s] triangles built: %d", job.name, len(triangles))

		if err := common.SaveTrianglesCSV(job.outFile, triangles); err != nil {
			log.Fatalf("[%s] save csv error: %v", job.name, err)
		}

		log.Printf("[%s] csv saved: %s", job.name, job.outFile)
	}
}
