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



[{
	"resource": "/home/gaz358/myprog/crypt_proto/internal/calculator/arb.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UndeclaredName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredName"
		}
	},
	"severity": 8,
	"message": "undefined: triName",
	"source": "compiler",
	"startLineNumber": 397,
	"startColumn": 6,
	"endLineNumber": 397,
	"endColumn": 13,
	"modelVersionId": 3,
	"origin": "extHost1"
}]

