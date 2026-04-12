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





gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
2026/04/12 09:11:55 pprof on http://localhost:6060/debug/pprof/
2026/04/12 09:11:55 [KuCoin] started with 2 WS
2026/04/12 09:11:55 [Main] KuCoinCollector started
2026/04/12 09:11:55 [Calculator] indexed 206 symbols
2026/04/12 09:11:56 [KuCoin WS 0] connected
2026/04/12 09:11:56 [KuCoin WS 1] connected
2026/04/12 09:12:06 [KuCoin WS 1] bootstrap complete symbol=ONT-BTC seq=599180861 bid=0.0000012080 ask=0.0000012200
2026/04/12 09:12:06 [KuCoin WS 1] bootstrap complete symbol=ONE-USDT seq=2149427738 bid=0.0020520000 ask=0.0020620000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=PAXG-USDT seq=1638391345 bid=4708.5300000000 ask=4708.5400000000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=POND-BTC seq=123772376 bid=0.0000000302 ask=0.0000000331
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=POND-USDT seq=434246626 bid=0.0021880000 ask=0.0022020000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=ONT-ETH seq=618948897 bid=0.0000390600 ask=0.0000442000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=PEPE-USDT seq=5048138119 bid=0.0000035280 ask=0.0000035290
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=RLC-BTC seq=144020811 bid=0.0000057900 ask=0.0000058200
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=ONT-USDT seq=1628832441 bid=0.0866600000 ask=0.0869400000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=RSR-BTC seq=132295411 bid=0.0000000205 ask=0.0000000207
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap progress 10/80 last=RSR-BTC
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=RLC-USDT seq=1052836520 bid=0.4158000000 ask=0.4166000000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=RSR-USDT seq=1855612347 bid=0.0014810000 ask=0.0014820000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=RUNE-BTC seq=1114756784 bid=0.0000054200 ask=0.0000054700
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=RUNE-USDT seq=3496762669 bid=0.3914000000 ask=0.3917000000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=SCRT-BTC seq=69769370 bid=0.0000012740 ask=0.0000012770
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=SHIB-DOGE seq=676603715 bid=0.0000637900 ask=0.0000640000
2026/04/12 09:12:07 [KuCoin WS 1] bootstrap complete symbol=PAXG-BTC seq=361701101 bid=0.0655950000 ask=0.0656820000
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=SHIB-USDT seq=5998501512 bid=0.0000058480 ask=0.0000058490
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=SNX-USDT seq=1536431520 bid=0.2852000000 ask=0.2861000000
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=SOL-KCS seq=458282158 bid=9.7990000000 ask=9.8300000000
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap progress 20/80 last=SOL-KCS
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=SOL-USDT seq=19083253488 bid=82.4300000000 ask=82.4400000000
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=STORJ-ETH seq=100783311 bid=0.0000441000 ask=0.0000443000
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=STX-BTC seq=442994654 bid=0.0000029800 ask=0.0000029900
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=STX-USDT seq=3150664867 bid=0.2144000000 ask=0.2145000000
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=PEPE-KCS seq=146635076 bid=0.0000004182 ask=0.0000004221
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=SUI-USDT seq=7715624116 bid=0.9162000000 ask=0.9163000000
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=TEL-ETH seq=1361297603 bid=0.0000009456 ask=0.0000009840
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=SNX-ETH seq=505529443 bid=0.0001280000 ask=0.0001350000
2026/04/12 09:12:08 [KuCoin WS 1] bootstrap complete symbol=SCRT-USDT seq=341461778 bid=0.0915000000 ask=0.0916000000
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=TEL-USDT seq=5957432975 bid=0.0021030000 ask=0.0021070000
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap progress 30/80 last=TEL-USDT
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=TRAC-BTC seq=324998565 bid=0.0000040200 ask=0.0000040900
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=TRX-BTC seq=978002849 bid=0.0000044630 ask=0.0000044660
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=TRAC-USDT seq=1072069699 bid=0.2893000000 ask=0.2905000000
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=TEL-BTC seq=2901170457 bid=0.0000000293 ask=0.0000000296
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=STORJ-USDT seq=901050923 bid=0.0979000000 ask=0.0983000000
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=TRX-USDT seq=1869311865 bid=0.3202000000 ask=0.3203000000
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=TWT-BTC seq=146487320 bid=0.0000057200 ask=0.0000058000
2026/04/12 09:12:09 [KuCoin WS 1] bootstrap complete symbol=TRAC-ETH seq=398232859 bid=0.0001294000 ask=0.0001368000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=USDT-EUR seq=457758862 bid=0.8550000000 ask=0.8552000000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=TWT-USDT seq=909086631 bid=0.4146000000 ask=0.4147000000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap progress 40/80 last=TWT-USDT
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=VET-ETH seq=571320393 bid=0.0000030900 ask=0.0000031100
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=USDT-BRL seq=1139005793 bid=5.0311000000 ask=5.0425000000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=SUI-KCS seq=243393244 bid=0.1087000000 ask=0.1096000000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=TRX-ETH seq=1200035189 bid=0.0001442700 ask=0.0001443700
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=VET-USDT seq=2105057886 bid=0.0068700000 ask=0.0068800000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=VSYS-BTC seq=160762276 bid=0.0000000032 ask=0.0000000032
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=VSYS-USDT seq=806215130 bid=0.0002299000 ask=0.0002301000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=WAN-BTC seq=272957209 bid=0.0000007830 ask=0.0000008000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=WAVES-BTC seq=465413201 bid=0.0000057000 ask=0.0000057500
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=WAVES-USDT seq=2035036074 bid=0.4102000000 ask=0.4109000000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap progress 50/80 last=WAVES-USDT
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=WBTC-BTC seq=705055021 bid=0.9956500000 ask=0.9994400000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=SNX-BTC seq=424271469 bid=0.0000039500 ask=0.0000041800
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=WBTC-USDT seq=219798840 bid=71413.7300000000 ask=71679.2100000000
2026/04/12 09:12:10 [KuCoin WS 1] bootstrap complete symbol=WIN-BTC seq=278143967 bid=0.0000000003 ask=0.0000000003
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=WIN-TRX seq=291943885 bid=0.0000590000 ask=0.0000597000
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=WIN-USDT seq=937727268 bid=0.0000189400 ask=0.0000189900
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=VET-BTC seq=599279052 bid=0.0000000954 ask=0.0000000964
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XDC-BTC seq=546102342 bid=0.0000004200 ask=0.0000004290
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XDC-USDT seq=1533456565 bid=0.0304200000 ask=0.0304700000
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XLM-BTC seq=871815058 bid=0.0000021160 ask=0.0000021200
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap progress 60/80 last=XLM-BTC
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XLM-ETH seq=843167029 bid=0.0000684000 ask=0.0000685400
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XLM-USDT seq=2521859909 bid=0.1519000000 ask=0.1520000000
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XDC-ETH seq=651344453 bid=0.0000136400 ask=0.0000137700
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XMR-ETH seq=2731134987 bid=0.1535800000 ask=0.1537700000
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XRP-BTC seq=1480164272 bid=0.0000185500 ask=0.0000185600
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XMR-USDT seq=13314105764 bid=341.0500000000 ask=341.0800000000
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XRP-ETH seq=1789492472 bid=0.0005997000 ask=0.0006001000
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XMR-BTC seq=2062907772 bid=0.0047520000 ask=0.0047560000
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XRP-KCS seq=1001165663 bid=0.1583300000 ask=0.1585900000
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap complete symbol=XTZ-BTC seq=435765093 bid=0.0000048100 ask=0.0000048400
2026/04/12 09:12:11 [KuCoin WS 1] bootstrap progress 70/80 last=XTZ-BTC
2026/04/12 09:12:12 [KuCoin WS 1] bootstrap complete symbol=ZEC-BTC seq=1266585481 bid=0.0050329000 ask=0.0050384000
2026/04/12 09:12:12 [KuCoin WS 1] bootstrap complete symbol=ZEC-USDT seq=3472328975 bid=361.2070000000 ask=361.2970000000
2026/04/12 09:12:12 [KuCoin WS 1] bootstrap complete symbol=XTZ-USDT seq=1249576112 bid=0.3461000000 ask=0.3464000000
2026/04/12 09:12:12 [KuCoin WS 1] bootstrap complete symbol=WAN-USDT seq=340585360 bid=0.0563900000 ask=0.0569400000
2026/04/12 09:12:12 [KuCoin WS 0] bootstrap complete symbol=A-ETH seq=583851631 bid=0.0000350000 ask=0.0000443000
2026/04/12 09:12:12 [KuCoin WS 1] bootstrap complete symbol=XYO-ETH seq=140274798 bid=0.0000016030 ask=0.0000016250
2026/04/12 09:12:12 [KuCoin WS 0] bootstrap complete symbol=A-BTC seq=168742336 bid=0.0000010890 ask=0.0000010990
2026/04/12 09:12:12 [KuCoin WS 1] bootstrap complete symbol=XRP-USDT seq=18978859892 bid=1.3313400000 ask=1.3313500000
2026/04/12 09:12:12 [KuCoin WS 0] bootstrap complete symbol=ALGO-BTC seq=749826925 bid=0.0000014790 ask=0.0000014810
2026/04/12 09:12:12 [KuCoin WS 0] bootstrap complete symbol=ALGO-ETH seq=943902233 bid=0.0000476800 ask=0.0000478100
2026/04/12 09:12:13 [KuCoin WS 1] bootstrap complete symbol=ZIL-USDT seq=1525726212 bid=0.0038640000 ask=0.0038660000
2026/04/12 09:12:13 [KuCoin WS 0] bootstrap complete symbol=ANKR-BTC seq=338223389 bid=0.0000000728 ask=0.0000000739
2026/04/12 09:12:13 [KuCoin WS 0] bootstrap complete symbol=ALGO-USDT seq=2602532353 bid=0.1059000000 ask=0.1060000000
2026/04/12 09:12:13 [KuCoin WS 1] bootstrap complete symbol=XYO-USDT seq=2119255718 bid=0.0035770000 ask=0.0035920000
2026/04/12 09:12:13 [KuCoin WS 0] bootstrap complete symbol=ADA-USDT seq=8277424679 bid=0.2431000000 ask=0.2432000000
2026/04/12 09:12:13 [KuCoin WS 0] bootstrap complete symbol=AR-BTC seq=147347616 bid=0.0000236000 ask=0.0000238000
2026/04/12 09:12:13 [KuCoin WS 0] bootstrap complete symbol=AR-USDT seq=2386089048 bid=1.6980000000 ask=1.7000000000
2026/04/12 09:12:13 [KuCoin WS 1] bootstrap complete symbol=XYO-BTC seq=64518743 bid=0.0000000497 ask=0.0000000503
2026/04/12 09:12:13 [KuCoin WS 0] bootstrap complete symbol=ATOM-BTC seq=1755305410 bid=0.0000243800 ask=0.0000244200
2026/04/12 09:12:13 [KuCoin WS 0] bootstrap progress 10/126 last=ATOM-BTC
2026/04/12 09:12:13 [KuCoin WS 0] bootstrap complete symbol=ATOM-ETH seq=1564368089 bid=0.0007880000 ask=0.0007890000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=ATOM-USDT seq=6314445975 bid=1.7502000000 ask=1.7505000000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=ADA-KCS seq=527463544 bid=0.0289700000 ask=0.0290700000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=ADA-BTC seq=1105259524 bid=0.0000033900 ask=0.0000034000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=AVA-ETH seq=427919746 bid=0.0000899000 ask=0.0000917000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=AVA-USDT seq=1112968637 bid=0.2001000000 ask=0.2010000000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=AVAX-BTC seq=1689078011 bid=0.0001264800 ask=0.0001266700
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=ANKR-USDT seq=800937780 bid=0.0052600000 ask=0.0052900000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=BCH-BTC seq=1261434452 bid=0.0059270000 ask=0.0059330000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=AAVE-USDT seq=9002494961 bid=90.1590000000 ask=90.1600000000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap progress 20/126 last=AAVE-USDT
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=AAVE-BTC seq=1082001215 bid=0.0012560000 ask=0.0012580000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=BCH-USDT seq=4654195022 bid=425.4100000000 ask=425.4200000000
2026/04/12 09:12:14 [KuCoin WS 1] bootstrap complete symbol=ZIL-ETH seq=680986353 bid=0.0000017390 ask=0.0000017490
2026/04/12 09:12:14 [KuCoin WS 1] bootstrap progress 80/80 last=ZIL-ETH
2026/04/12 09:12:14 [KuCoin WS 1] bootstrap finished 80/80 in 8.268529624s
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=AVAX-USDT seq=9253060579 bid=9.0800000000 ask=9.0810000000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=BCHSV-ETH seq=1101761828 bid=0.0068900000 ask=0.0073100000
2026/04/12 09:12:14 [KuCoin WS 0] bootstrap complete symbol=BCHSV-BTC seq=974982978 bid=0.0002166000 ask=0.0002175000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=AVA-BTC seq=340827696 bid=0.0000027800 ask=0.0000028300
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=A-USDT seq=143902706 bid=0.0783000000 ask=0.0786000000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=BDX-BTC seq=331366839 bid=0.0000011120 ask=0.0000011170
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=BDX-USDT seq=1920956282 bid=0.0798500000 ask=0.0799000000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=BNB-BTC seq=2283284105 bid=0.0083000000 ask=0.0083088000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap progress 30/126 last=BNB-BTC
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=BCHSV-USDT seq=1534657313 bid=15.5500000000 ask=15.5800000000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=BNB-KCS seq=1205141092 bid=70.8161000000 ask=71.0042000000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=BTC-BRL seq=415484944 bid=350019.4000000000 ask=365295.2000000000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=BTC-EUR seq=2881003565 bid=61217.6800000000 ask=61452.8700000000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=CHZ-BTC seq=500797246 bid=0.0000005308 ask=0.0000005349
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=CHZ-USDT seq=1957184934 bid=0.0382900000 ask=0.0383000000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=CKB-BTC seq=103429700 bid=0.0000000200 ask=0.0000000206
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=CKB-USDT seq=1270019213 bid=0.0014540000 ask=0.0014570000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=BTC-USDT seq=31660808716 bid=71736.8000000000 ask=71736.9000000000
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=CRO-BTC seq=10919269778 bid=0.0000009590 ask=0.0000009690
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap progress 40/126 last=CRO-BTC
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=CSPR-ETH seq=260338219 bid=0.0000013330 ask=0.0000014180
2026/04/12 09:12:15 [KuCoin WS 0] bootstrap complete symbol=CSPR-USDT seq=624783131 bid=0.0029830000 ask=0.0029880000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DAG-USDT seq=3818904434 bid=0.0089420000 ask=0.0089580000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DASH-BTC seq=601280288 bid=0.0005853000 ask=0.0005865000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DASH-ETH seq=602717483 bid=0.0188900000 ask=0.0189600000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=BNB-USDT seq=11268371831 bid=595.7690000000 ask=595.7700000000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=COTI-BTC seq=416058880 bid=0.0000001876 ask=0.0000001904
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=CRO-USDT seq=1388288857 bid=0.0689000000 ask=0.0689200000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DASH-USDT seq=1550465059 bid=42.0100000000 ask=42.0400000000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DOGE-USDT seq=11682666919 bid=0.0915300000 ask=0.0915400000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap progress 50/126 last=DOGE-USDT
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=COTI-USDT seq=1180254072 bid=0.0135000000 ask=0.0135200000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DOT-BTC seq=1709675701 bid=0.0000171100 ask=0.0000171400
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DAG-ETH seq=712436185 bid=0.0000039900 ask=0.0000042100
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DOGE-KCS seq=848142523 bid=0.0108820000 ask=0.0108990000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DOT-KCS seq=909536837 bid=0.1451000000 ask=0.1473000000
2026/04/12 09:12:16 [KuCoin WS 0] bootstrap complete symbol=DOT-USDT seq=8373350839 bid=1.2286000000 ask=1.2287000000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=EGLD-BTC seq=307767890 bid=0.0000523600 ask=0.0000532600
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ENJ-ETH seq=317092693 bid=0.0000141600 ask=0.0000145000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=EGLD-USDT seq=1349612931 bid=3.7600000000 ask=3.7700000000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ERG-USDT seq=1277115969 bid=0.3058000000 ask=0.3069000000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap progress 60/126 last=ERG-USDT
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ETC-BTC seq=1244690002 bid=0.0001147000 ask=0.0001153000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ETC-ETH seq=1354805632 bid=0.0037110000 ask=0.0037190000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ENJ-USDT seq=892113672 bid=0.0315200000 ask=0.0316000000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ETH-EUR seq=3228756860 bid=1895.2800000000 ask=1903.3400000000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ETH-USDT seq=20714266822 bid=2219.4200000000 ask=2219.4300000000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ETC-USDT seq=3784718993 bid=8.2443000000 ask=8.2546000000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=ETH-BTC seq=4759388318 bid=0.0309400000 ask=0.0309500000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=EWT-USDT seq=726889485 bid=0.4549000000 ask=0.4556300000
2026/04/12 09:12:17 [KuCoin WS 0] bootstrap complete symbol=FET-BTC seq=633672145 bid=0.0000033000 ask=0.0000033020
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=FET-ETH seq=755137610 bid=0.0001067300 ask=0.0001071700
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap progress 70/126 last=FET-ETH
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=DOGE-BTC seq=1555219880 bid=0.0000012750 ask=0.0000012770
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=ERG-BTC seq=191240167 bid=0.0000042400 ask=0.0000042600
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=HBAR-BTC seq=915675313 bid=0.0000012050 ask=0.0000012070
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=HBAR-USDT seq=5311722131 bid=0.0865400000 ask=0.0865500000
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=ETH-BRL seq=409345392 bid=11016.1900000000 ask=11237.3900000000
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=HYPE-KCS seq=391393576 bid=4.8650000000 ask=4.8790000000
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=ICP-USDT seq=3333536029 bid=2.4670000000 ask=2.4690000000
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=INJ-BTC seq=792699697 bid=0.0000405200 ask=0.0000409500
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=HYPE-USDT seq=4416880052 bid=40.9300000000 ask=40.9340000000
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap complete symbol=INJ-USDT seq=4688072132 bid=2.9220000000 ask=2.9230000000
2026/04/12 09:12:18 [KuCoin WS 0] bootstrap progress 80/126 last=INJ-USDT
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=IOST-ETH seq=379313490 bid=0.0000004740 ask=0.0000004770
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=ICX-USDT seq=527126385 bid=0.0354800000 ask=0.0355000000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=ICX-ETH seq=127689751 bid=0.0000154300 ask=0.0001500000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=IOTA-BTC seq=262784237 bid=0.0000007780 ask=0.0000007890
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=FET-USDT seq=3632565365 bid=0.2372000000 ask=0.2373000000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=IOTA-USDT seq=916244231 bid=0.0560000000 ask=0.0562000000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=IOTX-ETH seq=216157084 bid=0.0000020710 ask=0.0000021360
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=IOTX-USDT seq=774734539 bid=0.0045900000 ask=0.0046100000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=EWT-BTC seq=322796159 bid=0.0000063530 ask=0.0000063720
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=KAS-BTC seq=308907278 bid=0.0000004540 ask=0.0000004560
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap progress 90/126 last=KAS-BTC
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=ICP-BTC seq=549941352 bid=0.0000343200 ask=0.0000344800
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=KAS-USDT seq=2788984275 bid=0.0326400000 ask=0.0326500000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=KCS-BTC seq=3180867138 bid=0.0001170000 ask=0.0001172000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=KCS-USDT seq=2805341901 bid=8.4070000000 ask=8.4080000000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=KLV-BTC seq=532029054 bid=0.0000000141 ask=0.0000000142
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=KLV-TRX seq=461424472 bid=0.0031500000 ask=0.0031700000
2026/04/12 09:12:19 [KuCoin WS 0] bootstrap complete symbol=KLV-USDT seq=1239217253 bid=0.0010130000 ask=0.0010140000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=KNC-BTC seq=218611267 bid=0.0000018500 ask=0.0000018700
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=KCS-ETH seq=1587649629 bid=0.0037790000 ask=0.0037890000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=IOST-USDT seq=1004504260 bid=0.0010540000 ask=0.0010580000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap progress 100/126 last=IOST-USDT
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=KNC-ETH seq=272163215 bid=0.0000598000 ask=0.0000600000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=KNC-USDT seq=1125313092 bid=0.1327000000 ask=0.1330000000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=KRL-BTC seq=28940378 bid=0.0000020700 ask=0.0000020900
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=KRL-USDT seq=970005659 bid=0.1494600000 ask=0.1496800000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=IOTX-BTC seq=149762993 bid=0.0000000628 ask=0.0000000645
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=LINK-BTC seq=2528149423 bid=0.0001227600 ask=0.0001229000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=LTC-BTC seq=1656815904 bid=0.0007520000 ask=0.0007530000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=LTC-KCS seq=875790482 bid=6.4170000000 ask=6.4360000000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=LTC-ETH seq=1640128383 bid=0.0243200000 ask=0.0243400000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=LYX-ETH seq=608140169 bid=0.0001096000 ask=0.0001097000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap progress 110/126 last=LYX-ETH
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=LYX-USDT seq=1658793644 bid=0.2433000000 ask=0.2438000000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=MANA-ETH seq=645718899 bid=0.0000397000 ask=0.0000399000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=LINK-USDT seq=13123239257 bid=8.8113000000 ask=8.8114000000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=MANTRA-USDT seq=6497153 bid=0.0104600000 ask=0.0104800000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=MANTRA-BTC seq=1738034 bid=0.0000001459 ask=0.0000001463
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=LTC-USDT seq=8540935274 bid=54.0000000000 ask=54.0100000000
2026/04/12 09:12:20 [KuCoin WS 0] bootstrap complete symbol=NEAR-BTC seq=1349433209 bid=0.0000187100 ask=0.0000187400
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap complete symbol=NEO-BTC seq=662028392 bid=0.0000389000 ask=0.0000391000
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap complete symbol=MANA-USDT seq=3223134507 bid=0.0881900000 ask=0.0883400000
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap complete symbol=NFT-TRX seq=379423465 bid=0.0000010340 ask=0.0000010400
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap progress 120/126 last=NFT-TRX
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap complete symbol=NEAR-USDT seq=6887139863 bid=1.3431000000 ask=1.3432000000
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap complete symbol=OGN-USDT seq=797214890 bid=0.0212300000 ask=0.0212900000
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap complete symbol=ONE-BTC seq=498885231 bid=0.0000000284 ask=0.0000000290
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap complete symbol=NEO-USDT seq=2154969697 bid=2.7943000000 ask=2.7983000000
2026/04/12 09:12:21 [KuCoin WS 0] bootstrap complete symbol=NFT-USDT seq=1155261639 bid=0.0000003325 ask=0.0000003331
2026/04/12 09:12:22 [KuCoin WS 0] bootstrap complete symbol=OGN-BTC seq=191483678 bid=0.0000002950 ask=0.0000003080
2026/04/12 09:12:22 [KuCoin WS 0] bootstrap progress 126/126 last=OGN-BTC
2026/04/12 09:12:22 [KuCoin WS 0] bootstrap finished 126/126 in 10.178652091s

