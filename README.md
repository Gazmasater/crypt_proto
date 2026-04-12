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



2026/04/12 05:14:12 pprof on http://localhost:6060/debug/pprof/
2026/04/12 05:14:12 [KuCoin] started with 2 WS
2026/04/12 05:14:12 [Main] KuCoinCollector started
2026/04/12 05:14:12 [Calculator] indexed 206 symbols
2026/04/12 05:14:13 [KuCoin WS 0] connected
2026/04/12 05:14:13 [KuCoin WS 1] connected
2026/04/12 05:14:23 [KuCoin WS 1] bootstrap complete symbol=ONE-USDT seq=2149409236 bid=0.0020600000 ask=0.0020650000
2026/04/12 05:14:23 [KuCoin WS 1] bootstrap complete symbol=ONT-BTC seq=599150049 bid=0.0000011810 ask=0.0000012000
2026/04/12 05:14:24 [KuCoin WS 1] bootstrap complete symbol=ONT-ETH seq=618904410 bid=0.0000382800 ask=0.0000398400
2026/04/12 05:14:24 [KuCoin WS 1] bootstrap complete symbol=ONT-USDT seq=1628683920 bid=0.0853400000 ask=0.0855900000
2026/04/12 05:14:24 [KuCoin WS 1] bootstrap complete symbol=PAXG-BTC seq=361658049 bid=0.0656450000 ask=0.0657330000
2026/04/12 05:14:24 [KuCoin WS 1] bootstrap complete symbol=PAXG-USDT seq=1638364015 bid=4714.4300000000 ask=4714.4400000000
2026/04/12 05:14:25 [KuCoin WS 1] bootstrap complete symbol=PEPE-KCS seq=146570611 bid=0.0000004166 ask=0.0000004198
2026/04/12 05:14:25 [KuCoin WS 1] bootstrap complete symbol=PEPE-USDT seq=5047832142 bid=0.0000035110 ask=0.0000035120
2026/04/12 05:14:25 [KuCoin WS 1] bootstrap complete symbol=POND-BTC seq=123767028 bid=0.0000000302 ask=0.0000000312
2026/04/12 05:14:26 [KuCoin WS 1] bootstrap complete symbol=POND-USDT seq=434215349 bid=0.0021990000 ask=0.0022090000
2026/04/12 05:14:26 [KuCoin WS 1] bootstrap complete symbol=RLC-BTC seq=144009186 bid=0.0000057800 ask=0.0000058000
2026/04/12 05:14:26 [KuCoin WS 1] bootstrap complete symbol=RLC-USDT seq=1052807489 bid=0.4149000000 ask=0.4158000000
2026/04/12 05:14:26 [KuCoin WS 1] bootstrap complete symbol=RSR-BTC seq=132290553 bid=0.0000000204 ask=0.0000000207
2026/04/12 05:14:27 [KuCoin WS 1] bootstrap complete symbol=RSR-USDT seq=1855588478 bid=0.0014760000 ask=0.0014770000
2026/04/12 05:14:27 [KuCoin WS 1] bootstrap complete symbol=RUNE-BTC seq=1114737078 bid=0.0000054300 ask=0.0000054700
2026/04/12 05:14:27 [KuCoin WS 1] bootstrap complete symbol=RUNE-USDT seq=3496719833 bid=0.3911000000 ask=0.3914000000
2026/04/12 05:14:27 [KuCoin WS 1] bootstrap complete symbol=SCRT-BTC seq=69764974 bid=0.0000012740 ask=0.0000012770
2026/04/12 05:14:28 [KuCoin WS 1] bootstrap complete symbol=SCRT-USDT seq=341446510 bid=0.0914000000 ask=0.0916000000
2026/04/12 05:14:28 [KuCoin WS 1] bootstrap complete symbol=SHIB-DOGE seq=676578400 bid=0.0000635300 ask=0.0000637800
2026/04/12 05:14:28 [KuCoin WS 1] bootstrap complete symbol=SHIB-USDT seq=5998377022 bid=0.0000058160 ask=0.0000058170
2026/04/12 05:14:28 [KuCoin WS 0] bootstrap complete symbol=A-BTC seq=168712050 bid=0.0000010850 ask=0.0000010990
2026/04/12 05:14:29 [KuCoin WS 1] bootstrap complete symbol=SNX-BTC seq=424219495 bid=0.0000039400 ask=0.0000040300
2026/04/12 05:14:29 [KuCoin WS 0] bootstrap complete symbol=A-ETH seq=583807907 bid=0.0000352000 ask=0.0000414000
2026/04/12 05:14:29 [KuCoin WS 1] bootstrap complete symbol=SNX-ETH seq=505468717 bid=0.0001270000 ask=0.0001330000
2026/04/12 05:14:29 [KuCoin WS 0] bootstrap complete symbol=A-USDT seq=143866488 bid=0.0782000000 ask=0.0785000000
2026/04/12 05:14:29 [KuCoin WS 1] bootstrap complete symbol=SNX-USDT seq=1536358804 bid=0.2857000000 ask=0.2873000000
2026/04/12 05:14:29 [KuCoin WS 0] bootstrap complete symbol=AAVE-BTC seq=1081936286 bid=0.0012580000 ask=0.0012600000
2026/04/12 05:14:29 [KuCoin WS 1] bootstrap complete symbol=SOL-KCS seq=458209243 bid=9.8280000000 ask=9.8490000000
2026/04/12 05:14:30 [KuCoin WS 0] bootstrap complete symbol=AAVE-USDT seq=9001767701 bid=90.3560000000 ask=90.3570000000
2026/04/12 05:14:30 [KuCoin WS 1] bootstrap complete symbol=SOL-USDT seq=19082261690 bid=82.5600000000 ask=82.5700000000
2026/04/12 05:14:30 [KuCoin WS 0] bootstrap complete symbol=ADA-BTC seq=1105181143 bid=0.0000034000 ask=0.0000034100
2026/04/12 05:14:30 [KuCoin WS 1] bootstrap complete symbol=STORJ-ETH seq=100779664 bid=0.0000437000 ask=0.0000442000
2026/04/12 05:14:30 [KuCoin WS 0] bootstrap complete symbol=ADA-KCS seq=527411578 bid=0.0291000000 ask=0.0291800000
2026/04/12 05:14:30 [KuCoin WS 1] bootstrap complete symbol=STORJ-USDT seq=901041830 bid=0.0979000000 ask=0.0982000000
2026/04/12 05:14:30 [KuCoin WS 0] bootstrap complete symbol=ADA-USDT seq=8277282446 bid=0.2445000000 ask=0.2446000000
2026/04/12 05:14:30 [KuCoin WS 1] bootstrap complete symbol=STX-BTC seq=442983227 bid=0.0000029700 ask=0.0000030000

