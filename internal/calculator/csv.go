package calculator

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func ParseTrianglesFromCSV(path string) ([]*Triangle, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	headerIndex := make(map[string]int, len(rows[0]))
	for i, col := range rows[0] {
		headerIndex[strings.TrimSpace(col)] = i
	}

	res := make([]*Triangle, 0, len(rows)-1)
	for rowIdx, row := range rows[1:] {
		if len(row) < 15 {
			continue
		}

		t := &Triangle{
			A: strings.TrimSpace(row[0]),
			B: strings.TrimSpace(row[1]),
			C: strings.TrimSpace(row[2]),
		}

		for i, col := range triangleLegColumns {
			leg, err := parseTriangleLeg(row[col])
			if err != nil {
				return nil, fmt.Errorf("row %d leg %d: %w", rowIdx+2, i+1, err)
			}
			t.Legs[i] = leg
			t.Rules[i] = parseLegRules(row, headerIndex, i+1, leg)
		}

		res = append(res, t)
	}
	return res, nil
}

func parseTriangleLeg(raw string) (LegIndex, error) {
	leg := strings.ToUpper(strings.TrimSpace(raw))
	parts := strings.Fields(leg)
	if len(parts) != 2 {
		return LegIndex{}, fmt.Errorf("bad leg format: %q", raw)
	}

	isBuy := parts[0] == "BUY"
	if parts[0] != "BUY" && parts[0] != "SELL" {
		return LegIndex{}, fmt.Errorf("bad leg side: %q", raw)
	}

	pair := strings.Split(parts[1], "/")
	if len(pair) != 2 {
		return LegIndex{}, fmt.Errorf("bad pair format: %q", raw)
	}

	symbol := pair[0] + "-" + pair[1]
	return LegIndex{Key: "KuCoin|" + symbol, Symbol: symbol, IsBuy: isBuy}, nil
}

func parseLegRules(row []string, headerIndex map[string]int, legNum int, leg LegIndex) LegRules {
	rules := LegRules{
		Symbol: leg.Symbol,
		Side:   map[bool]string{true: "BUY", false: "SELL"}[leg.IsBuy],
		Fee:    defaultTakerFee,
	}

	prefix := fmt.Sprintf("Leg%d", legNum)
	rules.Symbol = csvString(row, headerIndex, prefix+"Symbol", rules.Symbol)
	rules.Side = csvString(row, headerIndex, prefix+"Side", rules.Side)
	rules.Base = csvString(row, headerIndex, prefix+"Base", "")
	rules.Quote = csvString(row, headerIndex, prefix+"Quote", "")
	rules.QtyStep = csvFloat(row, headerIndex, prefix+"QtyStep", csvFloat(row, headerIndex, fmt.Sprintf("Step%d", legNum), 0))
	rules.QuoteStep = csvFloat(row, headerIndex, prefix+"QuoteStep", 0)
	rules.PriceStep = csvFloat(row, headerIndex, prefix+"PriceStep", 0)
	rules.MinQty = csvFloat(row, headerIndex, prefix+"MinQty", csvFloat(row, headerIndex, fmt.Sprintf("MinQty%d", legNum), 0))
	rules.MinQuote = csvFloat(row, headerIndex, prefix+"MinQuote", 0)
	rules.MinNotional = csvFloat(row, headerIndex, prefix+"MinNotional", csvFloat(row, headerIndex, fmt.Sprintf("MinNotional%d", legNum), 0))
	rules.Fee = csvFloat(row, headerIndex, prefix+"Fee", rules.Fee)
	return rules
}

func csvString(row []string, headerIndex map[string]int, key, fallback string) string {
	idx, ok := headerIndex[key]
	if !ok || idx >= len(row) {
		return fallback
	}
	value := strings.TrimSpace(row[idx])
	if value == "" {
		return fallback
	}
	return value
}

func csvFloat(row []string, headerIndex map[string]int, key string, fallback float64) float64 {
	idx, ok := headerIndex[key]
	if !ok || idx >= len(row) {
		return fallback
	}
	value := strings.TrimSpace(row[idx])
	if value == "" {
		return fallback
	}
	var parsed float64
	_, err := fmt.Sscanf(value, "%f", &parsed)
	if err != nil {
		return fallback
	}
	return parsed
}
