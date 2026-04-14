package calculator

import (
	"math"
	"strings"

	"crypt_proto/internal/collector"
)

func simulateTriangleDepth(startUSDT float64, tri *Triangle, books [3]collector.BookSnapshot) (float64, [3]legExecution, bool) {
	var diag [3]legExecution
	amount := startUSDT
	for i := 0; i < 3; i++ {
		out, d, ok := executeLegByDepth(amount, tri.Legs[i], books[i])
		if !ok {
			return 0, diag, false
		}
		diag[i] = d
		amount = out
	}
	return amount, diag, true
}

func executeLegByDepth(in float64, leg LegRule, book collector.BookSnapshot) (float64, legExecution, bool) {
	side := strings.ToUpper(strings.TrimSpace(leg.Side))
	if side == "" {
		side = detectSideFromRawLeg(leg.RawLeg)
	}
	if side != "BUY" && side != "SELL" {
		return 0, legExecution{}, false
	}
	if in <= 0 {
		return 0, legExecution{}, false
	}

	qtyStep := firstPositive(leg.QtyStep, leg.Step)
	minQty := firstPositive(leg.LegMinQty, leg.MinQty)
	minQuote := leg.LegMinQuote
	minNotional := firstPositive(leg.LegMinNotnl, leg.MinNotional)

	switch side {
	case "BUY":
		if len(book.Asks) == 0 {
			return 0, legExecution{}, false
		}
		remainingQuote := in
		totalQty := 0.0
		totalNotional := 0.0
		bookLimitIn := 0.0
		for _, lvl := range book.Asks {
			if lvl.Price <= 0 || lvl.Size <= 0 {
				continue
			}
			bookLimitIn += lvl.Price * lvl.Size
			if remainingQuote <= eps() {
				break
			}
			maxSpend := lvl.Price * lvl.Size
			spend := math.Min(remainingQuote, maxSpend)
			qty := spend / lvl.Price
			totalQty += qty
			totalNotional += qty * lvl.Price
			remainingQuote -= qty * lvl.Price
		}
		tradeQty := floorToStep(totalQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}
		avgPrice := totalNotional / totalQty
		tradeNotional := tradeQty * avgPrice
		if tradeNotional <= 0 || tradeNotional > in+1e-9 {
			return 0, legExecution{}, false
		}
		if minQty > 0 && tradeQty+eps() < minQty {
			return 0, legExecution{}, false
		}
		if minQuote > 0 && tradeNotional+eps() < minQuote {
			return 0, legExecution{}, false
		}
		if minNotional > 0 && tradeNotional+eps() < minNotional {
			return 0, legExecution{}, false
		}
		outBase := tradeQty * feeM
		if outBase <= 0 {
			return 0, legExecution{}, false
		}
		return outBase, legExecution{In: in, Out: outBase, Price: avgPrice, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	case "SELL":
		if len(book.Bids) == 0 {
			return 0, legExecution{}, false
		}
		remainingBase := in
		totalQty := 0.0
		totalNotional := 0.0
		bookLimitIn := 0.0
		for _, lvl := range book.Bids {
			if lvl.Price <= 0 || lvl.Size <= 0 {
				continue
			}
			bookLimitIn += lvl.Size
			if remainingBase <= eps() {
				break
			}
			qty := math.Min(remainingBase, lvl.Size)
			totalQty += qty
			totalNotional += qty * lvl.Price
			remainingBase -= qty
		}
		tradeQty := floorToStep(totalQty, qtyStep)
		if tradeQty <= 0 {
			return 0, legExecution{}, false
		}
		avgPrice := totalNotional / totalQty
		tradeNotional := tradeQty * avgPrice
		if tradeNotional <= 0 {
			return 0, legExecution{}, false
		}
		if minQty > 0 && tradeQty+eps() < minQty {
			return 0, legExecution{}, false
		}
		if minQuote > 0 && tradeNotional+eps() < minQuote {
			return 0, legExecution{}, false
		}
		if minNotional > 0 && tradeNotional+eps() < minNotional {
			return 0, legExecution{}, false
		}
		outQuote := tradeNotional * feeM
		if outQuote <= 0 {
			return 0, legExecution{}, false
		}
		return outQuote, legExecution{In: in, Out: outQuote, Price: avgPrice, BookLimitIn: bookLimitIn, TradeQty: tradeQty, TradeNotional: tradeNotional}, true
	}
	return 0, legExecution{}, false
}

func computeMaxStartByDepth(tri *Triangle, books [3]collector.BookSnapshot) (float64, bool) {
	low, high := 0.0, math.MaxFloat64
	for i := 0; i < 3; i++ {
		leg := tri.Legs[i]
		side := strings.ToUpper(strings.TrimSpace(leg.Side))
		if side == "" {
			side = detectSideFromRawLeg(leg.RawLeg)
		}
		var bookCap float64
		switch side {
		case "BUY":
			for _, lvl := range books[i].Asks {
				bookCap += lvl.Price * lvl.Size
			}
		case "SELL":
			for _, lvl := range books[i].Bids {
				bookCap += lvl.Size
			}
		default:
			return 0, false
		}
		if bookCap <= 0 {
			return 0, false
		}
		if bookCap < high {
			high = bookCap
		}
	}
	if !isFinite(high) || high <= 0 {
		return 0, false
	}
	for i := 0; i < 20; i++ {
		mid := (low + high) / 2
		if mid <= 0 {
			break
		}
		if _, _, ok := simulateTriangleDepth(mid, tri, books); ok {
			low = mid
		} else {
			high = mid
		}
	}
	if low <= 0 {
		return 0, false
	}
	return low, true
}
