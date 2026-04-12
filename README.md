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



2026/04/12 08:09:30 [KuCoin] started with 2 WS
2026/04/12 08:09:30 [Main] KuCoinCollector started
2026/04/12 08:09:30 [Calculator] indexed 206 symbols
2026/04/12 08:09:32 [KuCoin WS 1] connected
2026/04/12 08:09:32 [KuCoin WS 0] connected
2026/04/12 08:09:42 [KuCoin WS 1] bootstrap complete symbol=ONE-USDT seq=2149422925 bid=0.0020510000 ask=0.0020620000
2026/04/12 08:09:42 [KuCoin WS 1] bootstrap complete symbol=ONT-BTC seq=599176359 bid=0.0000012260 ask=0.0000012460
2026/04/12 08:09:42 [KuCoin WS 1] bootstrap complete symbol=ONT-ETH seq=618942460 bid=0.0000396200 ask=0.0001249600
2026/04/12 08:09:43 [KuCoin WS 1] bootstrap complete symbol=ONT-USDT seq=1628798154 bid=0.0877500000 ask=0.0880900000
2026/04/12 08:09:43 [KuCoin WS 1] bootstrap complete symbol=PAXG-BTC seq=361693979 bid=0.0658140000 ask=0.0658730000
2026/04/12 08:09:44 [KuCoin WS 1] bootstrap complete symbol=PAXG-USDT seq=1638386569 bid=4712.8600000000 ask=4712.8700000000
2026/04/12 08:09:44 [KuCoin WS 1] bootstrap complete symbol=PEPE-KCS seq=146625026 bid=0.0000004155 ask=0.0000004191
2026/04/12 08:09:46 [KuCoin WS 1] bootstrap complete symbol=PEPE-USDT seq=5048072677 bid=0.0000034980 ask=0.0000034990
2026/04/12 08:09:46 [KuCoin WS 1] bootstrap complete symbol=POND-BTC seq=123770210 bid=0.0000000303 ask=0.0000000312
2026/04/12 08:09:46 [KuCoin WS 1] bootstrap complete symbol=POND-USDT seq=434238043 bid=0.0021890000 ask=0.0022020000
2026/04/12 08:09:47 [KuCoin WS 1] bootstrap complete symbol=RLC-BTC seq=144017974 bid=0.0000058000 ask=0.0000058200
2026/04/12 08:09:47 [KuCoin WS 0] bootstrap complete symbol=A-BTC seq=168735691 bid=0.0000010830 ask=0.0000010940
2026/04/12 08:09:48 [KuCoin WS 0] bootstrap complete symbol=A-ETH seq=583843485 bid=0.0000350000 ask=0.0000442000
2026/04/12 08:09:48 [KuCoin WS 1] bootstrap complete symbol=RLC-USDT seq=1052830748 bid=0.4153000000 ask=0.4159000000
2026/04/12 08:09:48 [KuCoin WS 0] bootstrap complete symbol=A-USDT seq=143893478 bid=0.0776000000 ask=0.0779000000
2026/04/12 08:09:48 [KuCoin WS 1] bootstrap complete symbol=RSR-BTC seq=132293915 bid=0.0000000205 ask=0.0000000206
2026/04/12 08:09:48 [KuCoin WS 0] bootstrap complete symbol=AAVE-BTC seq=1081988584 bid=0.0012530000 ask=0.0012560000
2026/04/12 08:09:48 [KuCoin WS 1] bootstrap complete symbol=RSR-USDT seq=1855608055 bid=0.0014740000 ask=0.0014760000
2026/04/12 08:09:49 [KuCoin WS 0] bootstrap complete symbol=AAVE-USDT seq=9002348390 bid=89.7510000000 ask=89.7540000000
2026/04/12 08:09:49 [KuCoin WS 1] bootstrap complete symbol=RUNE-BTC seq=1114754169 bid=0.0000054500 ask=0.0000054700
2026/04/12 08:09:49 [KuCoin WS 0] bootstrap complete symbol=ADA-BTC seq=1105254527 bid=0.0000033900 ask=0.0000034100
2026/04/12 08:09:49 [KuCoin WS 1] bootstrap complete symbol=RUNE-USDT seq=3496756186 bid=0.3906000000 ask=0.3911000000
2026/04/12 08:09:49 [KuCoin WS 0] bootstrap complete symbol=ADA-KCS seq=527455063 bid=0.0289800000 ask=0.0291400000
2026/04/12 08:09:49 [KuCoin WS 1] bootstrap complete symbol=SCRT-BTC seq=69768350 bid=0.0000012690 ask=0.0000012710
2026/04/12 08:09:49 [KuCoin WS 0] bootstrap complete symbol=ADA-USDT seq=8277392844 bid=0.2432000000 ask=0.2433000000
2026/04/12 08:09:50 [KuCoin WS 1] bootstrap complete symbol=SCRT-USDT seq=341458606 bid=0.0907000000 ask=0.0908000000
2026/04/12 08:09:50 [KuCoin WS 0] bootstrap complete symbol=ALGO-BTC seq=749823720 bid=0.0000014730 ask=0.0000014890
2026/04/12 08:09:50 [KuCoin WS 1] bootstrap complete symbol=SHIB-DOGE seq=676597972 bid=0.0000637800 ask=0.0000639500
2026/04/12 08:09:50 [KuCoin WS 0] bootstrap complete symbol=ALGO-ETH seq=943896478 bid=0.0000478000 ask=0.0000479500
2026/04/12 08:09:50 [KuCoin WS 0] bootstrap complete symbol=ALGO-USDT seq=2602524624 bid=0.1058000000 ask=0.1059000000
2026/04/12 08:09:51 [KuCoin WS 0] bootstrap complete symbol=ANKR-BTC seq=338220871 bid=0.0000000730 ask=0.0000000739
2026/04/12 08:09:51 [KuCoin WS 0] bootstrap complete symbol=ANKR-USDT seq=800932323 bid=0.0052300000 ask=0.0052700000
2026/04/12 08:09:51 [KuCoin WS 0] bootstrap complete symbol=AR-BTC seq=147346631 bid=0.0000236000 ask=0.0000238000
2026/04/12 08:09:52 [KuCoin WS 1] bootstrap complete symbol=SHIB-USDT seq=5998473062 bid=0.0000058240 ask=0.0000058250
2026/04/12 08:09:52 [KuCoin WS 0] bootstrap complete symbol=AR-USDT seq=2386082367 bid=1.6960000000 ask=1.6980000000
2026/04/12 08:09:52 [KuCoin WS 1] bootstrap complete symbol=SNX-BTC seq=424262557 bid=0.0000039500 ask=0.0000041100
2026/04/12 08:09:52 [KuCoin WS 0] bootstrap complete symbol=ATOM-BTC seq=1755280149 bid=0.0000244100 ask=0.0000244200
2026/04/12 08:09:52 [KuCoin WS 0] bootstrap complete symbol=ATOM-ETH seq=1564345163 bid=0.0007880000 ask=0.0007900000
2026/04/12 08:09:53 [KuCoin WS 0] bootstrap complete symbol=ATOM-USDT seq=6314399893 bid=1.7455000000 ask=1.7456000000
2026/04/12 08:09:53 [KuCoin WS 0] bootstrap complete symbol=AVA-BTC seq=340825831 bid=0.0000027900 ask=0.0000028200
2026/04/12 08:09:53 [KuCoin WS 1] bootstrap complete symbol=SNX-ETH seq=505517469 bid=0.0001280000 ask=0.0001340000
2026/04/12 08:09:53 [KuCoin WS 0] bootstrap complete symbol=AVA-ETH seq=427916945 bid=0.0000901000 ask=0.0000915000
2026/04/12 08:09:53 [KuCoin WS 1] bootstrap complete symbol=SNX-USDT seq=1536416821 bid=0.2838000000 ask=0.2850000000
2026/04/12 08:09:53 [KuCoin WS 0] bootstrap complete symbol=AVA-USDT seq=1112963797 bid=0.1998000000 ask=0.2004000000
2026/04/12 08:09:54 [KuCoin WS 1] bootstrap complete symbol=SOL-KCS seq=458269564 bid=9.8060000000 ask=9.8340000000
2026/04/12 08:09:54 [KuCoin WS 0] bootstrap complete symbol=AVAX-BTC seq=1689066360 bid=0.0001266900 ask=0.0001268700
2026/04/12 08:09:54 [KuCoin WS 0] bootstrap complete symbol=AVAX-USDT seq=9252887794 bid=9.0720000000 ask=9.0730000000
2026/04/12 08:09:54 [KuCoin WS 0] bootstrap complete symbol=BCH-BTC seq=1261424704 bid=0.0059380000 ask=0.0059440000
2026/04/12 08:09:55 [KuCoin WS 1] bootstrap complete symbol=SOL-USDT seq=19083046859 bid=82.3100000000 ask=82.3200000000
2026/04/12 08:09:55 [KuCoin WS 0] bootstrap complete symbol=BCH-USDT seq=4654145409 bid=425.1600000000 ask=425.1700000000
2026/04/12 08:09:55 [KuCoin WS 1] bootstrap complete symbol=STORJ-ETH seq=100782486 bid=0.0000440000 ask=0.0000444000
2026/04/12 08:09:55 [KuCoin WS 0] bootstrap complete symbol=BCHSV-BTC seq=974979365 bid=0.0002174000 ask=0.0002185000
2026/04/12 08:09:55 [KuCoin WS 0] bootstrap complete symbol=BCHSV-ETH seq=1101757551 bid=0.0068800000 ask=0.0073500000
2026/04/12 08:09:56 [KuCoin WS 1] bootstrap complete symbol=STORJ-USDT seq=901049393 bid=0.0976000000 ask=0.0979000000
2026/04/12 08:09:56 [KuCoin WS 1] bootstrap complete symbol=STX-BTC seq=442992272 bid=0.0000029900 ask=0.0000030000
2026/04/12 08:09:57 [KuCoin WS 1] bootstrap complete symbol=STX-USDT seq=3150657750 bid=0.2143000000 ask=0.2144000000
2026/04/12 08:09:57 [KuCoin WS 0] bootstrap complete symbol=BCHSV-USDT seq=1534646353 bid=15.5700000000 ask=15.5900000000
2026/04/12 08:09:57 [KuCoin WS 1] bootstrap complete symbol=SUI-KCS seq=243383481 bid=0.1084000000 ask=0.1089000000
2026/04/12 08:09:57 [KuCoin WS 0] bootstrap complete symbol=BDX-BTC seq=331350670 bid=0.0000011160 ask=0.0000011190
2026/04/12 08:09:57 [KuCoin WS 0] bootstrap complete symbol=BDX-USDT seq=1920839204 bid=0.0799500000 ask=0.0799600000
2026/04/12 08:09:58 [KuCoin WS 0] bootstrap complete symbol=BNB-BTC seq=2283254377 bid=0.0083053000 ask=0.0083152000
2026/04/12 08:09:58 [KuCoin WS 0] bootstrap complete symbol=BNB-KCS seq=1205055932 bid=70.8264000000 ask=70.9810000000
2026/04/12 08:09:58 [KuCoin WS 1] bootstrap complete symbol=SUI-USDT seq=7715448234 bid=0.9107000000 ask=0.9108000000
2026/04/12 08:09:59 [KuCoin WS 0] bootstrap complete symbol=BNB-USDT seq=11268195981 bid=594.6950000000 ask=594.6960000000
2026/04/12 08:09:59 [KuCoin WS 1] bootstrap complete symbol=TEL-BTC seq=2901147268 bid=0.0000000293 ask=0.0000000297
2026/04/12 08:09:59 [KuCoin WS 0] bootstrap complete symbol=BTC-BRL seq=415483090 bid=348970.4000000000 ask=365295.2000000000
2026/04/12 08:09:59 [KuCoin WS 1] bootstrap complete symbol=TEL-ETH seq=1361282019 bid=0.0000009484 ask=0.0000009687
2026/04/12 08:09:59 [KuCoin WS 0] bootstrap complete symbol=BTC-EUR seq=2880869581 bid=61105.0900000000 ask=61324.2100000000
2026/04/12 08:10:00 [KuCoin WS 0] bootstrap complete symbol=BTC-USDT seq=31659769814 bid=71574.0000000000 ask=71574.1000000000
2026/04/12 08:10:00 [KuCoin WS 0] bootstrap complete symbol=CHZ-BTC seq=500794570 bid=0.0000005302 ask=0.0000005354
2026/04/12 08:10:00 [KuCoin WS 0] bootstrap complete symbol=CHZ-USDT seq=1957174672 bid=0.0381200000 ask=0.0381300000
2026/04/12 08:10:00 [KuCoin WS 1] bootstrap complete symbol=TEL-USDT seq=5957354853 bid=0.0021030000 ask=0.0021080000
2026/04/12 08:10:00 [KuCoin WS 0] bootstrap complete symbol=CKB-BTC seq=103428799 bid=0.0000000200 ask=0.0000000206
2026/04/12 08:10:00 [KuCoin WS 1] bootstrap complete symbol=TRAC-BTC seq=324992919 bid=0.0000040300 ask=0.0000041000
2026/04/12 08:10:01 [KuCoin WS 0] bootstrap complete symbol=CKB-USDT seq=1270016143 bid=0.0014560000 ask=0.0014590000
2026/04/12 08:10:01 [KuCoin WS 1] bootstrap complete symbol=TRAC-ETH seq=398220807 bid=0.0001292000 ask=0.0001365000
2026/04/12 08:10:01 [KuCoin WS 0] bootstrap complete symbol=COTI-BTC seq=416056951 bid=0.0000001876 ask=0.0000001901
2026/04/12 08:10:01 [KuCoin WS 1] bootstrap complete symbol=TRAC-USDT seq=1072024752 bid=0.2892000000 ask=0.2905000000
2026/04/12 08:10:01 [KuCoin WS 0] bootstrap complete symbol=COTI-USDT seq=1180249151 bid=0.0134900000 ask=0.0135100000
2026/04/12 08:10:01 [KuCoin WS 1] bootstrap complete symbol=TRX-BTC seq=977985994 bid=0.0000044770 ask=0.0000044800
2026/04/12 08:10:02 [KuCoin WS 0] bootstrap complete symbol=CRO-BTC seq=10919267363 bid=0.0000009590 ask=0.0000009680
2026/04/12 08:10:02 [KuCoin WS 1] bootstrap complete symbol=TRX-ETH seq=1200022285 bid=0.0001446900 ask=0.0001448400
2026/04/12 08:10:02 [KuCoin WS 0] bootstrap complete symbol=CRO-USDT seq=1388215543 bid=0.0687200000 ask=0.0687300000
2026/04/12 08:10:02 [KuCoin WS 1] bootstrap complete symbol=TRX-USDT seq=1869291963 bid=0.3204000000 ask=0.3205000000
2026/04/12 08:10:02 [KuCoin WS 0] bootstrap complete symbol=CSPR-ETH seq=260331838 bid=0.0000013400 ask=0.0000013800
2026/04/12 08:10:02 [KuCoin WS 1] bootstrap complete symbol=TWT-BTC seq=146486381 bid=0.0000056800 ask=0.0000058100
2026/04/12 08:10:03 [KuCoin WS 0] bootstrap complete symbol=CSPR-USDT seq=624717873 bid=0.0029820000 ask=0.0030010000
2026/04/12 08:10:03 [KuCoin WS 1] bootstrap complete symbol=TWT-USDT seq=909082440 bid=0.4125000000 ask=0.4130000000
2026/04/12 08:10:03 [KuCoin WS 0] bootstrap complete symbol=DAG-ETH seq=712403043 bid=0.0000040000 ask=0.0000042000
2026/04/12 08:10:03 [KuCoin WS 1] bootstrap complete symbol=USDT-BRL seq=1138949446 bid=5.0393000000 ask=5.0496000000
2026/04/12 08:10:03 [KuCoin WS 0] bootstrap complete symbol=DAG-USDT seq=3818748312 bid=0.0089330000 ask=0.0089360000
2026/04/12 08:10:03 [KuCoin WS 1] bootstrap complete symbol=USDT-EUR seq=457755354 bid=0.8550000000 ask=0.8552000000
2026/04/12 08:10:04 [KuCoin WS 0] bootstrap complete symbol=DASH-BTC seq=601235439 bid=0.0005860000 ask=0.0005871000
2026/04/12 08:10:04 [KuCoin WS 1] bootstrap complete symbol=VET-BTC seq=599276726 bid=0.0000000955 ask=0.0000000962
2026/04/12 08:10:04 [KuCoin WS 0] bootstrap complete symbol=DASH-ETH seq=602700868 bid=0.0189500000 ask=0.0189900000
2026/04/12 08:10:04 [KuCoin WS 1] bootstrap complete symbol=VET-ETH seq=571318623 bid=0.0000030900 ask=0.0000031000
2026/04/12 08:10:04 [KuCoin WS 0] bootstrap complete symbol=DASH-USDT seq=1550380970 bid=41.9600000000 ask=41.9900000000
2026/04/12 08:10:04 [KuCoin WS 1] bootstrap complete symbol=VET-USDT seq=2105048898 bid=0.0068500000 ask=0.0068600000
2026/04/12 08:10:05 [KuCoin WS 0] bootstrap complete symbol=DOGE-BTC seq=1555214929 bid=0.0000012720 ask=0.0000012750
2026/04/12 08:10:05 [KuCoin WS 1] bootstrap complete symbol=VSYS-BTC seq=160751972 bid=0.0000000032 ask=0.0000000032
2026/04/12 08:10:05 [KuCoin WS 0] bootstrap complete symbol=DOGE-KCS seq=848114549 bid=0.0108680000 ask=0.0108730000
2026/04/12 08:10:05 [KuCoin WS 1] bootstrap complete symbol=VSYS-USDT seq=806191770 bid=0.0002299000 ask=0.0002300000
2026/04/12 08:10:05 [KuCoin WS 0] bootstrap complete symbol=DOGE-USDT seq=11682472519 bid=0.0911600000 ask=0.0911700000
2026/04/12 08:10:05 [KuCoin WS 0] bootstrap complete symbol=DOT-BTC seq=1709658611 bid=0.0000171700 ask=0.0000171900
2026/04/12 08:10:06 [KuCoin WS 1] bootstrap complete symbol=WAN-BTC seq=272947127 bid=0.0000007750 ask=0.0000007850
2026/04/12 08:10:06 [KuCoin WS 0] bootstrap complete symbol=DOT-KCS seq=909526176 bid=0.1464000000 ask=0.1471000000
2026/04/12 08:10:06 [KuCoin WS 1] bootstrap complete symbol=WAN-USDT seq=340564067 bid=0.0553300000 ask=0.0557600000
2026/04/12 08:10:06 [KuCoin WS 0] bootstrap complete symbol=DOT-USDT seq=8373285049 bid=1.2294000000 ask=1.2295000000
2026/04/12 08:10:06 [KuCoin WS 0] bootstrap complete symbol=EGLD-BTC seq=307765890 bid=0.0000523600 ask=0.0000528000
2026/04/12 08:10:07 [KuCoin WS 0] bootstrap complete symbol=EGLD-USDT seq=1349609991 bid=3.7500000000 ask=3.7600000000
2026/04/12 08:10:07 [KuCoin WS 1] bootstrap complete symbol=WAVES-BTC seq=465412063 bid=0.0000057000 ask=0.0000057600
2026/04/12 08:10:07 [KuCoin WS 0] bootstrap complete symbol=ENJ-ETH seq=317080933 bid=0.0000142600 ask=0.0000147400
2026/04/12 08:10:07 [KuCoin WS 1] bootstrap complete symbol=WAVES-USDT seq=2035030251 bid=0.4096000000 ask=0.4102000000
2026/04/12 08:10:07 [KuCoin WS 0] bootstrap complete symbol=ENJ-USDT seq=891994162 bid=0.0316700000 ask=0.0317400000
2026/04/12 08:10:08 [KuCoin WS 0] bootstrap complete symbol=ERG-BTC seq=191237900 bid=0.0000042400 ask=0.0000042600
2026/04/12 08:10:08 [KuCoin WS 1] bootstrap complete symbol=WBTC-BTC seq=705047982 bid=0.9957500000 ask=0.9993400000
2026/04/12 08:10:08 [KuCoin WS 0] bootstrap complete symbol=ERG-USDT seq=1277112692 bid=0.3068000000 ask=0.3072000000
2026/04/12 08:10:08 [KuCoin WS 1] bootstrap complete symbol=WBTC-USDT seq=219797509 bid=71234.1600000000 ask=71486.8400000000
2026/04/12 08:10:08 [KuCoin WS 0] bootstrap complete symbol=ETC-BTC seq=1244670206 bid=0.0001148000 ask=0.0001152000
2026/04/12 08:10:08 [KuCoin WS 1] bootstrap complete symbol=WIN-BTC seq=278140691 bid=0.0000000003 ask=0.0000000003
2026/04/12 08:10:09 [KuCoin WS 0] bootstrap complete symbol=ETC-ETH seq=1354778767 bid=0.0037120000 ask=0.0037240000
2026/04/12 08:10:09 [KuCoin WS 1] bootstrap complete symbol=WIN-TRX seq=291940910 bid=0.0000591000 ask=0.0000597000
2026/04/12 08:10:09 [KuCoin WS 0] bootstrap complete symbol=ETC-USDT seq=3784679452 bid=8.2258000000 ask=8.2401000000
2026/04/12 08:10:09 [KuCoin WS 1] bootstrap complete symbol=WIN-USDT seq=937717014 bid=0.0000189700 ask=0.0000190300
2026/04/12 08:10:09 [KuCoin WS 0] bootstrap complete symbol=ETH-BRL seq=409340241 bid=10986.0600000000 ask=11216.1200000000
2026/04/12 08:10:09 [KuCoin WS 1] bootstrap complete symbol=XDC-BTC seq=546100470 bid=0.0000004210 ask=0.0000004270
2026/04/12 08:10:10 [KuCoin WS 1] bootstrap complete symbol=XDC-ETH seq=651341634 bid=0.0000136600 ask=0.0000137500
2026/04/12 08:10:10 [KuCoin WS 1] bootstrap complete symbol=XDC-USDT seq=1533449847 bid=0.0303200000 ask=0.0303500000
2026/04/12 08:10:10 [KuCoin WS 0] bootstrap complete symbol=ETH-BTC seq=4759309347 bid=0.0309200000 ask=0.0309300000
2026/04/12 08:10:10 [KuCoin WS 1] bootstrap complete symbol=XLM-BTC seq=871804294 bid=0.0000021140 ask=0.0000021180
2026/04/12 08:10:11 [KuCoin WS 0] bootstrap complete symbol=ETH-EUR seq=3228550654 bid=1890.3800000000 ask=1898.3200000000
2026/04/12 08:10:11 [KuCoin WS 1] bootstrap complete symbol=XLM-ETH seq=843161020 bid=0.0000683800 ask=0.0000685500
2026/04/12 08:10:11 [KuCoin WS 0] bootstrap complete symbol=ETH-USDT seq=20713881520 bid=2212.6900000000 ask=2212.7000000000
2026/04/12 08:10:11 [KuCoin WS 1] bootstrap complete symbol=XLM-USDT seq=2521847599 bid=0.1514000000 ask=0.1515000000
2026/04/12 08:10:11 [KuCoin WS 0] bootstrap complete symbol=EWT-BTC seq=322771728 bid=0.0000062650 ask=0.0000062870
2026/04/12 08:10:11 [KuCoin WS 1] bootstrap complete symbol=XMR-BTC seq=2062792734 bid=0.0047440000 ask=0.0047490000
2026/04/12 08:10:11 [KuCoin WS 0] bootstrap complete symbol=EWT-USDT seq=726846454 bid=0.4485700000 ask=0.4493300000
2026/04/12 08:10:12 [KuCoin WS 1] bootstrap complete symbol=XMR-ETH seq=2730961526 bid=0.1534300000 ask=0.1536300000
2026/04/12 08:10:12 [KuCoin WS 0] bootstrap complete symbol=FET-BTC seq=633665635 bid=0.0000032720 ask=0.0000032790
2026/04/12 08:10:12 [KuCoin WS 1] bootstrap complete symbol=XMR-USDT seq=13313125361 bid=339.6600000000 ask=339.7600000000
2026/04/12 08:10:12 [KuCoin WS 1] bootstrap complete symbol=XRP-BTC seq=1480155791 bid=0.0000185800 ask=0.0000185900
2026/04/12 08:10:12 [KuCoin WS 0] bootstrap complete symbol=FET-ETH seq=755125656 bid=0.0001058200 ask=0.0001065100
2026/04/12 08:10:13 [KuCoin WS 0] bootstrap complete symbol=FET-USDT seq=3632534439 bid=0.2345000000 ask=0.2346000000
2026/04/12 08:10:13 [KuCoin WS 0] bootstrap complete symbol=HBAR-BTC seq=915666992 bid=0.0000012080 ask=0.0000012100
2026/04/12 08:10:13 [KuCoin WS 0] bootstrap complete symbol=HBAR-USDT seq=5311683454 bid=0.0865400000 ask=0.0865500000
2026/04/12 08:10:14 [KuCoin WS 0] bootstrap complete symbol=HYPE-KCS seq=391375341 bid=4.8650000000 ask=4.8750000000
2026/04/12 08:10:14 [KuCoin WS 0] bootstrap complete symbol=HYPE-USDT seq=4416742961 bid=40.8470000000 ask=40.8480000000
2026/04/12 08:10:14 [KuCoin WS 1] bootstrap complete symbol=XRP-ETH seq=1789478373 bid=0.0006004000 ask=0.0006007000
2026/04/12 08:10:14 [KuCoin WS 0] bootstrap complete symbol=ICP-BTC seq=549937343 bid=0.0000343000 ask=0.0000343500
2026/04/12 08:10:15 [KuCoin WS 0] bootstrap complete symbol=ICP-USDT seq=3333527215 bid=2.4560000000 ask=2.4570000000
2026/04/12 08:10:15 [KuCoin WS 0] bootstrap complete symbol=ICX-ETH seq=127688067 bid=0.0000156400 ask=0.0000185500
2026/04/12 08:10:15 [KuCoin WS 0] bootstrap complete symbol=ICX-USDT seq=527117550 bid=0.0353700000 ask=0.0354100000
2026/04/12 08:10:16 [KuCoin WS 0] bootstrap complete symbol=INJ-BTC seq=792696436 bid=0.0000405900 ask=0.0000410500
2026/04/12 08:10:16 [KuCoin WS 0] bootstrap complete symbol=INJ-USDT seq=4688065377 bid=2.9200000000 ask=2.9220000000
2026/04/12 08:10:16 [KuCoin WS 1] bootstrap complete symbol=XRP-KCS seq=1001129553 bid=0.1584100000 ask=0.1586900000
2026/04/12 08:10:16 [KuCoin WS 0] bootstrap complete symbol=IOST-ETH seq=379311939 bid=0.0000004750 ask=0.0000004770
2026/04/12 08:10:16 [KuCoin WS 1] bootstrap complete symbol=XRP-USDT seq=18978657191 bid=1.3286200000 ask=1.3286300000
2026/04/12 08:10:17 [KuCoin WS 0] bootstrap complete symbol=IOST-USDT seq=1004501665 bid=0.0010500000 ask=0.0010540000
2026/04/12 08:10:17 [KuCoin WS 1] bootstrap complete symbol=XTZ-BTC seq=435763461 bid=0.0000048100 ask=0.0000048400
2026/04/12 08:10:17 [KuCoin WS 0] bootstrap complete symbol=IOTA-BTC seq=262782736 bid=0.0000007820 ask=0.0000007880
2026/04/12 08:10:17 [KuCoin WS 1] bootstrap complete symbol=XTZ-USDT seq=1249571598 bid=0.3451000000 ask=0.3456000000
2026/04/12 08:10:17 [KuCoin WS 0] bootstrap complete symbol=IOTA-USDT seq=916238857 bid=0.0560000000 ask=0.0562000000
2026/04/12 08:10:17 [KuCoin WS 1] bootstrap complete symbol=XYO-BTC seq=64517814 bid=0.0000000498 ask=0.0000000503
2026/04/12 08:10:17 [KuCoin WS 0] bootstrap complete symbol=IOTX-BTC seq=149761021 bid=0.0000000633 ask=0.0000000648
2026/04/12 08:10:18 [KuCoin WS 1] bootstrap complete symbol=XYO-ETH seq=140273103 bid=0.0000016110 ask=0.0000016330
2026/04/12 08:10:18 [KuCoin WS 1] bootstrap complete symbol=XYO-USDT seq=2119253703 bid=0.0035760000 ask=0.0035880000
2026/04/12 08:10:18 [KuCoin WS 1] bootstrap complete symbol=ZEC-BTC seq=1266503984 bid=0.0050218000 ask=0.0050272000
2026/04/12 08:10:18 [KuCoin WS 1] bootstrap complete symbol=ZEC-USDT seq=3472172973 bid=359.4820000000 ask=359.6020000000
2026/04/12 08:10:19 [KuCoin WS 0] bootstrap complete symbol=IOTX-ETH seq=216153590 bid=0.0000020600 ask=0.0000020840
2026/04/12 08:10:19 [KuCoin WS 0] bootstrap complete symbol=IOTX-USDT seq=774732002 bid=0.0045800000 ask=0.0046000000
2026/04/12 08:10:20 [KuCoin WS 1] bootstrap complete symbol=ZIL-ETH seq=680983326 bid=0.0000017440 ask=0.0000017510
2026/04/12 08:10:20 [KuCoin WS 0] bootstrap complete symbol=KAS-BTC seq=308903831 bid=0.0000004520 ask=0.0000004540
2026/04/12 08:10:20 [KuCoin WS 1] bootstrap complete symbol=ZIL-USDT seq=1525721370 bid=0.0038640000 ask=0.0038670000
2026/04/12 08:10:20 [KuCoin WS 0] bootstrap complete symbol=KAS-USDT seq=2788972175 bid=0.0324300000 ask=0.0324600000
2026/04/12 08:10:20 [KuCoin WS 0] bootstrap complete symbol=KCS-BTC seq=3180842490 bid=0.0001171000 ask=0.0001172000
2026/04/12 08:10:21 [KuCoin WS 0] bootstrap complete symbol=KCS-ETH seq=1587628032 bid=0.0037850000 ask=0.0037920000
2026/04/12 08:10:21 [KuCoin WS 0] bootstrap complete symbol=KCS-USDT seq=2805295086 bid=8.3790000000 ask=8.3900000000
2026/04/12 08:10:21 [KuCoin WS 0] bootstrap complete symbol=KLV-BTC seq=532026206 bid=0.0000000141 ask=0.0000000142
2026/04/12 08:10:22 [KuCoin WS 0] bootstrap complete symbol=KLV-TRX seq=461421447 bid=0.0031500000 ask=0.0031700000
2026/04/12 08:10:22 [KuCoin WS 0] bootstrap complete symbol=KLV-USDT seq=1239214644 bid=0.0010130000 ask=0.0010140000
2026/04/12 08:10:22 [KuCoin WS 0] bootstrap complete symbol=KNC-BTC seq=218608346 bid=0.0000018440 ask=0.0000018660
2026/04/12 08:10:22 [KuCoin WS 0] bootstrap complete symbol=KNC-ETH seq=272160317 bid=0.0000600000 ask=0.0000602000
2026/04/12 08:10:23 [KuCoin WS 0] bootstrap complete symbol=KNC-USDT seq=1125310751 bid=0.1328000000 ask=0.1330000000
2026/04/12 08:10:23 [KuCoin WS 0] bootstrap complete symbol=KRL-BTC seq=28940059 bid=0.0000020709 ask=0.0000020900
2026/04/12 08:10:23 [KuCoin WS 0] bootstrap complete symbol=KRL-USDT seq=970005158 bid=0.1494500000 ask=0.1496900000
2026/04/12 08:10:24 [KuCoin WS 0] bootstrap complete symbol=LINK-BTC seq=2528122264 bid=0.0001226100 ask=0.0001227700
2026/04/12 08:10:24 [KuCoin WS 0] bootstrap complete symbol=LINK-USDT seq=13123099939 bid=8.7798000000 ask=8.7799000000
2026/04/12 08:10:26 [KuCoin WS 0] bootstrap complete symbol=LTC-BTC seq=1656806686 bid=0.0007540000 ask=0.0007550000
2026/04/12 08:10:27 [KuCoin WS 0] bootstrap complete symbol=LTC-ETH seq=1640113259 bid=0.0243900000 ask=0.0244200000
2026/04/12 08:10:27 [KuCoin WS 0] bootstrap complete symbol=LTC-KCS seq=875775655 bid=6.4290000000 ask=6.4460000000
2026/04/12 08:10:29 [KuCoin WS 0] bootstrap complete symbol=LTC-USDT seq=8540908112 bid=53.9900000000 ask=54.0000000000
2026/04/12 08:10:29 [KuCoin WS 0] bootstrap complete symbol=LYX-ETH seq=608133191 bid=0.0001094000 ask=0.0001111000
2026/04/12 08:10:30 [KuCoin WS 0] bootstrap complete symbol=LYX-USDT seq=1658763277 bid=0.2428000000 ask=0.2454000000
2026/04/12 08:10:30 [KuCoin WS 0] bootstrap complete symbol=MANA-ETH seq=645717588 bid=0.0000395000 ask=0.0000397000
2026/04/12 08:10:30 [KuCoin WS 0] bootstrap complete symbol=MANA-USDT seq=3223128029 bid=0.0876000000 ask=0.0877200000
2026/04/12 08:10:31 [KuCoin WS 0] bootstrap complete symbol=MANTRA-BTC seq=1737063 bid=0.0000001460 ask=0.0000001463
2026/04/12 08:10:31 [KuCoin WS 0] bootstrap complete symbol=MANTRA-USDT seq=6494512 bid=0.0104300000 ask=0.0104600000
2026/04/12 08:10:31 [KuCoin WS 0] bootstrap complete symbol=NEAR-BTC seq=1349420423 bid=0.0000187300 ask=0.0000187600
2026/04/12 08:10:32 [KuCoin WS 0] bootstrap complete symbol=NEAR-USDT seq=6887104141 bid=1.3416000000 ask=1.3418000000
2026/04/12 08:10:32 [KuCoin WS 0] bootstrap complete symbol=NEO-BTC seq=662026891 bid=0.0000389000 ask=0.0000392000
2026/04/12 08:10:33 [KuCoin WS 0] bootstrap complete symbol=NEO-USDT seq=2154960374 bid=2.7917000000 ask=2.7941000000
2026/04/12 08:10:34 [KuCoin WS 0] bootstrap complete symbol=NFT-TRX seq=379418124 bid=0.0000010360 ask=0.0000010430
2026/04/12 08:10:35 [KuCoin WS 0] bootstrap complete symbol=NFT-USDT seq=1155255448 bid=0.0000003328 ask=0.0000003336
2026/04/12 08:10:35 [KuCoin WS 0] bootstrap complete symbol=OGN-BTC seq=191481576 bid=0.0000002950 ask=0.0000002980
2026/04/12 08:10:36 [KuCoin WS 0] bootstrap complete symbol=OGN-USDT seq=797205380 bid=0.0212100000 ask=0.0212800000
2026/04/12 08:10:36 [KuCoin WS 0] bootstrap complete symbol=ONE-BTC seq=498883075 bid=0.0000000285 ask=0.0000000289

