package main

import (
	"log"

	"crypt_proto/cmd/exchange/builder"
	"crypt_proto/cmd/exchange/common"
	"crypt_proto/cmd/exchange/kucoin"
	"crypt_proto/cmd/exchange/mexc"
	"crypt_proto/cmd/exchange/okx"
)

func main() {

	// ---------- KUCOIN ----------
	kucoinMarkets := kucoin.LoadMarkets()
	kucoinTriangles := builder.BuildTriangles(kucoinMarkets, "USDC")
	if err := common.SaveTrianglesCSV(
		"data/kucoin_triangles_usdt.csv",
		kucoinTriangles,
	); err != nil {
		log.Fatalf("kucoin csv error: %v", err)
	}

	// ---------- MEXC ----------
	mexcMarkets := mexc.LoadMarkets()
	mexcTriangles := builder.BuildTriangles(mexcMarkets, "USDT")
	if err := common.SaveTrianglesCSV(
		"data/mexc_triangles_usdt.csv",
		mexcTriangles,
	); err != nil {
		log.Fatalf("mexc csv error: %v", err)
	}

	// ---------- OKX ----------
	okxMarkets := okx.LoadMarkets()
	okxTriangles := builder.BuildTriangles(okxMarkets, "USDT")
	if err := common.SaveTrianglesCSV(
		"data/okx_triangles_usdt.csv",
		okxTriangles,
	); err != nil {
		log.Fatalf("okx csv error: %v", err)
	}
}
