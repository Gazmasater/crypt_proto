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



