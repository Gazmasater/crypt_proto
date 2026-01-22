package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

/* ================= CONFIG ================= */

const (
	baseURL = "https://api.kucoin.com"

	sym1 = "WAVES-USDT"
	sym2 = "WAVES-BTC"
	sym3 = "BTC-USDT"
)

/* ================= SYMBOL INFO ================= */

type SymbolInfo struct {
	Symbol       string  `json:"symbol"`
	BaseMinSize  float64 `json:"baseMinSize,string"`
	BaseMaxSize  float64 `json:"baseMaxSize,string"`
	QuoteMinSize float64 `json:"quoteMinSize,string"`
}

func getStep(symbol string) (float64, error) {
	resp, err := http.Get(baseURL + "/api/v1/symbols/" + symbol)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var r struct {
		Code string     `json:"code"`
		Data SymbolInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, err
	}

	if r.Code != "200000" {
		return 0, fmt.Errorf("failed to get symbol info for %s", symbol)
	}

	return r.Data.BaseMinSize, nil
}

/* ================= MAIN ================= */

func main() {
	log.Println("Getting steps for triangle: USDT → DASH → BTC → USDT")

	step1, err := getStep(sym1)
	if err != nil {
		log.Fatal(err)
	}
	step2, err := getStep(sym2)
	if err != nil {
		log.Fatal(err)
	}
	step3, err := getStep(sym3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("====== TRADING STEPS ======")
	fmt.Printf("%s step: %.8f\n", sym1, step1)
	fmt.Printf("%s step: %.8f\n", sym2, step2)
	fmt.Printf("%s step: %.8f\n", sym3, step3)
}
