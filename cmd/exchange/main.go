package main

import (
	"log"

	"crypt_proto/cmd/exchange/builder"
	"crypt_proto/cmd/exchange/common"
	"crypt_proto/cmd/exchange/kucoin"
)

func main() {

	// ---------- KUCOIN ----------
	kucoinMarkets := kucoin.LoadMarkets()
	kucoinTriangles := builder.BuildTriangles(kucoinMarkets, "USDT")
	if err := common.SaveTrianglesCSV(
		"data/kucoin_triangles_usdt.csv",
		kucoinTriangles,
	); err != nil {
		log.Fatalf("kucoin csv error: %v", err)
	}

}
