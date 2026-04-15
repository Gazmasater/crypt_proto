git rm --cached cmd/arb/metrics/arb_metrics.csv

echo "cmd/arb/metrics/*.csv" >> .gitignore

git add .gitignore
git commit --amend --no-edit


git push origin new_arh --force



git filter-branch --force --index-filter \
'git rm --cached --ignore-unmatch cmd/arb/metrics/arb_metrics.csv' \
--prune-empty --tag-name-filter cat -- new_arh


rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now --aggressive


git push origin new_arh --force






func (c *Calculator) calcTriangle(md *models.MarketData, tri *Triangle) {

	triName := tri.String()

	anchorTS := md.Timestamp
	if anchorTS <= 0 {
		return
	}

	var q [3]queue.Quote
	for i, leg := range tri.Legs {
		quote, ok := c.mem.GetLatestBefore(md.Exchange, leg.Symbol, anchorTS, c.maxQuoteAgeMS)
		if !ok {
			return
		}
		if quote.Bid <= 0 || quote.Ask <= 0 || quote.BidSize <= 0 || quote.AskSize <= 0 {
			return
		}
		q[i] = quote
	}

	ages := [3]int64{
		quoteLagMS(anchorTS, q[0]),
		quoteLagMS(anchorTS, q[1]),
		quoteLagMS(anchorTS, q[2]),
	}
	minAge, maxAge, spreadAge := minMaxSpread(ages)

	maxStart, ok := computeMaxStartTopOfBook(tri, q)
	if !ok || maxStart <= 0 {
		return
	}

	if maxStart < 50 {
		return
	}

	finalAmount, diag, ok := simulateTriangle(maxStart, tri, q)
	if !ok || finalAmount <= 0 {
		return
	}

	profitUSDT := finalAmount - maxStart
	profitPct := profitUSDT / maxStart

	//if profitPct < 0 {
	//	return
	//}

	var books [3]collector.BookSnapshot
	if c.bookSource != nil {
		for i, leg := range tri.Legs {
			b, ok := c.bookSource.GetBookSnapshot(leg.Symbol, 32)
			if !ok {
				log.Printf("[DEPTH REJECT] tri=%s reason=no_snapshot leg=%d symbol=%s",
					triName, i+1, leg.Symbol)
				return
			}
			books[i] = b
		}

		depthMaxStart, ok := computeMaxStartByDepth(tri, books)
		if !ok || depthMaxStart <= 0 {
			log.Printf("[DEPTH REJECT] tri=%s reason=depth_max_start tri=%s depthMaxStart=%.8f",
				triName, triName, depthMaxStart)
			return
		}

		if depthMaxStart < maxStart {
			maxStart = depthMaxStart
		}

		if maxStart < 50 {
			log.Printf("[DEPTH REJECT] tri=%s reason=small_depth_volume maxStart=%.8f",
				triName, maxStart)
			return
		}

		depthFinal, depthDiag, ok := simulateTriangleDepth(maxStart, tri, books)
		if !ok || depthFinal <= 0 {
			log.Printf("[DEPTH REJECT] tri=%s reason=depth_sim_failed maxStart=%.8f depthFinal=%.8f",
				triName, maxStart, depthFinal)
			return
		}

		depthProfitUSDT := depthFinal - maxStart
		depthProfitPct := depthProfitUSDT / maxStart

		if depthProfitPct <= 0 {
			log.Printf("[DEPTH REJECT] tri=%s reason=depth_non_positive maxStart=%.8f final=%.8f profit=%.8f pct=%.6f%%",
				triName, maxStart, depthFinal, depthProfitUSDT, depthProfitPct*100)
			return
		}

		log.Printf("[DEPTH OK] tri=%s maxStart=%.8f final=%.8f profit=%.8f pct=%.6f%%",
			triName, maxStart, depthFinal, depthProfitUSDT, depthProfitPct*100)

		finalAmount = depthFinal
		diag = depthDiag
		profitUSDT = depthProfitUSDT
		profitPct = depthProfitPct
	}

	strength := computeOpportunityStrength(profitPct, maxStart, spreadAge, maxAge)
	triName := fmt.Sprintf("%s->%s->%s", tri.A, tri.B, tri.C)
	maxLegTS := q[0].Timestamp
	if q[1].Timestamp > maxLegTS {
		maxLegTS = q[1].Timestamp
	}
	if q[2].Timestamp > maxLegTS {
		maxLegTS = q[2].Timestamp
	}

	var pendingOp *executor.Opportunity
	if c.oppOut != nil {
		pendingOp = &executor.Opportunity{
			Exchange: md.Exchange,
			Triangle: toExecutorTriangle(tri),
			AnchorTS: anchorTS,
			Quotes: [3]queue.Quote{
				q[0],
				q[1],
				q[2],
			},
			AgesMS: [3]int64{
				ages[0],
				ages[1],
				ages[2],
			},
			MaxStart:   maxStart,
			Books:      books,
			ProfitPct:  profitPct,
			ProfitUSDT: profitUSDT,
			FinalUSDT:  finalAmount,
		}
	}

	written := false

	shouldWrite := c.metrics != nil && c.metrics.ShouldWrite(anchorTS, triName, profitPct, maxStart, maxAge)
	if shouldWrite {
		record := []string{
			strconv.FormatInt(anchorTS, 10),
			triName,
			tri.A, tri.B, tri.C,
			fmtFloat(profitPct), fmtFloat(profitUSDT), fmtFloat(maxStart), fmtFloat(finalAmount),
			fmtFloat(strength), strconv.FormatInt(minAge, 10), strconv.FormatInt(maxAge, 10), strconv.FormatInt(spreadAge, 10),

			tri.Legs[0].Symbol, tri.Legs[0].Side,
			fmtFloat(q[0].Bid), fmtFloat(q[0].Ask), fmtFloat(q[0].BidSize), fmtFloat(q[0].AskSize),
			strconv.FormatInt(ages[0], 10),
			fmtFloat(diag[0].In), fmtFloat(diag[0].Out), fmtFloat(diag[0].TradeQty), fmtFloat(diag[0].TradeNotional), fmtFloat(diag[0].BookLimitIn),

			tri.Legs[1].Symbol, tri.Legs[1].Side,
			fmtFloat(q[1].Bid), fmtFloat(q[1].Ask), fmtFloat(q[1].BidSize), fmtFloat(q[1].AskSize),
			strconv.FormatInt(ages[1], 10),
			fmtFloat(diag[1].In), fmtFloat(diag[1].Out), fmtFloat(diag[1].TradeQty), fmtFloat(diag[1].TradeNotional), fmtFloat(diag[1].BookLimitIn),

			tri.Legs[2].Symbol, tri.Legs[2].Side,
			fmtFloat(q[2].Bid), fmtFloat(q[2].Ask), fmtFloat(q[2].BidSize), fmtFloat(q[2].AskSize),
			strconv.FormatInt(ages[2], 10),
			fmtFloat(diag[2].In), fmtFloat(diag[2].Out), fmtFloat(diag[2].TradeQty), fmtFloat(diag[2].TradeNotional), fmtFloat(diag[2].BookLimitIn),
		}

		if err := c.metrics.Write(record); err != nil {
			log.Printf("[Calculator] metrics write error: %v", err)
		} else {
			written = true
		}
	}

	if c.summary != nil {
		c.summary.Observe(triName, profitPct, profitUSDT, written)
	}

	//if profitPct > 0.0 && maxStart > 50 {
	msg := fmt.Sprintf(
		"[ARB] %s→%s→%s | %.4f%% | volume=%.2f USDT | profit=%.6f USDT | "+
			"anchor=%d | l1=%s %s out=%.8f age=%dms | "+
			"l2=%s %s out=%.8f age=%dms | "+
			"l3=%s %s out=%.8f age=%dms",
		tri.A, tri.B, tri.C,
		profitPct*100, maxStart, profitUSDT,
		anchorTS,
		tri.Legs[0].Symbol, tri.Legs[0].Side, diag[0].Out, ages[0],
		tri.Legs[1].Symbol, tri.Legs[1].Side, diag[1].Out, ages[1],
		tri.Legs[2].Symbol, tri.Legs[2].Side, diag[2].Out, ages[2],
	)
	log.Println(msg)
	c.fileLog.Println(msg)
	//}

	if pendingOp != nil {
		if !c.shouldEmitOpportunity(triName, maxLegTS, maxStart) {
			return
		}

		select {
		case c.oppOut <- pendingOp:
		default:
		}
	}
}




  [{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/calculator/arb.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "tri.String undefined (type *Triangle has no field or method String)",
	"source": "compiler",
	"startLineNumber": 346,
	"startColumn": 17,
	"endLineNumber": 346,
	"endColumn": 23,
	"modelVersionId": 5,
	"origin": "extHost1"
}]


[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/calculator/arb.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "NoNewVar",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "NoNewVar"
		}
	},
	"severity": 8,
	"message": "no new variables on left side of :=",
	"source": "compiler",
	"startLineNumber": 448,
	"startColumn": 2,
	"endLineNumber": 448,
	"endColumn": 59,
	"modelVersionId": 5,
	"origin": "extHost1"
}]




