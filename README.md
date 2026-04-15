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

