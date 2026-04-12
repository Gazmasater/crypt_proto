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


Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.358.pb.gz
File: arb
Build ID: 0ba99bd9ae2047c4bad9c3d8309b6d2d1b541df0
Type: cpu
Time: 2026-04-12 16:56:50 MSK
Duration: 30s, Total samples = 3.32s (11.07%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1900ms, 57.23% of 3320ms total
Dropped 76 nodes (cum <= 16.60ms)
Showing top 10 nodes out of 149
      flat  flat%   sum%        cum   cum%
     820ms 24.70% 24.70%      820ms 24.70%  internal/runtime/syscall.Syscall6
     310ms  9.34% 34.04%      430ms 12.95%  internal/runtime/maps.(*Iter).Next
     170ms  5.12% 39.16%      450ms 13.55%  github.com/tidwall/gjson.parseObject
     170ms  5.12% 44.28%      170ms  5.12%  runtime.futex
     110ms  3.31% 47.59%      110ms  3.31%  github.com/tidwall/gjson.parseSquash
      90ms  2.71% 50.30%      570ms 17.17%  crypt_proto/internal/collector.computeTop
      80ms  2.41% 52.71%       80ms  2.41%  runtime.duffcopy
      50ms  1.51% 54.22%       50ms  1.51%  github.com/tidwall/gjson.parseObjectPath
      50ms  1.51% 55.72%       50ms  1.51%  github.com/tidwall/gjson.parseString
      50ms  1.51% 57.23%       50ms  1.51%  internal/runtime/maps.ctrlGroup.matchFull (inline)
(pprof) 





Saved profile in /home/gaz358/pprof/pprof.arb.samples.cpu.359.pb.gz
File: arb
Build ID: ec679a626462caba32cf0b7059a06bc6a77ccf33
Type: cpu
Time: 2026-04-12 17:15:14 MSK
Duration: 30s, Total samples = 3.13s (10.43%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1.76s, 56.23% of 3.13s total
Dropped 93 nodes (cum <= 0.02s)
Showing top 10 nodes out of 145
      flat  flat%   sum%        cum   cum%
     1.04s 33.23% 33.23%      1.04s 33.23%  internal/runtime/syscall.Syscall6
     0.22s  7.03% 40.26%      0.22s  7.03%  runtime.futex
     0.12s  3.83% 44.09%      0.31s  9.90%  github.com/tidwall/gjson.parseObject
     0.07s  2.24% 46.33%      0.07s  2.24%  github.com/tidwall/gjson.parseSquash
     0.07s  2.24% 48.56%      0.07s  2.24%  runtime.duffcopy
     0.07s  2.24% 50.80%      0.07s  2.24%  runtime.nextFreeFast
     0.06s  1.92% 52.72%      0.06s  1.92%  strconv.readFloat
     0.05s  1.60% 54.31%      0.05s  1.60%  runtime.casgstatus
     0.03s  0.96% 55.27%      0.03s  0.96%  crypto/tls.(*halfConn).explicitNonceLen
     0.03s  0.96% 56.23%      0.03s  0.96%  github.com/tidwall/gjson.Result.String
(pprof) 






2026/04/12 18:07:51 [Calculator] summary checked=286 written=0 profitable=0 best_pct=-0.2310% best_usdt=-0.137316 best_tri=USDT->ETH->KCS
2026/04/12 18:07:57 [Executor] summary accepted=0 reject_duplicate=14 reject_stale=0 reject_desync=378 reject_small_safe=0 reject_non_positive=609 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:08:01 [Calculator] summary checked=1305 written=0 profitable=0 best_pct=-0.2264% best_usdt=-0.341329 best_tri=USDT->ETH->BTC
2026/04/12 18:08:07 [Executor] summary accepted=0 reject_duplicate=31 reject_stale=0 reject_desync=700 reject_small_safe=0 reject_non_positive=999 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:08:11 [Calculator] summary checked=2155 written=0 profitable=0 best_pct=-0.2264% best_usdt=-0.341329 best_tri=USDT->ETH->BTC
2026/04/12 18:08:17 [Executor] summary accepted=0 reject_duplicate=66 reject_stale=0 reject_desync=1366 reject_small_safe=0 reject_non_positive=1517 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:08:21 [Calculator] summary checked=3500 written=0 profitable=0 best_pct=-0.2264% best_usdt=-0.341329 best_tri=USDT->ETH->BTC
2026/04/12 18:08:27 [Executor] summary accepted=0 reject_duplicate=128 reject_stale=0 reject_desync=2340 reject_small_safe=0 reject_non_positive=2706 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:08:31 [Calculator] summary checked=5686 written=0 profitable=0 best_pct=-0.2245% best_usdt=-0.234158 best_tri=USDT->BTC->ETH
2026/04/12 18:08:37 [Executor] summary accepted=0 reject_duplicate=182 reject_stale=0 reject_desync=3454 reject_small_safe=0 reject_non_positive=3452 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:08:41 [Calculator] summary checked=7437 written=0 profitable=0 best_pct=-0.1742% best_usdt=-0.100962 best_tri=USDT->KCS->ETH
2026/04/12 18:08:47 [Executor] summary accepted=0 reject_duplicate=250 reject_stale=0 reject_desync=4284 reject_small_safe=0 reject_non_positive=4321 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:08:51 [Calculator] summary checked=9560 written=0 profitable=0 best_pct=-0.1742% best_usdt=-0.100962 best_tri=USDT->KCS->ETH
2026/04/12 18:08:57 [Executor] summary accepted=0 reject_duplicate=294 reject_stale=0 reject_desync=5507 reject_small_safe=0 reject_non_positive=5284 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:09:01 [Calculator] summary checked=11235 written=0 profitable=0 best_pct=-0.1742% best_usdt=-0.100962 best_tri=USDT->KCS->ETH
2026/04/12 18:09:08 [Executor] summary accepted=0 reject_duplicate=339 reject_stale=0 reject_desync=6604 reject_small_safe=0 reject_non_positive=6012 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:09:11 [Calculator] summary checked=13349 written=0 profitable=0 best_pct=-0.1742% best_usdt=-0.100962 best_tri=USDT->KCS->ETH
2026/04/12 18:09:18 [Executor] summary accepted=0 reject_duplicate=371 reject_stale=0 reject_desync=7088 reject_small_safe=0 reject_non_positive=6322 reject_too_small_profit=0 reject_sim_failed=0
2026/04/12 18:09:21 [Calculator] summary checked=14182 written=0 profitable=0 best_pct=-0.1742% best_usdt=-0.100962 best_tri=USDT->KCS->ETH
