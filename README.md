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
				return
			}
			books[i] = b
		}

		depthMaxStart, ok := computeMaxStartByDepth(tri, books)
		if !ok || depthMaxStart <= 0 {
			return
		}

		if depthMaxStart < maxStart {
			maxStart = depthMaxStart
		}

		if maxStart < 50 {
			return
		}

		depthFinal, depthDiag, ok := simulateTriangleDepth(maxStart, tri, books)
		if !ok || depthFinal <= 0 {
			return
		}

		depthProfitUSDT := depthFinal - maxStart
		depthProfitPct := depthProfitUSDT / maxStart

		if depthProfitPct <= 0 {
			return
		}

		fmt.Println("00000000000000000000000000000000000000000000000")

		finalAmount = depthFinal
		diag = depthDiag
		profitUSDT = depthProfitUSDT
		profitPct = depthProfitPct
	}



