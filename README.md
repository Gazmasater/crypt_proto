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





if maxStart < 20 {
    return
}




