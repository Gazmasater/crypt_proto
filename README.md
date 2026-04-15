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







func computeMaxStartByDepth(tri *Triangle, quotes [3]queue.Quote, books [3]collector.BookSnapshot) (float64, bool) {
	high, ok := computeMaxStartTopOfBook(tri, quotes)
	if !ok || high <= 0 || !isFinite(high) {
		return 0, false
	}

	low := 0.0

	for i := 0; i < 32; i++ {
		mid := (low + high) / 2
		if mid <= 0 || !isFinite(mid) {
			break
		}

		if _, _, ok := simulateTriangleDepth(mid, tri, books); ok {
			low = mid
		} else {
			high = mid
		}
	}

	if low <= 0 || !isFinite(low) {
		return 0, false
	}

	return low, true
}




depthMaxStart, ok := computeMaxStartByDepth(tri, q, books)
