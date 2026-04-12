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


