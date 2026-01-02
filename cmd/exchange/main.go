package main

import (
	"crypt_proto/cmd/exchange/builder"
	"crypt_proto/cmd/exchange/common"
	"crypt_proto/cmd/exchange/kucoin"
)

func main() {
	markets := kucoin.LoadMarkets()

	triangles := builder.BuildTriangles(markets, "USDT")

	common.SaveTrianglesCSV("data/kucoin_triangles.csv", triangles)
}
