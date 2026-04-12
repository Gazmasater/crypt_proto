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





2026/04/12 09:35:10 [KuCoin] started with 2 WS
2026/04/12 09:35:10 [Main] KuCoinCollector started
2026/04/12 09:35:10 [Calculator] indexed 206 symbols
2026/04/12 09:35:12 [KuCoin WS 1] connected
2026/04/12 09:35:12 [KuCoin WS 0] connected
2026/04/12 09:35:22 [KuCoin WS 1] bootstrap complete symbol=PAXG-USDT seq=1638393549 bid=4709.2200000000 ask=4709.2300000000
2026/04/12 09:35:22 [KuCoin WS 1] bootstrap complete symbol=ONE-USDT seq=2149429672 bid=0.0020580000 ask=0.0020640000
2026/04/12 09:35:23 [KuCoin WS 1] bootstrap complete symbol=PEPE-USDT seq=5048158396 bid=0.0000035160 ask=0.0000035170
2026/04/12 09:35:23 [KuCoin WS 1] bootstrap complete symbol=PAXG-BTC seq=361703116 bid=0.0656690000 ask=0.0657520000
2026/04/12 09:35:23 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/12 09:35:23 [KuCoin WS 1] bootstrap complete symbol=RLC-BTC seq=144021637 bid=0.0000058000 ask=0.0000058200
2026/04/12 09:35:23 [KuCoin WS 1] bootstrap complete symbol=RLC-USDT seq=1052838003 bid=0.4156000000 ask=0.4164000000
2026/04/12 09:35:23 [KuCoin WS 1] bootstrap complete symbol=RSR-BTC seq=132295982 bid=0.0000000205 ask=0.0000000209
2026/04/12 09:35:23 [KuCoin WS 1] bootstrap complete symbol=RSR-USDT seq=1855614278 bid=0.0014830000 ask=0.0014840000
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=POND-USDT seq=434250017 bid=0.0021970000 ask=0.0022080000
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=RUNE-BTC seq=1114757557 bid=0.0000054300 ask=0.0000054600
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap progress 10/80 last=RUNE-BTC
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=PEPE-KCS seq=146639064 bid=0.0000004178 ask=0.0000004202
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=SCRT-BTC seq=69769660 bid=0.0000012770 ask=0.0000012810
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=RUNE-USDT seq=3496765179 bid=0.3904000000 ask=0.3907000000
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=POND-BTC seq=123773023 bid=0.0000000303 ask=0.0000000331
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=SHIB-DOGE seq=676605946 bid=0.0000636400 ask=0.0000637900
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=ONT-ETH seq=618951601 bid=0.0000389000 ask=0.0000395700
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=SHIB-USDT seq=5998509874 bid=0.0000058280 ask=0.0000058290
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=SNX-ETH seq=505533375 bid=0.0001270000 ask=0.0001340000
2026/04/12 09:35:24 [KuCoin WS 1] bootstrap complete symbol=SOL-KCS seq=458287133 bid=9.7970000000 ask=9.8340000000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=STORJ-ETH seq=100783604 bid=0.0000442000 ask=0.0000445000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap progress 20/80 last=STORJ-ETH
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=SNX-USDT seq=1536436169 bid=0.2843000000 ask=0.2857000000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=SOL-USDT seq=19083322340 bid=82.3700000000 ask=82.3800000000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=STORJ-USDT seq=901051552 bid=0.0981000000 ask=0.0984000000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=STX-BTC seq=442995752 bid=0.0000029700 ask=0.0000029900
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=STX-USDT seq=3150667376 bid=0.2140000000 ask=0.2141000000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=ONT-USDT seq=1628846993 bid=0.0863700000 ask=0.0865300000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=SCRT-USDT seq=341462763 bid=0.0916000000 ask=0.0918000000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=SUI-KCS seq=243396033 bid=0.1085000000 ask=0.1092000000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=ONT-BTC seq=599183383 bid=0.0000012040 ask=0.0000012270
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=SNX-BTC seq=424273805 bid=0.0000039200 ask=0.0000041200
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap progress 30/80 last=SNX-BTC
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=TEL-ETH seq=1361299838 bid=0.0000009473 ask=0.0000009717
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=TEL-USDT seq=5957451467 bid=0.0021030000 ask=0.0021060000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=TRAC-ETH seq=398234035 bid=0.0001294000 ask=0.0001378000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=SUI-USDT seq=7715687749 bid=0.9134000000 ask=0.9135000000
2026/04/12 09:35:25 [KuCoin WS 1] bootstrap complete symbol=TRAC-USDT seq=1072072633 bid=0.2897000000 ask=0.2906000000
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=TRX-USDT seq=1869321927 bid=0.3199000000 ask=0.3200000000
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=TRX-ETH seq=1200038867 bid=0.0001442800 ask=0.0001444200
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=TWT-USDT seq=909090262 bid=0.4169000000 ask=0.4178000000
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=USDT-BRL seq=1139045256 bid=5.0346000000 ask=5.0476000000
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=USDT-EUR seq=457759713 bid=0.8552000000 ask=0.8553000000
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap progress 40/80 last=USDT-EUR
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=VET-ETH seq=571321124 bid=0.0000030900 ask=0.0000031100
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=TEL-BTC seq=2901173358 bid=0.0000000293 ask=0.0000000296
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=TRAC-BTC seq=324999096 bid=0.0000040300 ask=0.0000040900
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=VET-USDT seq=2105061855 bid=0.0068600000 ask=0.0068700000
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=VSYS-BTC seq=160765725 bid=0.0000000032 ask=0.0000000032
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=VSYS-USDT seq=806221602 bid=0.0002300000 ask=0.0002301000
2026/04/12 09:35:26 [KuCoin WS 1] bootstrap complete symbol=VET-BTC seq=599279881 bid=0.0000000954 ask=0.0000000966
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=WAN-USDT seq=340588902 bid=0.0554500000 ask=0.0557300000
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=WAVES-BTC seq=465413611 bid=0.0000056900 ask=0.0000057600
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=WAVES-USDT seq=2035037166 bid=0.4102000000 ask=0.4109000000
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap progress 50/80 last=WAVES-USDT
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=TRX-BTC seq=978007403 bid=0.0000044630 ask=0.0000044660
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=WBTC-BTC seq=705056492 bid=0.9956500000 ask=0.9994400000
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=WBTC-USDT seq=219799219 bid=71347.8100000000 ask=71649.6700000000
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=WAN-BTC seq=272960122 bid=0.0000007750 ask=0.0000007850
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=WIN-BTC seq=278144429 bid=0.0000000003 ask=0.0000000003
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=WIN-TRX seq=291944649 bid=0.0000593000 ask=0.0000599000
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=XDC-BTC seq=546103105 bid=0.0000004210 ask=0.0000004280
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=XDC-USDT seq=1533458485 bid=0.0304200000 ask=0.0304500000
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=XLM-BTC seq=871817694 bid=0.0000021150 ask=0.0000021180
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=XLM-ETH seq=843168816 bid=0.0000683800 ask=0.0000685300
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap progress 60/80 last=XLM-ETH
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=XLM-USDT seq=2521863752 bid=0.1516000000 ask=0.1517000000
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=XMR-ETH seq=2731194708 bid=0.1543900000 ask=0.1544900000
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=XMR-USDT seq=13314510228 bid=342.3200000000 ask=342.3700000000
2026/04/12 09:35:27 [KuCoin WS 0] bootstrap complete symbol=A-BTC seq=168744760 bid=0.0000010900 ask=0.0000010970
2026/04/12 09:35:27 [KuCoin WS 1] bootstrap complete symbol=XRP-BTC seq=1480166774 bid=0.0000185500 ask=0.0000185700
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=XRP-ETH seq=1789496221 bid=0.0005999000 ask=0.0006000000
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=XRP-USDT seq=18978928841 bid=1.3296400000 ask=1.3296500000
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=XTZ-BTC seq=435765681 bid=0.0000048100 ask=0.0000048400
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=XTZ-USDT seq=1249577956 bid=0.3458000000 ask=0.3461000000
2026/04/12 09:35:28 [KuCoin WS 0] bootstrap complete symbol=AAVE-USDT seq=9002539060 bid=89.7600000000 ask=89.7610000000
2026/04/12 09:35:28 [KuCoin WS 0] bootstrap complete symbol=ADA-BTC seq=1105261501 bid=0.0000033900 ask=0.0000034000
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=WIN-USDT seq=937729655 bid=0.0000189900 ask=0.0000190500
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=TWT-BTC seq=146488032 bid=0.0000057600 ask=0.0000058900
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap progress 70/80 last=TWT-BTC
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=XYO-BTC seq=64518874 bid=0.0000000497 ask=0.0000000503
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=XYO-ETH seq=140275373 bid=0.0000016050 ask=0.0000016280
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=XYO-USDT seq=2119256510 bid=0.0035770000 ask=0.0035910000
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=XRP-KCS seq=1001175118 bid=0.1581400000 ask=0.1585800000
2026/04/12 09:35:28 [KuCoin WS 0] bootstrap complete symbol=ALGO-ETH seq=943903865 bid=0.0000476800 ask=0.0000478200
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=ZEC-BTC seq=1266620727 bid=0.0050642000 ask=0.0050719000
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=ZEC-USDT seq=3472401530 bid=363.1440000000 ask=363.2850000000
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=ZIL-ETH seq=680987366 bid=0.0000017360 ask=0.0000017420
2026/04/12 09:35:28 [KuCoin WS 1] bootstrap complete symbol=ZIL-USDT seq=1525727991 bid=0.0038530000 ask=0.0038570000
2026/04/12 09:35:29 [KuCoin WS 1] bootstrap complete symbol=XMR-BTC seq=2062947777 bid=0.0047750000 ask=0.0047800000
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=ANKR-BTC seq=338223993 bid=0.0000000728 ask=0.0000000739
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=ALGO-BTC seq=749828055 bid=0.0000014710 ask=0.0000014830
2026/04/12 09:35:29 [KuCoin WS 1] bootstrap complete symbol=XDC-ETH seq=651345414 bid=0.0000136900 ask=0.0000137500
2026/04/12 09:35:29 [KuCoin WS 1] bootstrap progress 80/80 last=XDC-ETH
2026/04/12 09:35:29 [KuCoin WS 1] bootstrap finished 80/80 in 7.236010397s
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=ANKR-USDT seq=800939280 bid=0.0052600000 ask=0.0052900000
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=ADA-KCS seq=527466296 bid=0.0289700000 ask=0.0290600000
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=AR-BTC seq=147347956 bid=0.0000236000 ask=0.0000237000
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=AR-USDT seq=2386091399 bid=1.6930000000 ask=1.6950000000
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap progress 10/126 last=AR-USDT
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=ATOM-BTC seq=1755317660 bid=0.0000243800 ask=0.0000244200
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=ATOM-ETH seq=1564375061 bid=0.0007880000 ask=0.0007900000
2026/04/12 09:35:29 [KuCoin WS 0] bootstrap complete symbol=ALGO-USDT seq=2602535400 bid=0.1057000000 ask=0.1058000000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=A-ETH seq=583854788 bid=0.0000352000 ask=0.0000412000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=AVA-BTC seq=340828463 bid=0.0000027900 ask=0.0000028200
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=AVA-ETH seq=427920780 bid=0.0000906000 ask=0.0000915000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=ATOM-USDT seq=6314464599 bid=1.7478000000 ask=1.7479000000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=A-USDT seq=143905859 bid=0.0783000000 ask=0.0784000000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=BCH-BTC seq=1261439320 bid=0.0059350000 ask=0.0059410000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=BCH-USDT seq=4654216600 bid=425.4700000000 ask=425.6200000000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap progress 20/126 last=BCH-USDT
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=BCHSV-BTC seq=974984401 bid=0.0002163000 ask=0.0002174000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=ADA-USDT seq=8277434853 bid=0.2430000000 ask=0.2431000000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=AVAX-BTC seq=1689081896 bid=0.0001263800 ask=0.0001265300
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=BCHSV-ETH seq=1101763087 bid=0.0069300000 ask=0.0073000000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=BCHSV-USDT seq=1534660453 bid=15.5100000000 ask=15.5600000000
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=BDX-BTC seq=331373366 bid=0.0000011120 ask=0.0000011170
2026/04/12 09:35:30 [KuCoin WS 0] bootstrap complete symbol=BDX-USDT seq=1921008018 bid=0.0798400000 ask=0.0798900000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=BNB-BTC seq=2283289986 bid=0.0083097000 ask=0.0083173000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=BNB-KCS seq=1205169513 bid=70.8161000000 ask=70.9826000000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=AVAX-USDT seq=9253122528 bid=9.0610000000 ask=9.0620000000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap progress 30/126 last=AVAX-USDT
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=BTC-EUR seq=2881059082 bid=61212.7900000000 ask=61349.7800000000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=AVA-USDT seq=1112970170 bid=0.2009000000 ask=0.2016000000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=CHZ-USDT seq=1957189025 bid=0.0381300000 ask=0.0381500000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=CHZ-BTC seq=500798222 bid=0.0000005302 ask=0.0000005351
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=BNB-USDT seq=11268418290 bid=595.6780000000 ask=595.6790000000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=CKB-BTC seq=103430024 bid=0.0000000200 ask=0.0000000206
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=AAVE-BTC seq=1082005395 bid=0.0012520000 ask=0.0012530000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=CRO-BTC seq=10919270820 bid=0.0000009600 ask=0.0000009690
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=CRO-USDT seq=1388325656 bid=0.0688800000 ask=0.0689100000
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap complete symbol=CSPR-ETH seq=260343380 bid=0.0000013370 ask=0.0000013590
2026/04/12 09:35:31 [KuCoin WS 0] bootstrap progress 40/126 last=CSPR-ETH
2026/04/12 09:35:32 [KuCoin WS 0] bootstrap complete symbol=DAG-ETH seq=712446002 bid=0.0000039900 ask=0.0000042100
2026/04/12 09:35:32 [KuCoin WS 0] bootstrap complete symbol=CSPR-USDT seq=624816687 bid=0.0029820000 ask=0.0029880000
2026/04/12 09:35:32 [KuCoin WS 0] bootstrap complete symbol=DASH-BTC seq=601294066 bid=0.0005868000 ask=0.0005879000
2026/04/12 09:35:32 [KuCoin WS 0] bootstrap complete symbol=COTI-USDT seq=1180256062 bid=0.0134900000 ask=0.0135100000
2026/04/12 09:35:32 [KuCoin WS 0] bootstrap complete symbol=DASH-ETH seq=602722743 bid=0.0189700000 ask=0.0190200000
2026/04/12 09:35:32 [KuCoin WS 0] bootstrap complete symbol=BTC-USDT seq=31661183074 bid=71663.1000000000 ask=71663.2000000000
2026/04/12 09:35:32 [KuCoin WS 0] bootstrap complete symbol=DOGE-BTC seq=1555221639 bid=0.0000012750 ask=0.0000012770
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=DOGE-KCS seq=848152600 bid=0.0108710000 ask=0.0108860000
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=CKB-USDT seq=1270020030 bid=0.0014530000 ask=0.0014560000
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=DOT-BTC seq=1709681934 bid=0.0000171000 ask=0.0000171200
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap progress 50/126 last=DOT-BTC
2026/04/12 09:35:33 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=COTI-BTC seq=416059622 bid=0.0000001877 ask=0.0000001896
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=DOT-KCS seq=909539827 bid=0.1457000000 ask=0.1467000000
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=DOGE-USDT seq=11682734531 bid=0.0914200000 ask=0.0914300000
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=DOT-USDT seq=8373368551 bid=1.2260000000 ask=1.2261000000
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=EGLD-BTC seq=307768561 bid=0.0000523500 ask=0.0000531200
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=ENJ-USDT seq=892167082 bid=0.0317100000 ask=0.0317800000
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=DAG-USDT seq=3818959957 bid=0.0089330000 ask=0.0089440000
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=ERG-BTC seq=191240787 bid=0.0000042400 ask=0.0000042600
2026/04/12 09:35:33 [KuCoin WS 0] bootstrap complete symbol=ERG-USDT seq=1277117529 bid=0.3058000000 ask=0.3069000000
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=DASH-USDT seq=1550491449 bid=42.0700000000 ask=42.1100000000
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap progress 60/126 last=DASH-USDT
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=ETC-BTC seq=1244700398 bid=0.0001147000 ask=0.0001154000
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=ETC-ETH seq=1354818624 bid=0.0037120000 ask=0.0037200000
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=ETH-BRL seq=409347943 bid=11008.0000000000 ask=11222.5600000000
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=ENJ-ETH seq=317097749 bid=0.0000142600 ask=0.0000144900
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=ETH-EUR seq=3228822961 bid=1895.4700000000 ask=1901.5900000000
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=EWT-BTC seq=322798790 bid=0.0000063530 ask=0.0000063740
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=EWT-USDT seq=726897652 bid=0.4549100000 ask=0.4554800000
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=FET-BTC seq=633674800 bid=0.0000032880 ask=0.0000033060
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=EGLD-USDT seq=1349613756 bid=3.7500000000 ask=3.7600000000
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=FET-ETH seq=755143312 bid=0.0001057200 ask=0.0001072700
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap progress 70/126 last=FET-ETH
2026/04/12 09:35:34 [KuCoin WS 0] bootstrap complete symbol=FET-USDT seq=3632579144 bid=0.2361000000 ask=0.2362000000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=ETC-USDT seq=3784737801 bid=8.2376000000 ask=8.2377000000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=HBAR-BTC seq=915677536 bid=0.0000012050 ask=0.0000012070
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=HBAR-USDT seq=5311733344 bid=0.0864100000 ask=0.0864200000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=ETH-USDT seq=20714393212 bid=2216.4200000000 ask=2216.4300000000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=HYPE-KCS seq=391400433 bid=4.8510000000 ask=4.8620000000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=HYPE-USDT seq=4416914356 bid=40.8090000000 ask=40.8160000000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=ICP-BTC seq=549942385 bid=0.0000343000 ask=0.0000343600
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=ICP-USDT seq=3333539059 bid=2.4560000000 ask=2.4570000000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=ICX-ETH seq=127690621 bid=0.0000159000 ask=0.0000188000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap progress 80/126 last=ICX-ETH
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=ETH-BTC seq=4759420277 bid=0.0309200000 ask=0.0309300000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=ICX-USDT seq=527131280 bid=0.0354400000 ask=0.0355000000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=INJ-BTC seq=792700877 bid=0.0000405600 ask=0.0000409000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=IOST-USDT seq=1004505172 bid=0.0010550000 ask=0.0010590000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=IOTA-BTC seq=262784913 bid=0.0000007780 ask=0.0000007890
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=IOTA-USDT seq=916246800 bid=0.0560000000 ask=0.0562000000
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=IOTX-BTC seq=149763620 bid=0.0000000635 ask=0.0000000644
2026/04/12 09:35:35 [KuCoin WS 0] bootstrap complete symbol=IOTX-ETH seq=216158250 bid=0.0000020640 ask=0.0000020780
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=IOTX-USDT seq=774735260 bid=0.0045800000 ask=0.0045900000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KAS-BTC seq=308908900 bid=0.0000004530 ask=0.0000004550
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap progress 90/126 last=KAS-BTC
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KCS-BTC seq=3180875875 bid=0.0001172000 ask=0.0001173000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KCS-ETH seq=1587657434 bid=0.0037880000 ask=0.0037920000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=IOST-ETH seq=379313984 bid=0.0000004760 ask=0.0000004780
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=INJ-USDT seq=4688074366 bid=2.9170000000 ask=2.9180000000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KLV-BTC seq=532030188 bid=0.0000000141 ask=0.0000000142
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KLV-TRX seq=461425346 bid=0.0031600000 ask=0.0031800000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KNC-BTC seq=218612268 bid=0.0000018500 ask=0.0000018600
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KNC-ETH seq=272164232 bid=0.0000598000 ask=0.0000601000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KCS-USDT seq=2805357073 bid=8.4020000000 ask=8.4070000000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KNC-USDT seq=1125313897 bid=0.1327000000 ask=0.1330000000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap progress 100/126 last=KNC-USDT
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KLV-USDT seq=1239218213 bid=0.0010130000 ask=0.0010140000
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=KRL-BTC seq=28940490 bid=0.0000020701 ask=0.0000020900
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=LINK-BTC seq=2528158323 bid=0.0001226500 ask=0.0001228100
2026/04/12 09:35:36 [KuCoin WS 0] bootstrap complete symbol=LINK-USDT seq=13123279941 bid=8.7942000000 ask=8.7943000000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=LTC-ETH seq=1640132167 bid=0.0243300000 ask=0.0243500000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=LTC-KCS seq=875795660 bid=6.4120000000 ask=6.4270000000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=KAS-USDT seq=2788990230 bid=0.0325800000 ask=0.0326000000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=KRL-USDT seq=970005808 bid=0.1494600000 ask=0.1496800000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=LYX-ETH seq=608143656 bid=0.0001099000 ask=0.0001102000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=LTC-USDT seq=8540945379 bid=53.9400000000 ask=53.9500000000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap progress 110/126 last=LTC-USDT
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=MANA-USDT seq=3223136761 bid=0.0879000000 ask=0.0880300000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=MANTRA-BTC seq=1738616 bid=0.0000001457 ask=0.0000001461
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=MANTRA-USDT seq=6498556 bid=0.0104300000 ask=0.0104700000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=BTC-BRL seq=415485636 bid=350019.4000000000 ask=365295.2000000000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=NEAR-BTC seq=1349437474 bid=0.0000186400 ask=0.0000186700
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=NEO-BTC seq=662028997 bid=0.0000389000 ask=0.0000392000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=NEO-USDT seq=2154972928 bid=2.7961000000 ask=2.7990000000
2026/04/12 09:35:37 [KuCoin WS 0] bootstrap complete symbol=NFT-TRX seq=379425025 bid=0.0000010350 ask=0.0000010430
2026/04/12 09:35:38 [KuCoin WS 0] bootstrap complete symbol=OGN-BTC seq=191484219 bid=0.0000002950 ask=0.0000003040
2026/04/12 09:35:38 [KuCoin WS 0] bootstrap complete symbol=ONE-BTC seq=498886003 bid=0.0000000286 ask=0.0000000289
2026/04/12 09:35:38 [KuCoin WS 0] bootstrap progress 120/126 last=ONE-BTC
2026/04/12 09:35:38 [KuCoin WS 0] bootstrap complete symbol=NFT-USDT seq=1155263622 bid=0.0000003321 ask=0.0000003328
2026/04/12 09:35:38 [KuCoin WS 0] bootstrap complete symbol=MANA-ETH seq=645719392 bid=0.0000396000 ask=0.0000398000
2026/04/12 09:35:38 [KuCoin WS 0] bootstrap complete symbol=NEAR-USDT seq=6887151423 bid=1.3365000000 ask=1.3368000000
2026/04/12 09:35:38 [KuCoin WS 0] bootstrap complete symbol=LTC-BTC seq=1656819419 bid=0.0007520000 ask=0.0007530000
2026/04/12 09:35:39 [KuCoin WS 0] bootstrap complete symbol=OGN-USDT seq=797217112 bid=0.0211500000 ask=0.0211800000
2026/04/12 09:35:39 [KuCoin WS 0] bootstrap complete symbol=LYX-USDT seq=1658807280 bid=0.2431000000 ask=0.2438000000
2026/04/12 09:35:39 [KuCoin WS 0] bootstrap progress 126/126 last=LYX-USDT
2026/04/12 09:35:39 [KuCoin WS 0] bootstrap finished 126/126 in 12.247447944s
2026/04/12 09:35:43 [Calculator] summary checked=3981 written=0 profitable=0 best_pct=-0.2240% best_usdt=-0.005051 best_tri=USDT->ICP->BTC
2026/04/12 09:35:53 [Calculator] summary checked=12044 written=0 profitable=0 best_pct=-0.2240% best_usdt=-0.005051 best_tri=USDT->ICP->BTC
2026/04/12 09:36:03 [Calculator] summary checked=15075 written=0 profitable=0 best_pct=-0.2240% best_usdt=-0.005051 best_tri=USDT->ICP->BTC
2026/04/12 09:36:13 [Calculator] summary checked=17099 written=0 profitable=0 best_pct=-0.2240% best_usdt=-0.005051 best_tri=USDT->ICP->BTC
2026/04/12 09:36:23 [Calculator] summary checked=18329 written=0 profitable=0 best_pct=-0.2119% best_usdt=-0.004778 best_tri=USDT->ICP->BTC
2026/04/12 09:36:33 [Calculator] summary checked=19345 written=0 profitable=0 best_pct=-0.2119% best_usdt=-0.004778 best_tri=USDT->ICP->BTC
2026/04/12 09:36:43 [Calculator] summary checked=20800 written=0 profitable=0 best_pct=-0.2119% best_usdt=-0.004778 best_tri=USDT->ICP->BTC
2026/04/12 09:36:53 [Calculator] summary checked=21483 written=0 profitable=0 best_pct=-0.2119% best_usdt=-0.004778 best_tri=USDT->ICP->BTC
2026/04/12 09:37:03 [Calculator] summary checked=23236 written=0 profitable=0 best_pct=-0.2119% best_usdt=-0.004778 best_tri=USDT->ICP->BTC
2026/04/12 09:37:13 [Calculator] summary checked=24119 written=0 profitable=0 best_pct=-0.2119% best_usdt=-0.004778 best_tri=USDT->ICP->BTC
2026/04/12 09:37:23 [Calculator] summary checked=25559 written=0 profitable=0 best_pct=-0.2119% best_usdt=-0.004778 best_tri=USDT->ICP->BTC
2026/04/12 09:37:33 [Calculator] summary checked=29161 written=0 profitable=0 best_pct=-0.2119% best_usdt=-0.004778 best_tri=USDT->ICP->BTC
2026/04/12 09:37:43 [Calculator] summary checked=31577 written=0 profitable=0 best_pct=-0.1920% best_usdt=-0.004328 best_tri=USDT->ICP->BTC
2026/04/12 09:37:53 [Calculator] summary checked=35752 written=0 profitable=0 best_pct=-0.1920% best_usdt=-0.004328 best_tri=USDT->ICP->BTC
2026/04/12 09:38:03 [Calculator] summary checked=37858 written=0 profitable=0 best_pct=-0.1920% best_usdt=-0.004328 best_tri=USDT->ICP->BTC
2026/04/12 09:38:13 [Calculator] summary checked=40965 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:38:23 [Calculator] summary checked=44022 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:38:33 [Calculator] summary checked=45073 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:38:43 [Calculator] summary checked=47081 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:38:53 [Calculator] summary checked=48412 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:39:03 [Calculator] summary checked=51016 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:39:13 [Calculator] summary checked=51892 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:39:23 [Calculator] summary checked=53320 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:39:33 [Calculator] summary checked=55545 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:39:43 [Calculator] summary checked=58700 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:39:53 [Calculator] summary checked=60623 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:40:03 [Calculator] summary checked=62044 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:40:13 [Calculator] summary checked=63508 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:40:23 [Calculator] summary checked=65053 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:40:33 [Calculator] summary checked=65882 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:40:43 [Calculator] summary checked=66731 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:40:53 [Calculator] summary checked=67249 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:41:03 [Calculator] summary checked=67780 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:41:13 [Calculator] summary checked=70315 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:41:23 [Calculator] summary checked=70473 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:41:33 [Calculator] summary checked=70554 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:41:43 [Calculator] summary checked=75058 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:41:53 [Calculator] summary checked=75952 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:42:03 [Calculator] summary checked=77233 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:42:13 [Calculator] summary checked=77410 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:42:23 [Calculator] summary checked=78605 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:42:33 [Calculator] summary checked=79396 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:42:43 [Calculator] summary checked=79748 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:42:53 [Calculator] summary checked=80162 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:43:03 [Calculator] summary checked=81100 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:43:13 [Calculator] summary checked=81850 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:43:23 [Calculator] summary checked=82616 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:43:33 [Calculator] summary checked=84191 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:43:43 [Calculator] summary checked=84660 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:43:53 [Calculator] summary checked=85081 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:44:03 [Calculator] summary checked=86452 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:44:13 [Calculator] summary checked=86778 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:44:23 [Calculator] summary checked=87390 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:44:33 [Calculator] summary checked=87670 written=0 profitable=0 best_pct=-0.1018% best_usdt=-0.001220 best_tri=USDT->BTC->ERG
2026/04/12 09:44:43 [Calculator] summary checked=87788 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:44:53 [Calculator] summary checked=88670 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:45:03 [Calculator] summary checked=89172 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:45:13 [Calculator] summary checked=91110 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:45:23 [Calculator] summary checked=91646 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:45:33 [Calculator] summary checked=91963 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:45:43 [Calculator] summary checked=92276 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:45:53 [Calculator] summary checked=92614 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:46:03 [Calculator] summary checked=93425 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:46:13 [Calculator] summary checked=94098 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:46:23 [Calculator] summary checked=95098 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:46:33 [Calculator] summary checked=96742 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:46:43 [Calculator] summary checked=97617 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:46:53 [Calculator] summary checked=98091 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:47:03 [Calculator] summary checked=100611 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:47:13 [Calculator] summary checked=101449 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:47:23 [Calculator] summary checked=102131 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:47:33 [Calculator] summary checked=102821 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:47:43 [Calculator] summary checked=103569 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:47:53 [Calculator] summary checked=104725 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:48:03 [Calculator] summary checked=109396 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:48:13 [Calculator] summary checked=110018 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:48:23 [Calculator] summary checked=110666 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:48:33 [Calculator] summary checked=117316 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:48:43 [Calculator] summary checked=118878 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:48:53 [Calculator] summary checked=120427 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:49:03 [Calculator] summary checked=120895 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:49:13 [Calculator] summary checked=121515 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:49:23 [Calculator] summary checked=124790 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:49:33 [Calculator] summary checked=126200 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:49:43 [Calculator] summary checked=127600 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:49:53 [Calculator] summary checked=127772 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:50:03 [Calculator] summary checked=128620 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:50:13 [Calculator] summary checked=131035 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:50:23 [Calculator] summary checked=139767 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:50:33 [Calculator] summary checked=140560 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:50:43 [Calculator] summary checked=141188 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:50:53 [Calculator] summary checked=142904 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:51:03 [Calculator] summary checked=144107 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:51:13 [Calculator] summary checked=146888 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:51:23 [Calculator] summary checked=156684 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:51:33 [Calculator] summary checked=158123 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:51:43 [Calculator] summary checked=161968 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:51:53 [Calculator] summary checked=164269 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:52:03 [Calculator] summary checked=165367 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:52:13 [Calculator] summary checked=167085 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:52:23 [Calculator] summary checked=168022 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:52:33 [Calculator] summary checked=170393 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:52:43 [Calculator] summary checked=171771 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:52:53 [Calculator] summary checked=171847 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:53:03 [Calculator] summary checked=172337 written=0 profitable=0 best_pct=-0.0864% best_usdt=-0.010910 best_tri=USDT->BDX->BTC
2026/04/12 09:53:13 [Calculator] summary checked=172754 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:53:23 [Calculator] summary checked=173795 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:53:33 [Calculator] summary checked=176555 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:53:43 [Calculator] summary checked=177602 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:53:53 [Calculator] summary checked=178128 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:54:03 [Calculator] summary checked=182474 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:54:14 [Calculator] summary checked=183075 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:54:24 [Calculator] summary checked=184506 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:54:34 [Calculator] summary checked=185241 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:54:44 [Calculator] summary checked=186653 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:54:54 [Calculator] summary checked=187182 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:55:04 [Calculator] summary checked=188334 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:55:14 [Calculator] summary checked=192384 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:55:24 [Calculator] summary checked=193405 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:55:34 [Calculator] summary checked=194181 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:55:44 [Calculator] summary checked=195180 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:55:54 [Calculator] summary checked=195627 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:56:04 [Calculator] summary checked=196096 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:56:14 [Calculator] summary checked=196494 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:56:24 [Calculator] summary checked=196890 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:56:34 [Calculator] summary checked=197092 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:56:44 [Calculator] summary checked=197544 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:56:54 [Calculator] summary checked=202103 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:57:04 [Calculator] summary checked=203171 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:57:14 [Calculator] summary checked=207166 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:57:24 [Calculator] summary checked=208431 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:57:34 [Calculator] summary checked=209007 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:57:44 [Calculator] summary checked=209359 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:57:54 [Calculator] summary checked=213107 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:58:04 [Calculator] summary checked=213792 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:58:14 [Calculator] summary checked=214707 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:58:24 [Calculator] summary checked=215395 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:58:34 [Calculator] summary checked=216192 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:58:44 [Calculator] summary checked=217277 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:58:54 [Calculator] summary checked=217843 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:59:04 [Calculator] summary checked=219558 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:59:14 [Calculator] summary checked=220057 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:59:24 [Calculator] summary checked=220309 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:59:34 [Calculator] summary checked=220563 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:59:44 [Calculator] summary checked=221101 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 09:59:54 [Calculator] summary checked=222046 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:00:04 [Calculator] summary checked=237721 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:00:14 [Calculator] summary checked=245132 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:00:24 [Calculator] summary checked=245677 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:00:34 [Calculator] summary checked=249741 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:00:44 [Calculator] summary checked=253656 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:00:54 [Calculator] summary checked=254510 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:01:04 [Calculator] summary checked=259326 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:01:14 [Calculator] summary checked=264649 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:01:24 [Calculator] summary checked=268364 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:01:34 [Calculator] summary checked=272037 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:01:44 [Calculator] summary checked=274118 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:01:54 [Calculator] summary checked=277209 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:02:04 [Calculator] summary checked=279484 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:02:14 [Calculator] summary checked=283358 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:02:24 [Calculator] summary checked=285126 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:02:34 [Calculator] summary checked=286576 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:02:44 [Calculator] summary checked=287523 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:02:54 [Calculator] summary checked=291873 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:03:04 [Calculator] summary checked=294167 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:03:14 [Calculator] summary checked=297701 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:03:24 [Calculator] summary checked=304260 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:03:34 [Calculator] summary checked=304942 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:03:44 [Calculator] summary checked=307416 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:03:54 [Calculator] summary checked=314942 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:04:04 [Calculator] summary checked=317081 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:04:14 [Calculator] summary checked=320857 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:04:24 [Calculator] summary checked=322966 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:04:34 [Calculator] summary checked=323530 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:04:44 [Calculator] summary checked=325187 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:04:54 [Calculator] summary checked=326239 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:05:04 [Calculator] summary checked=330789 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:05:14 [Calculator] summary checked=335072 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:05:24 [Calculator] summary checked=340821 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:05:34 [Calculator] summary checked=344787 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:05:44 [Calculator] summary checked=346058 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:05:54 [Calculator] summary checked=349496 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:06:04 [Calculator] summary checked=349739 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:06:14 [Calculator] summary checked=350299 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:06:24 [Calculator] summary checked=351099 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:06:34 [Calculator] summary checked=354275 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:06:44 [Calculator] summary checked=355399 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:06:54 [Calculator] summary checked=356464 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:07:04 [Calculator] summary checked=358488 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:07:14 [Calculator] summary checked=359197 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:07:24 [Calculator] summary checked=360420 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:07:34 [Calculator] summary checked=360721 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:07:44 [Calculator] summary checked=362208 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:07:54 [Calculator] summary checked=365844 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:08:04 [Calculator] summary checked=366372 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:08:14 [Calculator] summary checked=367852 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:08:24 [Calculator] summary checked=368599 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:08:34 [Calculator] summary checked=371688 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:08:44 [Calculator] summary checked=373674 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:08:54 [Calculator] summary checked=374298 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:09:04 [Calculator] summary checked=377870 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:09:14 [Calculator] summary checked=378428 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:09:24 [Calculator] summary checked=386635 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:09:34 [Calculator] summary checked=390170 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:09:44 [Calculator] summary checked=393044 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:09:54 [Calculator] summary checked=394268 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:10:04 [Calculator] summary checked=397932 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:10:14 [Calculator] summary checked=398332 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:10:24 [Calculator] summary checked=401437 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:10:34 [Calculator] summary checked=405833 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:10:44 [Calculator] summary checked=407673 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:10:54 [Calculator] summary checked=409033 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:11:04 [Calculator] summary checked=410732 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:11:14 [Calculator] summary checked=411219 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:11:24 [Calculator] summary checked=411649 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:11:34 [Calculator] summary checked=415442 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:11:44 [Calculator] summary checked=416568 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:11:54 [Calculator] summary checked=417872 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:12:04 [Calculator] summary checked=419995 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:12:14 [Calculator] summary checked=420863 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:12:24 [Calculator] summary checked=422777 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:12:34 [Calculator] summary checked=426514 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:12:44 [Calculator] summary checked=427302 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:12:54 [Calculator] summary checked=429319 written=4 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:13:04 [Calculator] summary checked=430654 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:13:14 [Calculator] summary checked=431816 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:13:24 [Calculator] summary checked=432733 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:13:34 [Calculator] summary checked=434817 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:13:44 [Calculator] summary checked=435669 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:13:54 [Calculator] summary checked=438552 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:14:04 [Calculator] summary checked=440152 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:14:14 [Calculator] summary checked=442933 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:14:24 [Calculator] summary checked=444083 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:14:34 [Calculator] summary checked=449914 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:14:44 [Calculator] summary checked=461932 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:14:54 [Calculator] summary checked=463070 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:15:04 [Calculator] summary checked=465379 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:15:14 [Calculator] summary checked=467434 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:15:24 [Calculator] summary checked=468312 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:15:34 [Calculator] summary checked=470726 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:15:44 [Calculator] summary checked=471313 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:15:54 [Calculator] summary checked=472087 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:16:04 [Calculator] summary checked=474355 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:16:14 [Calculator] summary checked=477536 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:16:24 [Calculator] summary checked=479980 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:16:34 [Calculator] summary checked=481585 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:16:44 [Calculator] summary checked=483296 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:16:54 [Calculator] summary checked=484243 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:17:04 [Calculator] summary checked=486562 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:17:14 [Calculator] summary checked=487478 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:17:24 [Calculator] summary checked=491490 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:17:34 [Calculator] summary checked=492745 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:17:44 [Calculator] summary checked=494297 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:17:54 [Calculator] summary checked=495963 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:18:04 [Calculator] summary checked=496818 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:18:14 [Calculator] summary checked=499257 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:18:24 [Calculator] summary checked=502237 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:18:34 [Calculator] summary checked=510079 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:18:44 [Calculator] summary checked=513174 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:18:54 [Calculator] summary checked=513647 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:19:04 [Calculator] summary checked=516919 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:19:14 [Calculator] summary checked=518297 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:19:24 [Calculator] summary checked=523816 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:19:34 [Calculator] summary checked=526305 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:19:44 [Calculator] summary checked=528885 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:19:54 [Calculator] summary checked=529636 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:20:04 [Calculator] summary checked=532959 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:20:14 [Calculator] summary checked=534055 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:20:24 [Calculator] summary checked=534510 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:20:34 [Calculator] summary checked=537045 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:20:44 [Calculator] summary checked=538546 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:20:54 [Calculator] summary checked=538933 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:21:04 [Calculator] summary checked=540317 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:21:14 [Calculator] summary checked=540910 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:21:24 [Calculator] summary checked=542419 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:21:34 [Calculator] summary checked=546019 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:21:44 [Calculator] summary checked=547046 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:21:54 [Calculator] summary checked=547459 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:22:04 [Calculator] summary checked=548963 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:22:14 [Calculator] summary checked=549704 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:22:24 [Calculator] summary checked=550709 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:22:34 [Calculator] summary checked=552762 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:22:44 [Calculator] summary checked=555153 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:22:54 [Calculator] summary checked=556984 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:23:04 [Calculator] summary checked=560009 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:23:14 [Calculator] summary checked=561352 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:23:24 [Calculator] summary checked=562195 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:23:34 [Calculator] summary checked=562927 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:23:44 [Calculator] summary checked=563893 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:23:54 [Calculator] summary checked=567352 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:24:04 [Calculator] summary checked=571521 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:24:14 [Calculator] summary checked=572906 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:24:24 [Calculator] summary checked=574670 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:24:34 [Calculator] summary checked=575456 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:24:44 [Calculator] summary checked=576424 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:24:54 [Calculator] summary checked=577592 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:25:04 [Calculator] summary checked=579762 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:25:14 [Calculator] summary checked=580509 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:25:24 [Calculator] summary checked=582168 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:25:34 [Calculator] summary checked=584559 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:25:44 [Calculator] summary checked=585516 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:25:54 [Calculator] summary checked=586349 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:26:04 [Calculator] summary checked=588524 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:26:14 [Calculator] summary checked=589511 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:26:24 [Calculator] summary checked=593101 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:26:35 [Calculator] summary checked=594003 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:26:45 [Calculator] summary checked=594118 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:26:55 [Calculator] summary checked=594384 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:27:05 [Calculator] summary checked=596310 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:27:15 [Calculator] summary checked=596937 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:27:25 [Calculator] summary checked=597433 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:27:35 [Calculator] summary checked=598760 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:27:45 [Calculator] summary checked=600356 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:27:55 [Calculator] summary checked=600840 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:28:05 [Calculator] summary checked=601364 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:28:15 [Calculator] summary checked=601886 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:28:25 [Calculator] summary checked=602656 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:28:35 [Calculator] summary checked=603523 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:28:45 [Calculator] summary checked=603958 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:28:55 [Calculator] summary checked=604672 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:29:05 [Calculator] summary checked=606660 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:29:15 [Calculator] summary checked=607482 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:29:25 [Calculator] summary checked=609122 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:29:35 [Calculator] summary checked=610511 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:29:45 [Calculator] summary checked=610867 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:29:55 [Calculator] summary checked=611086 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:30:05 [Calculator] summary checked=612709 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:30:15 [Calculator] summary checked=613204 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:30:25 [Calculator] summary checked=614422 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:30:35 [Calculator] summary checked=617815 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:30:45 [Calculator] summary checked=618900 written=5 profitable=0 best_pct=-0.0186% best_usdt=-0.000916 best_tri=USDT->BDX->BTC
2026/04/12 10:30:55 [Calculator] summary checked=620336 written=10 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:31:05 [Calculator] summary checked=622761 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:31:15 [Calculator] summary checked=623570 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:31:25 [Calculator] summary checked=623972 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:31:35 [Calculator] summary checked=624835 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:31:45 [Calculator] summary checked=625588 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:31:55 [Calculator] summary checked=626715 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:32:05 [Calculator] summary checked=628355 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:32:15 [Calculator] summary checked=628697 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:32:25 [Calculator] summary checked=629087 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:32:35 [Calculator] summary checked=629919 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:32:45 [Calculator] summary checked=630788 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:32:55 [Calculator] summary checked=634254 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:33:05 [Calculator] summary checked=635024 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:33:15 [Calculator] summary checked=639092 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:33:25 [Calculator] summary checked=639691 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:33:35 [Calculator] summary checked=640862 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:33:45 [Calculator] summary checked=641578 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:33:55 [Calculator] summary checked=642354 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:34:05 [Calculator] summary checked=643742 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:34:15 [Calculator] summary checked=644450 written=12 profitable=0 best_pct=-0.0040% best_usdt=-0.009437 best_tri=USDT->BCHSV->BTC
2026/04/12 10:34:17 [ARB] USDT→BTC→WAN | 0.1157% | volume=66.35 USDT | profit=0.076743 USDT | anchor=1775979257694 | l1=BTC-USDT BUY out=0.00092532 age=48ms | l2=WAN-BTC BUY out=1162.99226128 age=22ms | l3=WAN-USDT SELL out=66.42345284 age=0ms
2026/04/12 10:34:17 [ARB] USDT→BTC→WAN | 0.1157% | volume=66.35 USDT | profit=0.076743 USDT | anchor=1775979257694 | l1=BTC-USDT BUY out=0.00092532 age=48ms | l2=WAN-BTC BUY out=1162.99226128 age=0ms | l3=WAN-USDT SELL out=66.42345284 age=0ms
2026/04/12 10:34:17 [ARB] USDT→BTC→WAN | 0.1157% | volume=66.35 USDT | profit=0.076743 USDT | anchor=1775979257701 | l1=BTC-USDT BUY out=0.00092532 age=0ms | l2=WAN-BTC BUY out=1162.99226128 age=7ms | l3=WAN-USDT SELL out=66.42345284 age=7ms
2026/04/12 10:34:17 [ARB] USDT→BTC→WAN | 0.1157% | volume=66.35 USDT | profit=0.076743 USDT | anchor=1775979257700 | l1=BTC-USDT BUY out=0.00092532 age=54ms | l2=WAN-BTC BUY out=1162.99226128 age=0ms | l3=WAN-USDT SELL out=66.42345284 age=6ms
2026/04/12 10:34:25 [Calculator] summary checked=645653 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:34:35 [Calculator] summary checked=646325 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:34:45 [Calculator] summary checked=646878 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:34:55 [Calculator] summary checked=647496 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:35:05 [Calculator] summary checked=648083 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:35:15 [Calculator] summary checked=648684 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:35:25 [Calculator] summary checked=649367 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:35:35 [Calculator] summary checked=650299 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:35:45 [Calculator] summary checked=666796 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:35:55 [Calculator] summary checked=669298 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:36:05 [Calculator] summary checked=677878 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:36:15 [Calculator] summary checked=680858 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:36:25 [Calculator] summary checked=684338 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:36:35 [Calculator] summary checked=685533 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:36:45 [Calculator] summary checked=687776 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:36:55 [Calculator] summary checked=689751 written=16 profitable=4 best_pct=0.1157% best_usdt=0.076743 best_tri=USDT->BTC->WAN
2026/04/12 10:37:05 [Calculator] summary checked=692022 written=16 profitable=4 be









