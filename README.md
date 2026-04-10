2026/04/10 08:59:48 [KuCoin] started with 1 WS
2026/04/10 08:59:48 [Main] KuCoinCollector started
2026/04/10 08:59:48 [CALC] started | triangles indexed=56 | minVolume=100.00 USDT | minProfit=0.1000% | quoteAgeMax=400ms | logMode=debug
2026/04/10 08:59:49 [REJECT] stage=scan reason=no_quote_leg_1 count=1 tri=USDT->BTC->HBAR
2026/04/10 08:59:50 [REJECT] stage=scan reason=no_quote_leg_1 count=10 tri=USDT->BTC->ETH
2026/04/10 08:59:50 [REJECT] stage=scan reason=no_quote_leg_2 count=1 tri=USDT->BTC->LTC
2026/04/10 08:59:50 [REJECT] stage=scan reason=no_quote_leg_2 count=10 tri=USDT->BTC->ANKR
2026/04/10 08:59:50 [REJECT] stage=scan reason=no_quote_leg_3 count=1 tri=USDT->BTC->KNC
2026/04/10 08:59:53 [STATS] ticks=38 triangles_seen=87 cand=0 exec=0 | scan_rejects={no_quote_leg_3=1, no_quote_leg_1=43, no_quote_leg_2=43} | exec_rejects={none}
2026/04/10 08:59:56 [REJECT] stage=scan reason=no_quote_leg_1 count=100 tri=USDT->BTC->DASH
2026/04/10 08:59:56 [REJECT] stage=scan reason=no_quote_leg_3 count=10 tri=USDT->LYX->ETH
2026/04/10 08:59:57 [REJECT] stage=scan reason=max_start_lt_100 count=1 tri=USDT->KCS->SOL
2026/04/10 08:59:58 [STATS] ticks=146 triangles_seen=233 cand=0 exec=0 | scan_rejects={max_start_lt_100=1, no_quote_leg_3=19, no_quote_leg_2=69, no_quote_leg_1=144} | exec_rejects={none}
2026/04/10 09:00:00 [EXEC CMP] USDT→DOT→KCS | est=-0.5720% | ideal=-0.2731% | rounded=-0.2738% | strict=-0.5728% | tails=-0.5726% | tail_usdt=0.00017320 | tails_map={USDT=0.00017320}
2026/04/10 09:00:00 [REJECT] stage=exec reason=profit_below_threshold count=1 tri=USDT->DOT->KCS
2026/04/10 09:00:00 [EXEC CMP] USDT→DOT→KCS | est=-0.5602% | ideal=-0.2613% | rounded=-0.2619% | strict=-0.5610% | tails=-0.5608% | tail_usdt=0.00016140 | tails_map={USDT=0.00016140}
2026/04/10 09:00:01 [REJECT] stage=scan reason=no_quote_leg_2 count=100 tri=USDT->KCS->SOL
2026/04/10 09:00:03 [STATS] ticks=381 triangles_seen=538 cand=2 exec=0 | scan_rejects={max_start_lt_100=8, no_quote_leg_3=73, no_quote_leg_2=115, no_quote_leg_1=340} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:08 [STATS] ticks=517 triangles_seen=696 cand=2 exec=0 | scan_rejects={max_start_lt_100=8, no_quote_leg_3=96, no_quote_leg_2=144, no_quote_leg_1=446} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:09 [REJECT] stage=scan reason=no_quote_leg_3 count=100 tri=USDT->KCS->SOL
2026/04/10 09:00:10 [REJECT] stage=scan reason=max_start_lt_100 count=10 tri=USDT->DOT->KCS
2026/04/10 09:00:13 [STATS] ticks=649 triangles_seen=884 cand=2 exec=0 | scan_rejects={max_start_lt_100=21, no_quote_leg_3=118, no_quote_leg_2=175, no_quote_leg_1=568} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:18 [STATS] ticks=820 triangles_seen=1112 cand=2 exec=0 | scan_rejects={max_start_lt_100=25, no_quote_leg_3=173, no_quote_leg_2=194, no_quote_leg_1=718} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:23 [STATS] ticks=923 triangles_seen=1248 cand=2 exec=0 | scan_rejects={max_start_lt_100=25, no_quote_leg_2=202, no_quote_leg_3=208, no_quote_leg_1=811} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:28 [STATS] ticks=1016 triangles_seen=1359 cand=2 exec=0 | scan_rejects={max_start_lt_100=26, no_quote_leg_2=211, no_quote_leg_3=225, no_quote_leg_1=895} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:33 [STATS] ticks=1126 triangles_seen=1479 cand=2 exec=0 | scan_rejects={max_start_lt_100=26, no_quote_leg_2=214, no_quote_leg_3=240, no_quote_leg_1=997} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:33 [REJECT] stage=scan reason=no_quote_leg_1 count=1000 tri=USDT->BTC->SNX
2026/04/10 09:00:38 [STATS] ticks=1247 triangles_seen=1627 cand=2 exec=0 | scan_rejects={max_start_lt_100=31, no_quote_leg_2=231, no_quote_leg_3=252, no_quote_leg_1=1111} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:43 [STATS] ticks=1351 triangles_seen=1771 cand=2 exec=0 | scan_rejects={max_start_lt_100=31, no_quote_leg_2=239, no_quote_leg_3=280, no_quote_leg_1=1219} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:48 [STATS] ticks=1473 triangles_seen=1951 cand=2 exec=0 | scan_rejects={max_start_lt_100=49, no_quote_leg_2=257, no_quote_leg_3=304, no_quote_leg_1=1339} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:53 [STATS] ticks=1621 triangles_seen=2156 cand=2 exec=0 | scan_rejects={max_start_lt_100=53, no_quote_leg_2=271, no_quote_leg_3=331, no_quote_leg_1=1499} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:00:58 [STATS] ticks=1795 triangles_seen=2397 cand=2 exec=0 | scan_rejects={max_start_lt_100=73, no_quote_leg_2=289, no_quote_leg_3=379, no_quote_leg_1=1654} | exec_rejects={profit_below_threshold=2}
2026/04/10 09:01:01 [EXEC CMP] USDT→DOT→KCS | est=-0.6815% | ideal=-0.3829% | rounded=-0.3831% | strict=-0.6830% | tails=-0.6828% | tail_usdt=0.00018350 | tails_map={USDT=0.00018350}
2026/04/10 09:01:01 [EXEC CMP] USDT→DOT→KCS | est=-0.6967% | ideal=-0.3982% | rounded=-0.3991% | strict=-0.6981% | tails=-0.6979% | tail_usdt=0.00019860 | tails_map={USDT=0.00019860}
2026/04/10 09:01:03 [STATS] ticks=1954 triangles_seen=2601 cand=4 exec=0 | scan_rejects={max_start_lt_100=78, no_quote_leg_2=315, no_quote_leg_3=405, no_quote_leg_1=1799} | exec_rejects={profit_below_threshold=4}
2026/04/10 09:01:08 [EXEC CMP] USDT→LINK→BTC | est=-0.3742% | ideal=-0.0747% | rounded=-0.0753% | strict=-0.3755% | tails=-0.3753% | tail_usdt=0.00020082 | tails_map={USDT=0.00020082}
2026/04/10 09:01:08 [STATS] ticks=2093 triangles_seen=2812 cand=5 exec=0 | scan_rejects={max_start_lt_100=83, no_quote_leg_2=383, no_quote_leg_3=440, no_quote_leg_1=1901} | exec_rejects={profit_below_threshold=5}
2026/04/10 09:01:08 [EXEC CMP] USDT→LINK→BTC | est=-0.3742% | ideal=-0.0747% | rounded=-0.0753% | strict=-0.3755% | tails=-0.3753% | tail_usdt=0.00020082 | tails_map={USDT=0.00020082}
2026/04/10 09:01:08 [EXEC CMP] USDT→XLM→ETH | est=-0.3855% | ideal=-0.0860% | rounded=-0.0868% | strict=-0.3880% | tails=-0.3872% | tail_usdt=0.00080027 | tails_map={USDT=0.00080027}
2026/04/10 09:01:08 [EXEC CMP] USDT→LINK→BTC | est=-0.3742% | ideal=-0.0747% | rounded=-0.0753% | strict=-0.3755% | tails=-0.3753% | tail_usdt=0.00020082 | tails_map={USDT=0.00020082}
2026/04/10 09:01:08 [EXEC CMP] USDT→BTC→ICP | est=-0.7874% | ideal=-0.4891% | rounded=-0.4900% | strict=-0.7890% | tails=-0.7881% | tail_usdt=0.00087300 | tails_map={USDT=0.00087300}
2026/04/10 09:01:08 [EXEC CMP] USDT→LINK→BTC | est=-0.3697% | ideal=-0.0702% | rounded=-0.0710% | strict=-0.3712% | tails=-0.3710% | tail_usdt=0.00020050 | tails_map={USDT=0.00020050}
2026/04/10 09:01:08 [REJECT] stage=exec reason=profit_below_threshold count=10 tri=USDT->LINK->BTC
2026/04/10 09:01:08 [EXEC CMP] USDT→LINK→BTC | est=-0.3697% | ideal=-0.0702% | rounded=-0.0710% | strict=-0.3712% | tails=-0.3710% | tail_usdt=0.00020050 | tails_map={USDT=0.00020050}
2026/04/10 09:01:08 [EXEC CMP] USDT→BTC→ICP | est=-0.7312% | ideal=-0.4328% | rounded=-0.4340% | strict=-0.7330% | tails=-0.7322% | tail_usdt=0.00081700 | tails_map={USDT=0.00081700}
2026/04/10 09:01:08 [EXEC CMP] USDT→BTC→AVAX | est=-0.3910% | ideal=-0.0916% | rounded=-0.0930% | strict=-0.3930% | tails=-0.3925% | tail_usdt=0.00047700 | tails_map={USDT=0.00047700}
2026/04/10 09:01:08 [EXEC CMP] USDT→XLM→ETH | est=-0.3855% | ideal=-0.0860% | rounded=-0.0868% | strict=-0.3880% | tails=-0.3872% | tail_usdt=0.00080027 | tails_map={USDT=0.00080027}
2026/04/10 09:01:08 [EXEC CMP] USDT→BTC→XLM | est=-0.4361% | ideal=-0.1368% | rounded=-0.1382% | strict=-0.4376% | tails=-0.4374% | tail_usdt=0.00022190 | tails_map={USDT=0.00022190}
2026/04/10 09:01:08 [EXEC CMP] USDT→BTC→DOT | est=-0.4468% | ideal=-0.1475% | rounded=-0.1478% | strict=-0.4475% | tails=-0.4473% | tail_usdt=0.00023180 | tails_map={USDT=0.00023180}
2026/04/10 09:01:08 [EXEC CMP] USDT→BTC→AVAX | est=-0.3910% | ideal=-0.0916% | rounded=-0.0930% | strict=-0.3930% | tails=-0.3925% | tail_usdt=0.00047700 | tails_map={USDT=0.00047700}
2026/04/10 09:01:08 [EXEC CMP] USDT→BTC→DASH | est=-0.4347% | ideal=-0.1354% | rounded=-0.1400% | strict=-0.4400% | tails=-0.4395% | tail_usdt=0.00052400 | tails_map={USDT=0.00052400}
2026/04/10 09:01:08 [EXEC CMP] USDT→BTC→AVAX | est=-0.4016% | ideal=-0.1022% | rounded=-0.1040% | strict=-0.4040% | tails=-0.4035% | tail_usdt=0.00048800 | tails_map={USDT=0.00048800}
2026/04/10 09:01:08 [REJECT] stage=scan reason=max_start_lt_100 count=100 tri=USDT->LYX->ETH
2026/04/10 09:01:08 [EXEC CMP] USDT→DOT→KCS | est=-0.5892% | ideal=-0.2903% | rounded=-0.2913% | strict=-0.5904% | tails=-0.5902% | tail_usdt=0.00019080 | tails_map={USDT=0.00019080}
2026/04/10 09:01:10 [EXEC CMP] USDT→DOT→KCS | est=-0.8416% | ideal=-0.5435% | rounded=-0.5440% | strict=-0.8427% | tails=-0.8426% | tail_usdt=0.00014340 | tails_map={USDT=0.00014340}
2026/04/10 09:01:10 [EXEC CMP] USDT→LINK→BTC | est=-0.3851% | ideal=-0.0856% | rounded=-0.0860% | strict=-0.3862% | tails=-0.3861% | tail_usdt=0.00010050 | tails_map={USDT=0.00010050}
2026/04/10 09:01:10 [EXEC CMP] USDT→BTC→ATOM | est=-0.4948% | ideal=-0.1957% | rounded=-0.1958% | strict=-0.4954% | tails=-0.4952% | tail_usdt=0.00016670 | tails_map={USDT=0.00016670}
2026/04/10 09:01:10 [EXEC CMP] USDT→BTC→ATOM | est=-0.4948% | ideal=-0.1957% | rounded=-0.1958% | strict=-0.4954% | tails=-0.4952% | tail_usdt=0.00016670 | tails_map={USDT=0.00016670}
2026/04/10 09:01:10 [EXEC CMP] USDT→DOT→KCS | est=-0.8416% | ideal=-0.5435% | rounded=-0.5440% | strict=-0.8427% | tails=-0.8426% | tail_usdt=0.00014340 | tails_map={USDT=0.00014340}
2026/04/10 09:01:10 [EXEC CMP] USDT→DOT→KCS | est=-0.8416% | ideal=-0.5435% | rounded=-0.5440% | strict=-0.8427% | tails=-0.8426% | tail_usdt=0.00014340 | tails_map={USDT=0.00014340}
2026/04/10 09:01:10 [EXEC CMP] USDT→LINK→BTC | est=-0.3895% | ideal=-0.0901% | rounded=-0.0904% | strict=-0.3905% | tails=-0.3904% | tail_usdt=0.00010183 | tails_map={USDT=0.00010183}
2026/04/10 09:01:10 [EXEC CMP] USDT→LINK→BTC | est=-0.3895% | ideal=-0.0901% | rounded=-0.0904% | strict=-0.3905% | tails=-0.3904% | tail_usdt=0.00010183 | tails_map={USDT=0.00010183}
2026/04/10 09:01:11 [REJECT] stage=scan reason=no_quote_leg_1 count=2000 tri=USDT->KCS->SOL
2026/04/10 09:01:11 [EXEC CMP] USDT→DOT→KCS | est=-0.8298% | ideal=-0.5317% | rounded=-0.5322% | strict=-0.8309% | tails=-0.8308% | tail_usdt=0.00013160 | tails_map={USDT=0.00013160}
2026/04/10 09:01:11 [EXEC CMP] USDT→DOT→KCS | est=-0.8298% | ideal=-0.5317% | rounded=-0.5322% | strict=-0.8309% | tails=-0.8308% | tail_usdt=0.00013160 | tails_map={USDT=0.00013160}
2026/04/10 09:01:11 [EXEC CMP] USDT→DOT→KCS | est=-0.7653% | ideal=-0.4670% | rounded=-0.4674% | strict=-0.7671% | tails=-0.7669% | tail_usdt=0.00016770 | tails_map={USDT=0.00016770}
2026/04/10 09:01:12 [EXEC CMP] USDT→DOT→KCS | est=-0.7653% | ideal=-0.4670% | rounded=-0.4674% | strict=-0.7671% | tails=-0.7669% | tail_usdt=0.00016770 | tails_map={USDT=0.00016770}
2026/04/10 09:01:12 [EXEC CMP] USDT→DOT→KCS | est=-0.7653% | ideal=-0.4670% | rounded=-0.4674% | strict=-0.7671% | tails=-0.7669% | tail_usdt=0.00016770 | tails_map={USDT=0.00016770}
2026/04/10 09:01:12 [EXEC CMP] USDT→BTC→FET | est=-0.5807% | ideal=-0.2818% | rounded=-0.2820% | strict=-0.5812% | tails=-0.5810% | tail_usdt=0.00020860 | tails_map={USDT=0.00020860}
2026/04/10 09:01:12 [EXEC CMP] USDT→BTC→DASH | est=-0.4814% | ideal=-0.1822% | rounded=-0.1900% | strict=-0.4900% | tails=-0.4895% | tail_usdt=0.00051700 | tails_map={USDT=0.00051700}
2026/04/10 09:01:12 [EXEC CMP] USDT→BTC→FET | est=-0.5210% | ideal=-0.2219% | rounded=-0.2221% | strict=-0.5214% | tails=-0.5213% | tail_usdt=0.00014880 | tails_map={USDT=0.00014880}
2026/04/10 09:01:12 [EXEC CMP] USDT→DOT→KCS | est=-0.7459% | ideal=-0.4476% | rounded=-0.4480% | strict=-0.7469% | tails=-0.7468% | tail_usdt=0.00014750 | tails_map={USDT=0.00014750}
2026/04/10 09:01:12 [EXEC CMP] USDT→BTC→ICP | est=-0.6176% | ideal=-0.3188% | rounded=-0.3190% | strict=-0.6190% | tails=-0.6184% | tail_usdt=0.00064600 | tails_map={USDT=0.00064600}
2026/04/10 09:01:12 [EXEC CMP] USDT→DOT→KCS | est=-0.8104% | ideal=-0.5123% | rounded=-0.5128% | strict=-0.8116% | tails=-0.8115% | tail_usdt=0.00011230 | tails_map={USDT=0.00011230}
2026/04/10 09:01:13 [EXEC CMP] USDT→DOT→KCS | est=-0.8104% | ideal=-0.5123% | rounded=-0.5128% | strict=-0.8116% | tails=-0.8115% | tail_usdt=0.00011230 | tails_map={USDT=0.00011230}
2026/04/10 09:01:13 [EXEC CMP] USDT→BTC→ICP | est=-0.6176% | ideal=-0.3188% | rounded=-0.3190% | strict=-0.6190% | tails=-0.6184% | tail_usdt=0.00064600 | tails_map={USDT=0.00064600}
2026/04/10 09:01:13 [STATS] ticks=2452 triangles_seen=3325 cand=41 exec=0 | scan_rejects={max_start_lt_100=149, no_quote_leg_2=524, no_quote_leg_3=531, no_quote_leg_1=2080} | exec_rejects={profit_below_threshold=41}
2026/04/10 09:01:13 [EXEC CMP] USDT→DOT→KCS | est=-0.6662% | ideal=-0.3676% | rounded=-0.3680% | strict=-0.6678% | tails=-0.6676% | tail_usdt=0.00016830 | tails_map={USDT=0.00016830}
2026/04/10 09:01:18 [STATS] ticks=2613 triangles_seen=3527 cand=42 exec=0 | scan_rejects={max_start_lt_100=154, no_quote_leg_2=547, no_quote_leg_3=565, no_quote_leg_1=2219} | exec_rejects={profit_below_threshold=42}
2026/04/10 09:01:22 [EXEC CMP] USDT→DOT→KCS | est=-0.6613% | ideal=-0.3627% | rounded=-0.3628% | strict=-0.6627% | tails=-0.6625% | tail_usdt=0.00016320 | tails_map={USDT=0.00016320}
2026/04/10 09:01:23 [EXEC CMP] USDT→DOT→KCS | est=-0.6537% | ideal=-0.3550% | rounded=-0.3552% | strict=-0.6551% | tails=-0.6549% | tail_usdt=0.00015560 | tails_map={USDT=0.00015560}
2026/04/10 09:01:23 [STATS] ticks=2689 triangles_seen=3617 cand=44 exec=0 | scan_rejects={max_start_lt_100=154, no_quote_leg_2=560, no_quote_leg_3=568, no_quote_leg_1=2291} | exec_rejects={profit_below_threshold=44}
2026/04/10 09:01:27 [EXEC CMP] USDT→SHIB→DOGE | est=-0.4935% | ideal=-0.1944% | rounded=-0.1944% | strict=-0.4935% | tails=-0.4935% | tail_usdt=0.00000390 | tails_map={USDT=0.00000390}
2026/04/10 09:01:27 [EXEC CMP] USDT→SHIB→DOGE | est=-0.4935% | ideal=-0.1944% | rounded=-0.1944% | strict=-0.4935% | tails=-0.4935% | tail_usdt=0.00000390 | tails_map={USDT=0.00000390}
2026/04/10 09:01:28 [STATS] ticks=2820 triangles_seen=3772 cand=46 exec=0 | scan_rejects={max_start_lt_100=158, no_quote_leg_2=580, no_quote_leg_3=603, no_quote_leg_1=2385} | exec_rejects={profit_below_threshold=46}
2026/04/10 09:01:33 [STATS] ticks=2999 triangles_seen=3996 cand=46 exec=0 | scan_rejects={max_start_lt_100=160, no_quote_leg_2=597, no_quote_leg_3=639, no_quote_leg_1=2554} | exec_rejects={profit_below_threshold=46}
2026/04/10 09:01:34 [EXEC CMP] USDT→DOT→KCS | est=-0.8298% | ideal=-0.5317% | rounded=-0.5322% | strict=-0.8309% | tails=-0.8308% | tail_usdt=0.00013160 | tails_map={USDT=0.00013160}
2026/04/10 09:01:37 [EXEC CMP] USDT→LINK→BTC | est=-0.3662% | ideal=-0.0667% | rounded=-0.0674% | strict=-0.3676% | tails=-0.3674% | tail_usdt=0.00020089 | tails_map={USDT=0.00020089}
2026/04/10 09:01:37 [EXEC CMP] USDT→BTC→AVAX | est=-0.3663% | ideal=-0.0668% | rounded=-0.0680% | strict=-0.3690% | tails=-0.3685% | tail_usdt=0.00045300 | tails_map={USDT=0.00045300}
2026/04/10 09:01:37 [EXEC CMP] USDT→BTC→DOT | est=-0.4086% | ideal=-0.1092% | rounded=-0.1094% | strict=-0.4092% | tails=-0.4090% | tail_usdt=0.00019350 | tails_map={USDT=0.00019350}
2026/04/10 09:01:37 [EXEC CMP] USDT→BTC→DOT | est=-0.3856% | ideal=-0.0862% | rounded=-0.0864% | strict=-0.3863% | tails=-0.3860% | tail_usdt=0.00027050 | tails_map={USDT=0.00027050}
2026/04/10 09:01:37 [EXEC CMP] USDT→BTC→NEAR | est=-0.3922% | ideal=-0.0928% | rounded=-0.0931% | strict=-0.3930% | tails=-0.3927% | tail_usdt=0.00027720 | tails_map={USDT=0.00027720}
2026/04/10 09:01:37 [EXEC CMP] USDT→BTC→DOT | est=-0.3856% | ideal=-0.0862% | rounded=-0.0864% | strict=-0.3863% | tails=-0.3860% | tail_usdt=0.00027050 | tails_map={USDT=0.00027050}
2026/04/10 09:01:38 [STATS] ticks=3250 triangles_seen=4339 cand=53 exec=0 | scan_rejects={max_start_lt_100=190, no_quote_leg_2=674, no_quote_leg_3=677, no_quote_leg_1=2745} | exec_rejects={profit_below_threshold=53}
2026/04/10 09:01:38 [EXEC CMP] USDT→BTC→AVAX | est=-0.3734% | ideal=-0.0739% | rounded=-0.0750% | strict=-0.3760% | tails=-0.3752% | tail_usdt=0.00082200 | tails_map={USDT=0.00082200}
2026/04/10 09:01:40 [EXEC CMP] USDT→BTC→AVAX | est=-0.3881% | ideal=-0.0887% | rounded=-0.0900% | strict=-0.3910% | tails=-0.3902% | tail_usdt=0.00078900 | tails_map={USDT=0.00078900}
2026/04/10 09:01:43 [STATS] ticks=3486 triangles_seen=4736 cand=55 exec=0 | scan_rejects={max_start_lt_100=213, no_quote_leg_3=715, no_quote_leg_2=830, no_quote_leg_1=2923} | exec_rejects={profit_below_threshold=55}
2026/04/10 09:01:46 [REJECT] stage=scan reason=no_quote_leg_1 count=3000 tri=USDT->BTC->EWT
2026/04/10 09:01:47 [EXEC CMP] USDT→SHIB→DOGE | est=-0.6952% | ideal=-0.3967% | rounded=-0.3967% | strict=-0.6952% | tails=-0.6952% | tail_usdt=0.00000578 | tails_map={USDT=0.00000578}
2026/04/10 09:01:48 [EXEC CMP] USDT→SHIB→DOGE | est=-0.7119% | ideal=-0.4134% | rounded=-0.4134% | strict=-0.7119% | tails=-0.7119% | tail_usdt=0.00000253 | tails_map={USDT=0.00000253}
2026/04/10 09:01:48 [STATS] ticks=3640 triangles_seen=4951 cand=57 exec=0 | scan_rejects={max_start_lt_100=215, no_quote_leg_3=752, no_quote_leg_2=893, no_quote_leg_1=3034} | exec_rejects={profit_below_threshold=57}
2026/04/10 09:01:48 [EXEC CMP] USDT→SHIB→DOGE | est=-0.7226% | ideal=-0.4242% | rounded=-0.4242% | strict=-0.7226% | tails=-0.7226% | tail_usdt=0.00000324 | tails_map={USDT=0.00000324}
2026/04/10 09:01:48 [EXEC CMP] USDT→SHIB→DOGE | est=-0.7226% | ideal=-0.4242% | rounded=-0.4242% | strict=-0.7226% | tails=-0.7226% | tail_usdt=0.00000324 | tails_map={USDT=0.00000324}
2026/04/10 09:01:53 [EXEC CMP] USDT→DOT→KCS | est=-0.8062% | ideal=-0.5080% | rounded=-0.5085% | strict=-0.8073% | tails=-0.8072% | tail_usdt=0.00010800 | tails_map={USDT=0.00010800}
2026/04/10 09:01:53 [STATS] ticks=3818 triangles_seen=5160 cand=60 exec=0 | scan_rejects={max_start_lt_100=223, no_quote_leg_3=783, no_quote_leg_2=915, no_quote_leg_1=3179} | exec_rejects={profit_below_threshold=60}
2026/04/10 09:01:53 [EXEC CMP] USDT→SHIB→DOGE | est=-1.3717% | ideal=-1.0753% | rounded=-1.0753% | strict=-1.3717% | tails=-1.3717% | tail_usdt=0.00000301 | tails_map={USDT=0.00000301}
2026/04/10 09:01:54 [EXEC CMP] USDT→SHIB→DOGE | est=-1.3717% | ideal=-1.0753% | rounded=-1.0753% | strict=-1.3717% | tails=-1.3717% | tail_usdt=0.00000301 | tails_map={USDT=0.00000301}
2026/04/10 09:01:54 [EXEC CMP] USDT→DOT→KCS | est=-0.7009% | ideal=-0.4024% | rounded=-0.4025% | strict=-0.7023% | tails=-0.7022% | tail_usdt=0.00010290 | tails_map={USDT=0.00010290}
2026/04/10 09:01:54 [EXEC CMP] USDT→DOT→KCS | est=-0.6856% | ideal=-0.3871% | rounded=-0.3874% | strict=-0.6872% | tails=-0.6870% | tail_usdt=0.00018770 | tails_map={USDT=0.00018770}
2026/04/10 09:01:57 [EXEC CMP] USDT→DOT→KCS | est=-0.6856% | ideal=-0.3871% | rounded=-0.3874% | strict=-0.6872% | tails=-0.6870% | tail_usdt=0.00018770 | tails_map={USDT=0.00018770}
2026/04/10 09:01:57 [EXEC CMP] USDT→DOT→KCS | est=-0.6856% | ideal=-0.3871% | rounded=-0.3874% | strict=-0.6872% | tails=-0.6870% | tail_usdt=0.00018770 | tails_map={USDT=0.00018770}
2026/04/10 09:01:57 [EXEC CMP] USDT→DOT→KCS | est=-0.6856% | ideal=-0.3871% | rounded=-0.3874% | strict=-0.6872% | tails=-0.6870% | tail_usdt=0.00018770 | tails_map={USDT=0.00018770}
2026/04/10 09:01:58 [STATS] ticks=4008 triangles_seen=5401 cand=67 exec=0 | scan_rejects={max_start_lt_100=236, no_quote_leg_3=799, no_quote_leg_2=952, no_quote_leg_1=3347} | exec_rejects={profit_below_threshold=67}
2026/04/10 09:02:03 [STATS] ticks=4180 triangles_seen=5613 cand=67 exec=0 | scan_rejects={max_start_lt_100=241, no_quote_leg_3=826, no_quote_leg_2=976, no_quote_leg_1=3503} | exec_rejects={profit_below_threshold=67}
2026/04/10 09:02:07 [REJECT] stage=scan reason=no_quote_leg_2 count=1000 tri=USDT->BTC->ICP
2026/04/10 09:02:07 [EXEC CMP] USDT→XRP→ETH | est=-0.3402% | ideal=-0.0407% | rounded=-0.0410% | strict=-0.3407% | tails=-0.3406% | tail_usdt=0.00013294 | tails_map={USDT=0.00013294}
2026/04/10 09:02:07 [EXEC CMP] USDT→XRP→ETH | est=-0.3547% | ideal=-0.0552% | rounded=-0.0555% | strict=-0.3552% | tails=-0.3551% | tail_usdt=0.00013247 | tails_map={USDT=0.00013247}
2026/04/10 09:02:08 [EXEC CMP] USDT→XRP→ETH | est=-0.3540% | ideal=-0.0545% | rounded=-0.0546% | strict=-0.3543% | tails=-0.3543% | tail_usdt=0.00007059 | tails_map={USDT=0.00007059}
2026/04/10 09:02:08 [STATS] ticks=4354 triangles_seen=5955 cand=70 exec=0 | scan_rejects={max_start_lt_100=260, no_quote_leg_3=874, no_quote_leg_2=1124, no_quote_leg_1=3627} | exec_rejects={profit_below_threshold=70}
2026/04/10 09:02:13 [STATS] ticks=4482 triangles_seen=6132 cand=70 exec=0 | scan_rejects={max_start_lt_100=264, no_quote_leg_3=905, no_quote_leg_2=1158, no_quote_leg_1=3735} | exec_rejects={profit_below_threshold=70}
2026/04/10 09:02:18 [STATS] ticks=4601 triangles_seen=6280 cand=70 exec=0 | scan_rejects={max_start_lt_100=264, no_quote_leg_3=926, no_quote_leg_2=1170, no_quote_leg_1=3850} | exec_rejects={profit_below_threshold=70}
2026/04/10 09:02:23 [STATS] ticks=4693 triangles_seen=6388 cand=70 exec=0 | scan_rejects={max_start_lt_100=264, no_quote_leg_3=937, no_quote_leg_2=1180, no_quote_leg_1=3937} | exec_rejects={profit_below_threshold=70}
2026/04/10 09:02:26 [REJECT] stage=scan reason=no_quote_leg_1 count=4000 tri=USDT->BTC->XDC
2026/04/10 09:02:28 [STATS] ticks=4781 triangles_seen=6493 cand=70 exec=0 | scan_rejects={max_start_lt_100=264, no_quote_leg_3=947, no_quote_leg_2=1193, no_quote_leg_1=4019} | exec_rejects={profit_below_threshold=70}
2026/04/10 09:02:30 [EXEC CMP] USDT→BTC→AVAX | est=-0.3820% | ideal=-0.0826% | rounded=-0.0840% | strict=-0.3840% | tails=-0.3836% | tail_usdt=0.00044200 | tails_map={USDT=0.00044200}
2026/04/10 09:02:31 [EXEC CMP] USDT→BTC→AVAX | est=-0.3711% | ideal=-0.0716% | rounded=-0.0730% | strict=-0.3730% | tails=-0.3724% | tail_usdt=0.00055800 | tails_map={USDT=0.00055800}
2026/04/10 09:02:31 [EXEC CMP] USDT→BTC→AVAX | est=-0.4010% | ideal=-0.1016% | rounded=-0.1030% | strict=-0.4030% | tails=-0.4023% | tail_usdt=0.00072400 | tails_map={USDT=0.00072400}
2026/04/10 09:02:31 [EXEC CMP] USDT→BTC→AVAX | est=-0.4117% | ideal=-0.1123% | rounded=-0.1140% | strict=-0.4140% | tails=-0.4133% | tail_usdt=0.00073500 | tails_map={USDT=0.00073500}
2026/04/10 09:02:31 [REJECT] stage=scan reason=max_start_zero count=1 tri=USDT->BTC->BDX
2026/04/10 09:02:31 [EXEC CMP] USDT→BTC→AVAX | est=-0.4010% | ideal=-0.1016% | rounded=-0.1030% | strict=-0.4030% | tails=-0.4023% | tail_usdt=0.00072400 | tails_map={USDT=0.00072400}
2026/04/10 09:02:33 [STATS] ticks=4946 triangles_seen=6841 cand=75 exec=0 | scan_rejects={max_start_zero=1, max_start_lt_100=302, no_quote_leg_3=972, no_quote_leg_2=1327, no_quote_leg_1=4164} | exec_rejects={profit_below_threshold=75}
2026/04/10 09:02:33 [EXEC CMP] USDT→DOT→KCS | est=-0.5032% | ideal=-0.2042% | rounded=-0.2047% | strict=-0.5047% | tails=-0.5046% | tail_usdt=0.00010510 | tails_map={USDT=0.00010510}
2026/04/10 09:02:37 [REJECT] stage=scan reason=no_quote_leg_3 count=1000 tri=USDT->LYX->ETH
2026/04/10 09:02:38 [STATS] ticks=5078 triangles_seen=7022 cand=76 exec=0 | scan_rejects={max_start_zero=1, max_start_lt_100=309, no_quote_leg_3=1004, no_quote_leg_2=1350, no_quote_leg_1=4282} | exec_rejects={profit_below_threshold=76}
2026/04/10 09:02:38 [EXEC CMP] USDT→BTC→AVAX | est=-0.4057% | ideal=-0.1063% | rounded=-0.1070% | strict=-0.4080% | tails=-0.4072% | tail_usdt=0.00077300 | tails_map={USDT=0.00077300}
2026/04/10 09:02:38 [EXEC CMP] USDT→BTC→AVAX | est=-0.3951% | ideal=-0.0957% | rounded=-0.0970% | strict=-0.3970% | tails=-0.3962% | tail_usdt=0.00076200 | tails_map={USDT=0.00076200}
2026/04/10 09:02:40 [EXEC CMP] USDT→BTC→DASH | est=-0.5102% | ideal=-0.2112% | rounded=-0.2200% | strict=-0.5200% | tails=-0.5193% | tail_usdt=0.00071700 | tails_map={USDT=0.00071700}
2026/04/10 09:02:40 [EXEC CMP] USDT→BTC→KNC | est=-0.9274% | ideal=-0.6296% | rounded=-0.6298% | strict=-0.9279% | tails=-0.9277% | tail_usdt=0.00015570 | tails_map={USDT=0.00015570}
2026/04/10 09:02:40 [EXEC CMP] USDT→BTC→AVAX | est=-0.3773% | ideal=-0.0778% | rounded=-0.0790% | strict=-0.3800% | tails=-0.3796% | tail_usdt=0.00040700 | tails_map={USDT=0.00040700}
2026/04/10 09:02:40 [EXEC CMP] USDT→BTC→AVAX | est=-0.3666% | ideal=-0.0671% | rounded=-0.0680% | strict=-0.3690% | tails=-0.3686% | tail_usdt=0.00039600 | tails_map={USDT=0.00039600}
2026/04/10 09:02:43 [EXEC CMP] USDT→SHIB→DOGE | est=-1.2316% | ideal=-0.9347% | rounded=-0.9347% | strict=-1.2316% | tails=-1.2316% | tail_usdt=0.00000273 | tails_map={USDT=0.00000273}
2026/04/10 09:02:43 [STATS] ticks=5323 triangles_seen=7497 cand=83 exec=0 | scan_rejects={max_start_zero=1, max_start_lt_100=343, no_quote_leg_3=1062, no_quote_leg_2=1551, no_quote_leg_1=4457} | exec_rejects={profit_below_threshold=83}
2026/04/10 09:02:48 [STATS] ticks=5480 triangles_seen=7725 cand=83 exec=0 | scan_rejects={max_start_zero=1, max_start_lt_100=357, no_quote_leg_3=1086, no_quote_leg_2=1612, no_quote_leg_1=4586} | exec_rejects={profit_below_threshold=83}
2026/04/10 09:02:49 [EXEC CMP] USDT→BTC→AVAX | est=-0.3614% | ideal=-0.0619% | rounded=-0.0630% | strict=-0.3640% | tails=-0.3631% | tail_usdt=0.00094800 | tails_map={USDT=0.00094800}
2026/04/10 09:02:49 [EXEC CMP] USDT→BTC→AVAX | est=-0.3614% | ideal=-0.0619% | rounded=-0.0630% | strict=-0.3640% | tails=-0.3631% | tail_usdt=0.00094800 | tails_map={USDT=0.00094800}
2026/04/10 09:02:49 [EXEC CMP] USDT→BTC→AVAX | est=-0.3690% | ideal=-0.0695% | rounded=-0.0710% | strict=-0.3710% | tails=-0.3705% | tail_usdt=0.00054000 | tails_map={USDT=0.00054000}
2026/04/10 09:02:53 [STATS] ticks=5639 triangles_seen=8019 cand=86 exec=0 | scan_rejects={max_start_zero=1, max_start_lt_100=363, no_quote_leg_3=1109, no_quote_leg_2=1729, no_quote_leg_1=4731} | exec_rejects={profit_below_threshold=86}
2026/04/10 09:02:55 [EXEC CMP] USDT→DOT→KCS | est=-0.6689% | ideal=-0.3703% | rounded=-0.3712% | strict=-0.6703% | tails=-0.6701% | tail_usdt=0.00017080 | tails_map={USDT=0.00017080}
2026/04/10 09:02:57 [EXEC CMP] USDT→BTC→AVAX | est=-0.3621% | ideal=-0.0626% | rounded=-0.0640% | strict=-0.3640% | tails=-0.3636% | tail_usdt=0.00039700 | tails_map={USDT=0.00039700}
2026/04/10 09:02:57 [EXEC CMP] USDT→BTC→NEAR | est=-0.4177% | ideal=-0.1184% | rounded=-0.1186% | strict=-0.4183% | tails=-0.4181% | tail_usdt=0.00015160 | tails_map={USDT=0.00015160}
2026/04/10 09:02:57 [EXEC CMP] USDT→BTC→NEAR | est=-0.4177% | ideal=-0.1184% | rounded=-0.1186% | strict=-0.4183% | tails=-0.4181% | tail_usdt=0.00015160 | tails_map={USDT=0.00015160}
2026/04/10 09:02:57 [EXEC CMP] USDT→BTC→AVAX | est=-0.3515% | ideal=-0.0519% | rounded=-0.0530% | strict=-0.3540% | tails=-0.3536% | tail_usdt=0.00038700 | tails_map={USDT=0.00038700}
2026/04/10 09:02:57 [EXEC CMP] USDT→BTC→NEAR | est=-0.4177% | ideal=-0.1184% | rounded=-0.1186% | strict=-0.4183% | tails=-0.4181% | tail_usdt=0.00015160 | tails_map={USDT=0.00015160}
2026/04/10 09:02:57 [EXEC CMP] USDT→BTC→NEAR | est=-0.4177% | ideal=-0.1184% | rounded=-0.1186% | strict=-0.4183% | tails=-0.4181% | tail_usdt=0.00015160 | tails_map={USDT=0.00015160}
2026/04/10 09:02:58 [STATS] ticks=5764 triangles_seen=8283 cand=93 exec=0 | scan_rejects={max_start_zero=1, max_start_lt_100=376, no_quote_leg_3=1137, no_quote_leg_2=1830, no_quote_leg_1=4846} | exec_rejects={profit_below_threshold=93}
2026/04/10 09:02:59 [EXEC CMP] USDT→AVA→ETH | est=-1.2300% | ideal=-0.9331% | rounded=-0.9349% | strict=-1.2328% | tails=-1.2310% | tail_usdt=0.00180092 | tails_map={USDT=0.00180092}
2026/04/10 09:03:00 [EXEC CMP] USDT→DOT→KCS | est=-0.7105% | ideal=-0.4121% | rounded=-0.4125% | strict=-0.7114% | tails=-0.7113% | tail_usdt=0.00011200 | tails_map={USDT=0.00011200}
2026/04/10 09:03:00 [EXEC CMP] USDT→DOT→KCS | est=-0.7029% | ideal=-0.4044% | rounded=-0.4049% | strict=-0.7039% | tails=-0.7038% | tail_usdt=0.00010450 | tails_map={USDT=0.00010450}
2026/04/10 09:03:01 [EXEC CMP] USDT→DOT→KCS | est=-0.7258% | ideal=-0.4274% | rounded=-0.4276% | strict=-0.7274% | tails=-0.7273% | tail_usdt=0.00012800 | tails_map={USDT=0.00012800}
2026/04/10 09:03:01 [EXEC CMP] USDT→DOT→KCS | est=-0.7334% | ideal=-0.4350% | rounded=-0.4352% | strict=-0.7350% | tails=-0.7349% | tail_usdt=0.00013560 | tails_map={USDT=0.00013560}
2026/04/10 09:03:01 [EXEC CMP] USDT→DOT→KCS | est=-0.7182% | ideal=-0.4197% | rounded=-0.4200% | strict=-0.7190% | tails=-0.7189% | tail_usdt=0.00011960 | tails_map={USDT=0.00011960}
2026/04/10 09:03:03 [REJECT] stage=scan reason=no_quote_leg_1 count=5000 tri=USDT->XLM->ETH
2026/04/10 09:03:03 [EXEC CMP] USDT→BTC→ATOM | est=-0.4442% | ideal=-0.1449% | rounded=-0.1453% | strict=-0.4450% | tails=-0.4447% | tail_usdt=0.00034130 | tails_map={USDT=0.00034130}
2026/04/10 09:03:03 [REJECT] stage=exec reason=profit_below_threshold count=100 tri=USDT->BTC->ATOM
2026/04/10 09:03:03 [STATS] ticks=5999 triangles_seen=8690 cand=100 exec=0 | scan_rejects={max_start_zero=2, max_start_lt_100=413, no_quote_leg_3=1200, no_quote_leg_2=1956, no_quote_leg_1=5019} | exec_rejects={profit_below_threshold=100}
2026/04/10 09:03:03 [EXEC CMP] USDT→LINK→BTC | est=-0.3945% | ideal=-0.0951% | rounded=-0.0957% | strict=-0.3960% | tails=-0.3959% | tail_usdt=0.00010125 | tails_map={USDT=0.00010125}
2026/04/10 09:03:05 [EXEC CMP] USDT→DOT→KCS | est=-0.7334% | ideal=-0.4350% | rounded=-0.4352% | strict=-0.7350% | tails=-0.7349% | tail_usdt=0.00013560 | tails_map={USDT=0.00013560}
2026/04/10 09:03:05 [EXEC CMP] USDT→DOT→KCS | est=-0.7258% | ideal=-0.4274% | rounded=-0.4276% | strict=-0.7274% | tails=-0.7273% | tail_usdt=0.00012800 | tails_map={USDT=0.00012800}
2026/04/10 09:03:05 [EXEC CMP] USDT→DOT→KCS | est=-0.7258% | ideal=-0.4274% | rounded=-0.4276% | strict=-0.7274% | tails=-0.7273% | tail_usdt=0.00012800 | tails_map={USDT=0.00012800}
2026/04/10 09:03:08 [STATS] ticks=6195 triangles_seen=8927 cand=104 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=430, no_quote_leg_3=1247, no_quote_leg_2=1992, no_quote_leg_1=5151} | exec_rejects={profit_below_threshold=104}
2026/04/10 09:03:09 [REJECT] stage=scan reason=no_quote_leg_2 count=2000 tri=USDT->XRP->BTC
2026/04/10 09:03:13 [STATS] ticks=6366 triangles_seen=9136 cand=104 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=436, no_quote_leg_3=1279, no_quote_leg_2=2020, no_quote_leg_1=5294} | exec_rejects={profit_below_threshold=104}
2026/04/10 09:03:18 [STATS] ticks=6478 triangles_seen=9274 cand=104 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=440, no_quote_leg_3=1297, no_quote_leg_2=2038, no_quote_leg_1=5392} | exec_rejects={profit_below_threshold=104}
2026/04/10 09:03:23 [STATS] ticks=6574 triangles_seen=9381 cand=104 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=440, no_quote_leg_3=1311, no_quote_leg_2=2057, no_quote_leg_1=5466} | exec_rejects={profit_below_threshold=104}
2026/04/10 09:03:23 [EXEC CMP] USDT→LINK→BTC | est=-0.4054% | ideal=-0.1060% | rounded=-0.1067% | strict=-0.4076% | tails=-0.4070% | tail_usdt=0.00060196 | tails_map={USDT=0.00060196}
2026/04/10 09:03:23 [EXEC CMP] USDT→BTC→ICP | est=-0.7677% | ideal=-0.4694% | rounded=-0.4710% | strict=-0.7700% | tails=-0.7688% | tail_usdt=0.00117400 | tails_map={USDT=0.00117400}
2026/04/10 09:03:24 [EXEC CMP] USDT→BTC→DOT | est=-0.4078% | ideal=-0.1085% | rounded=-0.1088% | strict=-0.4085% | tails=-0.4083% | tail_usdt=0.00023180 | tails_map={USDT=0.00023180}
2026/04/10 09:03:24 [EXEC CMP] USDT→BTC→AVAX | est=-0.3954% | ideal=-0.0960% | rounded=-0.0970% | strict=-0.3980% | tails=-0.3975% | tail_usdt=0.00052100 | tails_map={USDT=0.00052100}
2026/04/10 09:03:28 [STATS] ticks=6770 triangles_seen=9788 cand=108 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=472, no_quote_leg_3=1377, no_quote_leg_2=2248, no_quote_leg_1=5580} | exec_rejects={profit_below_threshold=108}
2026/04/10 09:03:33 [STATS] ticks=6884 triangles_seen=9931 cand=108 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=474, no_quote_leg_3=1402, no_quote_leg_2=2269, no_quote_leg_1=5675} | exec_rejects={profit_below_threshold=108}
2026/04/10 09:03:34 [EXEC CMP] USDT→BTC→DASH | est=-0.4512% | ideal=-0.1520% | rounded=-0.1600% | strict=-0.4600% | tails=-0.4594% | tail_usdt=0.00062800 | tails_map={USDT=0.00062800}
2026/04/10 09:03:34 [EXEC CMP] USDT→BTC→ICP | est=-0.6036% | ideal=-0.3048% | rounded=-0.3060% | strict=-0.6050% | tails=-0.6042% | tail_usdt=0.00077300 | tails_map={USDT=0.00077300}
2026/04/10 09:03:34 [EXEC CMP] USDT→BTC→ICP | est=-0.6036% | ideal=-0.3048% | rounded=-0.3060% | strict=-0.6050% | tails=-0.6042% | tail_usdt=0.00077300 | tails_map={USDT=0.00077300}
2026/04/10 09:03:38 [STATS] ticks=7007 triangles_seen=10112 cand=111 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=480, no_quote_leg_3=1424, no_quote_leg_2=2336, no_quote_leg_1=5758} | exec_rejects={profit_below_threshold=111}
2026/04/10 09:03:39 [EXEC CMP] USDT→ATOM→ETH | est=-0.4665% | ideal=-0.1673% | rounded=-0.1673% | strict=-0.4690% | tails=-0.4689% | tail_usdt=0.00010039 | tails_map={USDT=0.00010039}
2026/04/10 09:03:43 [STATS] ticks=7189 triangles_seen=10358 cand=112 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=485, no_quote_leg_3=1489, no_quote_leg_2=2348, no_quote_leg_1=5921} | exec_rejects={profit_below_threshold=112}
2026/04/10 09:03:47 [REJECT] stage=scan reason=no_quote_leg_1 count=6000 tri=USDT->BTC->ANKR
2026/04/10 09:03:48 [STATS] ticks=7292 triangles_seen=10487 cand=112 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=487, no_quote_leg_3=1500, no_quote_leg_2=2369, no_quote_leg_1=6016} | exec_rejects={profit_below_threshold=112}
2026/04/10 09:03:49 [EXEC CMP] USDT→SHIB→DOGE | est=-1.0200% | ideal=-0.7225% | rounded=-0.7225% | strict=-1.0201% | tails=-1.0201% | tail_usdt=0.00000098 | tails_map={USDT=0.00000098}
2026/04/10 09:03:49 [EXEC CMP] USDT→SHIB→DOGE | est=-1.0200% | ideal=-0.7225% | rounded=-0.7225% | strict=-1.0201% | tails=-1.0201% | tail_usdt=0.00000098 | tails_map={USDT=0.00000098}
2026/04/10 09:03:53 [STATS] ticks=7431 triangles_seen=10691 cand=114 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=489, no_quote_leg_3=1530, no_quote_leg_2=2389, no_quote_leg_1=6166} | exec_rejects={profit_below_threshold=114}
2026/04/10 09:03:58 [STATS] ticks=7547 triangles_seen=10849 cand=114 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=489, no_quote_leg_3=1558, no_quote_leg_2=2403, no_quote_leg_1=6282} | exec_rejects={profit_below_threshold=114}
2026/04/10 09:04:00 [EXEC CMP] USDT→DOT→KCS | est=-0.5442% | ideal=-0.2452% | rounded=-0.2458% | strict=-0.5449% | tails=-0.5448% | tail_usdt=0.00014530 | tails_map={USDT=0.00014530}
2026/04/10 09:04:00 [EXEC CMP] USDT→DOT→KCS | est=-0.5442% | ideal=-0.2452% | rounded=-0.2458% | strict=-0.5449% | tails=-0.5448% | tail_usdt=0.00014530 | tails_map={USDT=0.00014530}
2026/04/10 09:04:02 [EXEC CMP] USDT→DOT→KCS | est=-0.5206% | ideal=-0.2215% | rounded=-0.2221% | strict=-0.5213% | tails=-0.5212% | tail_usdt=0.00012170 | tails_map={USDT=0.00012170}
2026/04/10 09:04:02 [EXEC CMP] USDT→DOT→KCS | est=-0.5442% | ideal=-0.2452% | rounded=-0.2458% | strict=-0.5449% | tails=-0.5448% | tail_usdt=0.00014530 | tails_map={USDT=0.00014530}
2026/04/10 09:04:03 [STATS] ticks=7725 triangles_seen=11074 cand=118 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=492, no_quote_leg_3=1611, no_quote_leg_2=2430, no_quote_leg_1=6420} | exec_rejects={profit_below_threshold=118}
2026/04/10 09:04:05 [EXEC CMP] USDT→DOT→KCS | est=-0.5595% | ideal=-0.2606% | rounded=-0.2609% | strict=-0.5609% | tails=-0.5607% | tail_usdt=0.00016130 | tails_map={USDT=0.00016130}
2026/04/10 09:04:08 [STATS] ticks=7881 triangles_seen=11269 cand=119 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=498, no_quote_leg_3=1660, no_quote_leg_2=2464, no_quote_leg_1=6525} | exec_rejects={profit_below_threshold=119}
2026/04/10 09:04:13 [STATS] ticks=8039 triangles_seen=11474 cand=119 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=507, no_quote_leg_3=1691, no_quote_leg_2=2504, no_quote_leg_1=6650} | exec_rejects={profit_below_threshold=119}
2026/04/10 09:04:18 [STATS] ticks=8161 triangles_seen=11624 cand=119 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=507, no_quote_leg_3=1707, no_quote_leg_2=2515, no_quote_leg_1=6773} | exec_rejects={profit_below_threshold=119}
2026/04/10 09:04:23 [STATS] ticks=8306 triangles_seen=11859 cand=119 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=513, no_quote_leg_3=1744, no_quote_leg_2=2587, no_quote_leg_1=6893} | exec_rejects={profit_below_threshold=119}
2026/04/10 09:04:23 [EXEC CMP] USDT→SHIB→DOGE | est=-0.8648% | ideal=-0.5667% | rounded=-0.5668% | strict=-0.8648% | tails=-0.8648% | tail_usdt=0.00000552 | tails_map={USDT=0.00000552}
2026/04/10 09:04:23 [EXEC CMP] USDT→SHIB→DOGE | est=-0.8541% | ideal=-0.5560% | rounded=-0.5560% | strict=-0.8541% | tails=-0.8541% | tail_usdt=0.00000484 | tails_map={USDT=0.00000484}
2026/04/10 09:04:28 [STATS] ticks=8425 triangles_seen=12006 cand=121 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=515, no_quote_leg_3=1763, no_quote_leg_2=2615, no_quote_leg_1=6989} | exec_rejects={profit_below_threshold=121}
2026/04/10 09:04:28 [REJECT] stage=scan reason=no_quote_leg_1 count=7000 tri=USDT->BTC->BDX
2026/04/10 09:04:33 [STATS] ticks=8572 triangles_seen=12194 cand=121 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=518, no_quote_leg_3=1790, no_quote_leg_2=2634, no_quote_leg_1=7128} | exec_rejects={profit_below_threshold=121}
2026/04/10 09:04:38 [STATS] ticks=8689 triangles_seen=12343 cand=121 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=518, no_quote_leg_3=1816, no_quote_leg_2=2656, no_quote_leg_1=7229} | exec_rejects={profit_below_threshold=121}
2026/04/10 09:04:43 [STATS] ticks=8781 triangles_seen=12466 cand=121 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=522, no_quote_leg_3=1835, no_quote_leg_2=2679, no_quote_leg_1=7306} | exec_rejects={profit_below_threshold=121}
2026/04/10 09:04:48 [STATS] ticks=8984 triangles_seen=12725 cand=121 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=533, no_quote_leg_3=1890, no_quote_leg_2=2708, no_quote_leg_1=7470} | exec_rejects={profit_below_threshold=121}
2026/04/10 09:04:48 [EXEC CMP] USDT→LINK→BTC | est=-0.3555% | ideal=-0.0560% | rounded=-0.0564% | strict=-0.3567% | tails=-0.3566% | tail_usdt=0.00010197 | tails_map={USDT=0.00010197}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→DOT | est=-0.3394% | ideal=-0.0398% | rounded=-0.0399% | strict=-0.3398% | tails=-0.3397% | tail_usdt=0.00006800 | tails_map={USDT=0.00006800}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→AVAX | est=-0.3828% | ideal=-0.0833% | rounded=-0.0840% | strict=-0.3850% | tails=-0.3847% | tail_usdt=0.00031300 | tails_map={USDT=0.00031300}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→NEAR | est=-0.4427% | ideal=-0.1434% | rounded=-0.1436% | strict=-0.4433% | tails=-0.4432% | tail_usdt=0.00007160 | tails_map={USDT=0.00007160}
2026/04/10 09:04:48 [EXEC CMP] USDT→LINK→BTC | est=-0.3795% | ideal=-0.0801% | rounded=-0.0802% | strict=-0.3812% | tails=-0.3811% | tail_usdt=0.00010050 | tails_map={USDT=0.00010050}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→FET | est=-0.5308% | ideal=-0.2318% | rounded=-0.2319% | strict=-0.5312% | tails=-0.5311% | tail_usdt=0.00005960 | tails_map={USDT=0.00005960}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→AVAX | est=-0.3722% | ideal=-0.0727% | rounded=-0.0730% | strict=-0.3740% | tails=-0.3737% | tail_usdt=0.00030200 | tails_map={USDT=0.00030200}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→DOT | est=-0.3317% | ideal=-0.0321% | rounded=-0.0322% | strict=-0.3322% | tails=-0.3321% | tail_usdt=0.00006040 | tails_map={USDT=0.00006040}
2026/04/10 09:04:48 [EXEC CMP] USDT→LINK→BTC | est=-0.3784% | ideal=-0.0790% | rounded=-0.0795% | strict=-0.3798% | tails=-0.3795% | tail_usdt=0.00030005 | tails_map={USDT=0.00030005}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→NEAR | est=-0.4281% | ideal=-0.1288% | rounded=-0.1290% | strict=-0.4287% | tails=-0.4286% | tail_usdt=0.00005700 | tails_map={USDT=0.00005700}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→DOT | est=-0.3164% | ideal=-0.0168% | rounded=-0.0169% | strict=-0.3169% | tails=-0.3169% | tail_usdt=0.00004510 | tails_map={USDT=0.00004510}
2026/04/10 09:04:48 [EXEC CMP] USDT→DOT→KCS | est=-0.6213% | ideal=-0.3226% | rounded=-0.3233% | strict=-0.6223% | tails=-0.6222% | tail_usdt=0.00012280 | tails_map={USDT=0.00012280}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ICP | est=-0.4029% | ideal=-0.1035% | rounded=-0.1040% | strict=-0.4040% | tails=-0.4037% | tail_usdt=0.00033200 | tails_map={USDT=0.00033200}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ICP | est=-0.8253% | ideal=-0.5272% | rounded=-0.5280% | strict=-0.8270% | tails=-0.8262% | tail_usdt=0.00075500 | tails_map={USDT=0.00075500}
2026/04/10 09:04:48 [EXEC CMP] USDT→ETC→ETH | est=-0.4507% | ideal=-0.1514% | rounded=-0.1525% | strict=-0.4520% | tails=-0.4520% | tail_usdt=0.00000339 | tails_map={USDT=0.00000339}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→FET | est=-0.4893% | ideal=-0.1902% | rounded=-0.1903% | strict=-0.4898% | tails=-0.4897% | tail_usdt=0.00011810 | tails_map={USDT=0.00011810}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ICP | est=-0.8253% | ideal=-0.5272% | rounded=-0.5280% | strict=-0.8270% | tails=-0.8262% | tail_usdt=0.00075500 | tails_map={USDT=0.00075500}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→NEAR | est=-0.4281% | ideal=-0.1288% | rounded=-0.1290% | strict=-0.4287% | tails=-0.4286% | tail_usdt=0.00005700 | tails_map={USDT=0.00005700}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→DOT | est=-0.3332% | ideal=-0.0336% | rounded=-0.0342% | strict=-0.3342% | tails=-0.3336% | tail_usdt=0.00060340 | tails_map={USDT=0.00060340}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→FET | est=-0.5060% | ideal=-0.2069% | rounded=-0.2076% | strict=-0.5069% | tails=-0.5063% | tail_usdt=0.00057630 | tails_map={USDT=0.00057630}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→NEAR | est=-0.4448% | ideal=-0.1456% | rounded=-0.1463% | strict=-0.4460% | tails=-0.4454% | tail_usdt=0.00061530 | tails_map={USDT=0.00061530}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ATOM | est=-0.5000% | ideal=-0.2009% | rounded=-0.2016% | strict=-0.5010% | tails=-0.5004% | tail_usdt=0.00057040 | tails_map={USDT=0.00057040}
2026/04/10 09:04:48 [EXEC CMP] USDT→DOT→KCS | est=-0.6441% | ideal=-0.3455% | rounded=-0.3460% | strict=-0.6451% | tails=-0.6450% | tail_usdt=0.00014560 | tails_map={USDT=0.00014560}
2026/04/10 09:04:48 [EXEC CMP] USDT→ATOM→ETH | est=-0.6056% | ideal=-0.3069% | rounded=-0.3085% | strict=-0.6078% | tails=-0.6076% | tail_usdt=0.00020036 | tails_map={USDT=0.00020036}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→FET | est=-0.4761% | ideal=-0.1770% | rounded=-0.1776% | strict=-0.4771% | tails=-0.4765% | tail_usdt=0.00064640 | tails_map={USDT=0.00064640}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ICP | est=-0.8419% | ideal=-0.5438% | rounded=-0.5450% | strict=-0.8440% | tails=-0.8427% | tail_usdt=0.00131300 | tails_map={USDT=0.00131300}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→XLM | est=-0.5120% | ideal=-0.2129% | rounded=-0.2150% | strict=-0.5131% | tails=-0.5125% | tail_usdt=0.00058250 | tails_map={USDT=0.00058250}
2026/04/10 09:04:48 [EXEC CMP] USDT→ATOM→ETH | est=-0.5871% | ideal=-0.2882% | rounded=-0.2898% | strict=-0.5893% | tails=-0.5891% | tail_usdt=0.00020079 | tails_map={USDT=0.00020079}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ICP | est=-0.8418% | ideal=-0.5437% | rounded=-0.5450% | strict=-0.8440% | tails=-0.8425% | tail_usdt=0.00145200 | tails_map={USDT=0.00145200}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ATOM | est=-0.4999% | ideal=-0.2008% | rounded=-0.2016% | strict=-0.5010% | tails=-0.5003% | tail_usdt=0.00070940 | tails_map={USDT=0.00070940}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→FET | est=-0.4760% | ideal=-0.1768% | rounded=-0.1776% | strict=-0.4771% | tails=-0.4763% | tail_usdt=0.00078540 | tails_map={USDT=0.00078540}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→XLM | est=-0.5119% | ideal=-0.2128% | rounded=-0.2150% | strict=-0.5131% | tails=-0.5124% | tail_usdt=0.00072150 | tails_map={USDT=0.00072150}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→AVAX | est=-0.3675% | ideal=-0.0680% | rounded=-0.0700% | strict=-0.3710% | tails=-0.3700% | tail_usdt=0.00097900 | tails_map={USDT=0.00097900}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ATOM | est=-0.4999% | ideal=-0.2008% | rounded=-0.2016% | strict=-0.5010% | tails=-0.5003% | tail_usdt=0.00070940 | tails_map={USDT=0.00070940}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ATOM | est=-0.4944% | ideal=-0.1953% | rounded=-0.1961% | strict=-0.4956% | tails=-0.4948% | tail_usdt=0.00080390 | tails_map={USDT=0.00080390}
2026/04/10 09:04:48 [EXEC CMP] USDT→ATOM→ETH | est=-0.5871% | ideal=-0.2882% | rounded=-0.2898% | strict=-0.5893% | tails=-0.5891% | tail_usdt=0.00020079 | tails_map={USDT=0.00020079}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ICP | est=-0.8137% | ideal=-0.5156% | rounded=-0.5170% | strict=-0.8160% | tails=-0.8146% | tail_usdt=0.00142400 | tails_map={USDT=0.00142400}
2026/04/10 09:04:48 [EXEC CMP] USDT→BTC→ICP | est=-0.7857% | ideal=-0.4874% | rounded=-0.4890% | strict=-0.7880% | tails=-0.7866% | tail_usdt=0.00139600 | tails_map={USDT=0.00139600}
2026/04/10 09:04:49 [EXEC CMP] USDT→BTC→NEAR | est=-0.3779% | ideal=-0.0784% | rounded=-0.0792% | strict=-0.3790% | tails=-0.3782% | tail_usdt=0.00078720 | tails_map={USDT=0.00078720}
2026/04/10 09:04:53 [EXEC CMP] USDT→DOT→KCS | est=-0.8095% | ideal=-0.5114% | rounded=-0.5120% | strict=-0.8108% | tails=-0.8107% | tail_usdt=0.00011150 | tails_map={USDT=0.00011150}
2026/04/10 09:04:53 [STATS] ticks=9287 triangles_seen=13260 cand=161 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=584, no_quote_leg_3=1994, no_quote_leg_2=2905, no_quote_leg_1=7613} | exec_rejects={profit_below_threshold=161}
2026/04/10 09:04:53 [EXEC CMP] USDT→DOT→KCS | est=-0.7472% | ideal=-0.4488% | rounded=-0.4496% | strict=-0.7485% | tails=-0.7484% | tail_usdt=0.00014910 | tails_map={USDT=0.00014910}
2026/04/10 09:04:53 [REJECT] stage=scan reason=no_quote_leg_3 count=2000 tri=USDT->KCS->SOL
2026/04/10 09:04:53 [EXEC CMP] USDT→DOT→KCS | est=-0.7943% | ideal=-0.4961% | rounded=-0.4969% | strict=-0.7957% | tails=-0.7955% | tail_usdt=0.00019630 | tails_map={USDT=0.00019630}
2026/04/10 09:04:55 [EXEC CMP] USDT→DOT→KCS | est=-0.7666% | ideal=-0.4683% | rounded=-0.4690% | strict=-0.7679% | tails=-0.7677% | tail_usdt=0.00016850 | tails_map={USDT=0.00016850}
2026/04/10 09:04:55 [EXEC CMP] USDT→DOT→KCS | est=-0.7818% | ideal=-0.4835% | rounded=-0.4841% | strict=-0.7830% | tails=-0.7828% | tail_usdt=0.00018360 | tails_map={USDT=0.00018360}
2026/04/10 09:04:56 [EXEC CMP] USDT→DOT→KCS | est=-0.8095% | ideal=-0.5114% | rounded=-0.5120% | strict=-0.8108% | tails=-0.8107% | tail_usdt=0.00011150 | tails_map={USDT=0.00011150}
2026/04/10 09:04:56 [EXEC CMP] USDT→DOT→KCS | est=-0.8095% | ideal=-0.5114% | rounded=-0.5120% | strict=-0.8108% | tails=-0.8107% | tail_usdt=0.00011150 | tails_map={USDT=0.00011150}
2026/04/10 09:04:58 [STATS] ticks=9433 triangles_seen=13493 cand=167 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=595, no_quote_leg_3=2029, no_quote_leg_2=2959, no_quote_leg_1=7740} | exec_rejects={profit_below_threshold=167}
2026/04/10 09:04:59 [EXEC CMP] USDT→AVA→ETH | est=-1.2034% | ideal=-0.9065% | rounded=-0.9083% | strict=-1.2063% | tails=-1.2045% | tail_usdt=0.00180137 | tails_map={USDT=0.00180137}
2026/04/10 09:04:59 [EXEC CMP] USDT→BTC→FET | est=-0.4731% | ideal=-0.1739% | rounded=-0.1746% | strict=-0.4742% | tails=-0.4735% | tail_usdt=0.00070950 | tails_map={USDT=0.00070950}
2026/04/10 09:04:59 [EXEC CMP] USDT→BTC→AVAX | est=-0.3644% | ideal=-0.0649% | rounded=-0.0670% | strict=-0.3670% | tails=-0.3661% | tail_usdt=0.00090200 | tails_map={USDT=0.00090200}
2026/04/10 09:04:59 [REJECT] stage=scan reason=no_quote_leg_2 count=3000 tri=USDT->BTC->RSR
2026/04/10 09:04:59 [EXEC CMP] USDT→BTC→DOT | est=-0.3558% | ideal=-0.0563% | rounded=-0.0571% | strict=-0.3571% | tails=-0.3564% | tail_usdt=0.00069230 | tails_map={USDT=0.00069230}
2026/04/10 09:04:59 [EXEC CMP] USDT→LINK→BTC | est=-0.3946% | ideal=-0.0952% | rounded=-0.0953% | strict=-0.3963% | tails=-0.3962% | tail_usdt=0.00010061 | tails_map={USDT=0.00010061}
2026/04/10 09:04:59 [EXEC CMP] USDT→BTC→DOT | est=-0.3482% | ideal=-0.0486% | rounded=-0.0495% | strict=-0.3495% | tails=-0.3488% | tail_usdt=0.00068470 | tails_map={USDT=0.00068470}
2026/04/10 09:04:59 [EXEC CMP] USDT→BTC→AVAX | est=-0.3644% | ideal=-0.0649% | rounded=-0.0670% | strict=-0.3670% | tails=-0.3661% | tail_usdt=0.00090200 | tails_map={USDT=0.00090200}
2026/04/10 09:05:00 [EXEC CMP] USDT→BTC→ICP | est=-0.5350% | ideal=-0.2360% | rounded=-0.2370% | strict=-0.5370% | tails=-0.5364% | tail_usdt=0.00055000 | tails_map={USDT=0.00055000}
2026/04/10 09:05:00 [EXEC CMP] USDT→BTC→DOT | est=-0.3530% | ideal=-0.0535% | rounded=-0.0538% | strict=-0.3538% | tails=-0.3536% | tail_usdt=0.00016700 | tails_map={USDT=0.00016700}
2026/04/10 09:05:00 [EXEC CMP] USDT→BTC→FET | est=-0.4779% | ideal=-0.1788% | rounded=-0.1789% | strict=-0.4785% | tails=-0.4783% | tail_usdt=0.00019180 | tails_map={USDT=0.00019180}
2026/04/10 09:05:00 [EXEC CMP] USDT→BTC→FET | est=-0.4365% | ideal=-0.1372% | rounded=-0.1373% | strict=-0.4370% | tails=-0.4368% | tail_usdt=0.00015030 | tails_map={USDT=0.00015030}
2026/04/10 09:05:00 [EXEC CMP] USDT→XLM→ETH | est=-0.4521% | ideal=-0.1529% | rounded=-0.1542% | strict=-0.4536% | tails=-0.4524% | tail_usdt=0.00120000 | tails_map={USDT=0.00120000}
2026/04/10 09:05:00 [EXEC CMP] USDT→BTC→DOT | est=-0.3378% | ideal=-0.0382% | rounded=-0.0385% | strict=-0.3385% | tails=-0.3383% | tail_usdt=0.00015170 | tails_map={USDT=0.00015170}
2026/04/10 09:05:00 [EXEC CMP] USDT→BTC→DOT | est=-0.3927% | ideal=-0.0933% | rounded=-0.0935% | strict=-0.3934% | tails=-0.3932% | tail_usdt=0.00020660 | tails_map={USDT=0.00020660}
2026/04/10 09:05:00 [EXEC CMP] USDT→BTC→HBAR | est=-0.4302% | ideal=-0.1309% | rounded=-0.1310% | strict=-0.4306% | tails=-0.4305% | tail_usdt=0.00011394 | tails_map={USDT=0.00011394}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→HBAR | est=-0.4077% | ideal=-0.1084% | rounded=-0.1085% | strict=-0.4082% | tails=-0.4080% | tail_usdt=0.00012147 | tails_map={USDT=0.00012147}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→AVAX | est=-0.3621% | ideal=-0.0626% | rounded=-0.0630% | strict=-0.3650% | tails=-0.3646% | tail_usdt=0.00037800 | tails_map={USDT=0.00037800}
2026/04/10 09:05:01 [EXEC CMP] USDT→LINK→BTC | est=-0.3728% | ideal=-0.0733% | rounded=-0.0739% | strict=-0.3749% | tails=-0.3745% | tail_usdt=0.00040017 | tails_map={USDT=0.00040017}
2026/04/10 09:05:01 [EXEC CMP] USDT→LINK→BTC | est=-0.3750% | ideal=-0.0755% | rounded=-0.0767% | strict=-0.3771% | tails=-0.3762% | tail_usdt=0.00090034 | tails_map={USDT=0.00090034}
2026/04/10 09:05:01 [EXEC CMP] USDT→LINK→BTC | est=-0.3602% | ideal=-0.0607% | rounded=-0.0619% | strict=-0.3623% | tails=-0.3614% | tail_usdt=0.00090055 | tails_map={USDT=0.00090055}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→DOT | est=-0.4075% | ideal=-0.1081% | rounded=-0.1087% | strict=-0.4085% | tails=-0.4081% | tail_usdt=0.00044880 | tails_map={USDT=0.00044880}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→HBAR | est=-0.4225% | ideal=-0.1232% | rounded=-0.1236% | strict=-0.4233% | tails=-0.4228% | tail_usdt=0.00044359 | tails_map={USDT=0.00044359}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→FET | est=-0.4512% | ideal=-0.1520% | rounded=-0.1525% | strict=-0.4521% | tails=-0.4516% | tail_usdt=0.00049240 | tails_map={USDT=0.00049240}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→AVAX | est=-0.3769% | ideal=-0.0774% | rounded=-0.0780% | strict=-0.3800% | tails=-0.3793% | tail_usdt=0.00072000 | tails_map={USDT=0.00072000}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→FET | est=-0.4512% | ideal=-0.1520% | rounded=-0.1525% | strict=-0.4521% | tails=-0.4516% | tail_usdt=0.00049240 | tails_map={USDT=0.00049240}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→HBAR | est=-0.4113% | ideal=-0.1119% | rounded=-0.1124% | strict=-0.4120% | tails=-0.4116% | tail_usdt=0.00044236 | tails_map={USDT=0.00044236}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→AVAX | est=-0.3769% | ideal=-0.0774% | rounded=-0.0780% | strict=-0.3800% | tails=-0.3793% | tail_usdt=0.00072000 | tails_map={USDT=0.00072000}
2026/04/10 09:05:01 [EXEC CMP] USDT→LINK→BTC | est=-0.3624% | ideal=-0.0629% | rounded=-0.0634% | strict=-0.3645% | tails=-0.3641% | tail_usdt=0.00040072 | tails_map={USDT=0.00040072}
2026/04/10 09:05:01 [EXEC CMP] USDT→LINK→BTC | est=-0.3613% | ideal=-0.0618% | rounded=-0.0626% | strict=-0.3630% | tails=-0.3624% | tail_usdt=0.00060027 | tails_map={USDT=0.00060027}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→ADA | est=-0.3436% | ideal=-0.0440% | rounded=-0.0468% | strict=-0.3466% | tails=-0.3461% | tail_usdt=0.00048680 | tails_map={USDT=0.00048680}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→AVAX | est=-0.3769% | ideal=-0.0774% | rounded=-0.0780% | strict=-0.3800% | tails=-0.3793% | tail_usdt=0.00072000 | tails_map={USDT=0.00072000}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→AVAX | est=-0.3814% | ideal=-0.0820% | rounded=-0.0830% | strict=-0.3840% | tails=-0.3835% | tail_usdt=0.00048200 | tails_map={USDT=0.00048200}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→AVAX | est=-0.3708% | ideal=-0.0713% | rounded=-0.0730% | strict=-0.3740% | tails=-0.3735% | tail_usdt=0.00047200 | tails_map={USDT=0.00047200}
2026/04/10 09:05:01 [EXEC CMP] USDT→ATOM→ETH | est=-0.5378% | ideal=-0.2388% | rounded=-0.2404% | strict=-0.5400% | tails=-0.5398% | tail_usdt=0.00020042 | tails_map={USDT=0.00020042}
2026/04/10 09:05:01 [EXEC CMP] USDT→BTC→ATOM | est=-0.4672% | ideal=-0.1680% | rounded=-0.1687% | strict=-0.4683% | tails=-0.4678% | tail_usdt=0.00052760 | tails_map={USDT=0.00052760}
2026/04/10 09:05:02 [EXEC CMP] USDT→DOT→KCS | est=-0.6619% | ideal=-0.3633% | rounded=-0.3635% | strict=-0.6635% | tails=-0.6633% | tail_usdt=0.00016400 | tails_map={USDT=0.00016400}
2026/04/10 09:05:02 [EXEC CMP] USDT→SHIB→DOGE | est=-0.4488% | ideal=-0.1496% | rounded=-0.1496% | strict=-0.4488% | tails=-0.4488% | tail_usdt=0.00000918 | tails_map={USDT=0.00000918}
2026/04/10 09:05:02 [EXEC CMP] USDT→SHIB→DOGE | est=-0.4488% | ideal=-0.1496% | rounded=-0.1496% | strict=-0.4488% | tails=-0.4488% | tail_usdt=0.00000918 | tails_map={USDT=0.00000918}
2026/04/10 09:05:02 [EXEC CMP] USDT→BTC→DOT | est=-0.3497% | ideal=-0.0502% | rounded=-0.0507% | strict=-0.3508% | tails=-0.3503% | tail_usdt=0.00051000 | tails_map={USDT=0.00051000}
2026/04/10 09:05:02 [EXEC CMP] USDT→BTC→ICP | est=-0.7745% | ideal=-0.4762% | rounded=-0.4770% | strict=-0.7770% | tails=-0.7759% | tail_usdt=0.00113600 | tails_map={USDT=0.00113600}
2026/04/10 09:05:02 [EXEC CMP] USDT→LINK→BTC | est=-0.3891% | ideal=-0.0897% | rounded=-0.0908% | strict=-0.3912% | tails=-0.3906% | tail_usdt=0.00060047 | tails_map={USDT=0.00060047}
2026/04/10 09:05:02 [EXEC CMP] USDT→BTC→AVAX | est=-0.3652% | ideal=-0.0657% | rounded=-0.0680% | strict=-0.3680% | tails=-0.3673% | tail_usdt=0.00072700 | tails_map={USDT=0.00072700}
2026/04/10 09:05:02 [EXEC CMP] USDT→BTC→NEAR | est=-0.3780% | ideal=-0.0785% | rounded=-0.0791% | strict=-0.3791% | tails=-0.3786% | tail_usdt=0.00053830 | tails_map={USDT=0.00053830}
2026/04/10 09:05:02 [EXEC CMP] USDT→BTC→ICP | est=-0.8026% | ideal=-0.5044% | rounded=-0.5050% | strict=-0.8050% | tails=-0.8038% | tail_usdt=0.00116400 | tails_map={USDT=0.00116400}
2026/04/10 09:05:02 [EXEC CMP] USDT→BTC→NEAR | est=-0.3635% | ideal=-0.0640% | rounded=-0.0646% | strict=-0.3646% | tails=-0.3641% | tail_usdt=0.00052380 | tails_map={USDT=0.00052380}
2026/04/10 09:05:03 [STATS] ticks=9884 triangles_seen=14294 cand=212 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=746, no_quote_leg_3=2144, no_quote_leg_2=3230, no_quote_leg_1=7959} | exec_rejects={profit_below_threshold=212}
2026/04/10 09:05:04 [REJECT] stage=scan reason=no_quote_leg_1 count=8000 tri=USDT->LYX->ETH
2026/04/10 09:05:08 [EXEC CMP] USDT→LINK→BTC | est=-0.3747% | ideal=-0.0752% | rounded=-0.0759% | strict=-0.3770% | tails=-0.3764% | tail_usdt=0.00060025 | tails_map={USDT=0.00060025}
2026/04/10 09:05:08 [STATS] ticks=10141 triangles_seen=14668 cand=213 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=766, no_quote_leg_3=2208, no_quote_leg_2=3315, no_quote_leg_1=8163} | exec_rejects={profit_below_threshold=213}
2026/04/10 09:05:08 [EXEC CMP] USDT→SHIB→DOGE | est=-0.4330% | ideal=-0.1337% | rounded=-0.1337% | strict=-0.4330% | tails=-0.4330% | tail_usdt=0.00000335 | tails_map={USDT=0.00000335}
2026/04/10 09:05:08 [EXEC CMP] USDT→SHIB→DOGE | est=-0.4163% | ideal=-0.1169% | rounded=-0.1169% | strict=-0.4163% | tails=-0.4163% | tail_usdt=0.00000658 | tails_map={USDT=0.00000658}
2026/04/10 09:05:13 [STATS] ticks=10279 triangles_seen=14842 cand=215 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=771, no_quote_leg_3=2240, no_quote_leg_2=3340, no_quote_leg_1=8273} | exec_rejects={profit_below_threshold=215}
2026/04/10 09:05:13 [EXEC CMP] USDT→DOT→KCS | est=-0.6628% | ideal=-0.3642% | rounded=-0.3644% | strict=-0.6644% | tails=-0.6642% | tail_usdt=0.00016490 | tails_map={USDT=0.00016490}
2026/04/10 09:05:16 [EXEC CMP] USDT→DOT→KCS | est=-0.5757% | ideal=-0.2768% | rounded=-0.2777% | strict=-0.5769% | tails=-0.5767% | tail_usdt=0.00017730 | tails_map={USDT=0.00017730}
2026/04/10 09:05:16 [EXEC CMP] USDT→DOT→KCS | est=-0.5757% | ideal=-0.2768% | rounded=-0.2777% | strict=-0.5769% | tails=-0.5767% | tail_usdt=0.00017730 | tails_map={USDT=0.00017730}
2026/04/10 09:05:16 [EXEC CMP] USDT→DOT→KCS | est=-0.5757% | ideal=-0.2768% | rounded=-0.2777% | strict=-0.5769% | tails=-0.5767% | tail_usdt=0.00017730 | tails_map={USDT=0.00017730}
2026/04/10 09:05:16 [EXEC CMP] USDT→BTC→ICP | est=-0.6025% | ideal=-0.3037% | rounded=-0.3050% | strict=-0.6050% | tails=-0.6038% | tail_usdt=0.00120100 | tails_map={USDT=0.00120100}
2026/04/10 09:05:16 [EXEC CMP] USDT→BTC→FET | est=-0.5726% | ideal=-0.2737% | rounded=-0.2745% | strict=-0.5738% | tails=-0.5730% | tail_usdt=0.00077020 | tails_map={USDT=0.00077020}
2026/04/10 09:05:16 [EXEC CMP] USDT→LINK→BTC | est=-0.3722% | ideal=-0.0728% | rounded=-0.0737% | strict=-0.3740% | tails=-0.3734% | tail_usdt=0.00060131 | tails_map={USDT=0.00060131}
2026/04/10 09:05:16 [EXEC CMP] USDT→BTC→ICP | est=-0.6025% | ideal=-0.3037% | rounded=-0.3050% | strict=-0.6050% | tails=-0.6038% | tail_usdt=0.00120100 | tails_map={USDT=0.00120100}
2026/04/10 09:05:18 [STATS] ticks=10489 triangles_seen=15160 cand=223 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=787, no_quote_leg_3=2280, no_quote_leg_2=3430, no_quote_leg_1=8437} | exec_rejects={profit_below_threshold=223}
2026/04/10 09:05:19 [EXEC CMP] USDT→LINK→BTC | est=-0.3688% | ideal=-0.0693% | rounded=-0.0701% | strict=-0.3704% | tails=-0.3699% | tail_usdt=0.00050066 | tails_map={USDT=0.00050066}
2026/04/10 09:05:19 [EXEC CMP] USDT→LINK→BTC | est=-0.3688% | ideal=-0.0693% | rounded=-0.0701% | strict=-0.3704% | tails=-0.3699% | tail_usdt=0.00050066 | tails_map={USDT=0.00050066}
2026/04/10 09:05:23 [STATS] ticks=10637 triangles_seen=15414 cand=225 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=811, no_quote_leg_3=2307, no_quote_leg_2=3495, no_quote_leg_1=8573} | exec_rejects={profit_below_threshold=225}
2026/04/10 09:05:28 [STATS] ticks=10752 triangles_seen=15542 cand=225 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=811, no_quote_leg_3=2332, no_quote_leg_2=3514, no_quote_leg_1=8657} | exec_rejects={profit_below_threshold=225}
2026/04/10 09:05:30 [EXEC CMP] USDT→LINK→BTC | est=-0.3652% | ideal=-0.0657% | rounded=-0.0662% | strict=-0.3665% | tails=-0.3664% | tail_usdt=0.00010076 | tails_map={USDT=0.00010076}
2026/04/10 09:05:30 [EXEC CMP] USDT→BTC→AVAX | est=-0.3936% | ideal=-0.0942% | rounded=-0.0950% | strict=-0.3960% | tails=-0.3951% | tail_usdt=0.00088200 | tails_map={USDT=0.00088200}
2026/04/10 09:05:30 [EXEC CMP] USDT→BTC→AVAX | est=-0.3830% | ideal=-0.0836% | rounded=-0.0850% | strict=-0.3860% | tails=-0.3851% | tail_usdt=0.00087200 | tails_map={USDT=0.00087200}
2026/04/10 09:05:30 [EXEC CMP] USDT→BTC→NEAR | est=-0.3591% | ideal=-0.0596% | rounded=-0.0604% | strict=-0.3602% | tails=-0.3596% | tail_usdt=0.00064640 | tails_map={USDT=0.00064640}
2026/04/10 09:05:30 [EXEC CMP] USDT→BTC→AVAX | est=-0.3843% | ideal=-0.0848% | rounded=-0.0860% | strict=-0.3870% | tails=-0.3867% | tail_usdt=0.00034800 | tails_map={USDT=0.00034800}
2026/04/10 09:05:30 [EXEC CMP] USDT→BTC→NEAR | est=-0.3604% | ideal=-0.0609% | rounded=-0.0610% | strict=-0.3610% | tails=-0.3609% | tail_usdt=0.00012220 | tails_map={USDT=0.00012220}
2026/04/10 09:05:30 [EXEC CMP] USDT→DOT→KCS | est=-0.6780% | ideal=-0.3794% | rounded=-0.3796% | strict=-0.6795% | tails=-0.6792% | tail_usdt=0.00028000 | tails_map={USDT=0.00028000}
2026/04/10 09:05:30 [EXEC CMP] USDT→SHIB→DOGE | est=-0.5965% | ideal=-0.2977% | rounded=-0.2977% | strict=-0.5965% | tails=-0.5965% | tail_usdt=0.00000702 | tails_map={USDT=0.00000702}
2026/04/10 09:05:30 [EXEC CMP] USDT→SHIB→DOGE | est=-0.5965% | ideal=-0.2977% | rounded=-0.2977% | strict=-0.5965% | tails=-0.5965% | tail_usdt=0.00000702 | tails_map={USDT=0.00000702}
2026/04/10 09:05:30 [EXEC CMP] USDT→BTC→DOT | est=-0.4026% | ideal=-0.1032% | rounded=-0.1034% | strict=-0.4031% | tails=-0.4030% | tail_usdt=0.00006440 | tails_map={USDT=0.00006440}
2026/04/10 09:05:30 [EXEC CMP] USDT→SHIB→DOGE | est=-0.5965% | ideal=-0.2977% | rounded=-0.2977% | strict=-0.5965% | tails=-0.5965% | tail_usdt=0.00000702 | tails_map={USDT=0.00000702}
2026/04/10 09:05:31 [EXEC CMP] USDT→LINK→BTC | est=-0.3509% | ideal=-0.0514% | rounded=-0.0522% | strict=-0.3525% | tails=-0.3521% | tail_usdt=0.00040179 | tails_map={USDT=0.00040179}
2026/04/10 09:05:31 [EXEC CMP] USDT→BTC→NEAR | est=-0.4126% | ideal=-0.1132% | rounded=-0.1134% | strict=-0.4131% | tails=-0.4130% | tail_usdt=0.00008440 | tails_map={USDT=0.00008440}
2026/04/10 09:05:31 [EXEC CMP] USDT→LINK→BTC | est=-0.3565% | ideal=-0.0569% | rounded=-0.0572% | strict=-0.3583% | tails=-0.3581% | tail_usdt=0.00020156 | tails_map={USDT=0.00020156}
2026/04/10 09:05:31 [EXEC CMP] USDT→LINK→BTC | est=-0.3565% | ideal=-0.0569% | rounded=-0.0572% | strict=-0.3583% | tails=-0.3581% | tail_usdt=0.00020156 | tails_map={USDT=0.00020156}
2026/04/10 09:05:31 [EXEC CMP] USDT→BTC→NEAR | est=-0.4126% | ideal=-0.1132% | rounded=-0.1134% | strict=-0.4131% | tails=-0.4130% | tail_usdt=0.00008440 | tails_map={USDT=0.00008440}
2026/04/10 09:05:31 [EXEC CMP] USDT→LINK→BTC | est=-0.3565% | ideal=-0.0569% | rounded=-0.0572% | strict=-0.3583% | tails=-0.3581% | tail_usdt=0.00020156 | tails_map={USDT=0.00020156}
2026/04/10 09:05:33 [STATS] ticks=10987 triangles_seen=16033 cand=242 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=857, no_quote_leg_3=2394, no_quote_leg_2=3714, no_quote_leg_1=8823} | exec_rejects={profit_below_threshold=242}
2026/04/10 09:05:38 [STATS] ticks=11166 triangles_seen=16225 cand=242 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=859, no_quote_leg_3=2417, no_quote_leg_2=3724, no_quote_leg_1=8980} | exec_rejects={profit_below_threshold=242}
2026/04/10 09:05:38 [REJECT] stage=scan reason=no_quote_leg_1 count=9000 tri=USDT->XRP->KCS
2026/04/10 09:05:42 [EXEC CMP] USDT→BTC→FET | est=-0.5371% | ideal=-0.2381% | rounded=-0.2388% | strict=-0.5381% | tails=-0.5375% | tail_usdt=0.00064450 | tails_map={USDT=0.00064450}
2026/04/10 09:05:42 [EXEC CMP] USDT→BTC→FET | est=-0.5350% | ideal=-0.2361% | rounded=-0.2366% | strict=-0.5360% | tails=-0.5354% | tail_usdt=0.00055740 | tails_map={USDT=0.00055740}
2026/04/10 09:05:42 [EXEC CMP] USDT→BTC→FET | est=-0.5765% | ideal=-0.2777% | rounded=-0.2782% | strict=-0.5775% | tails=-0.5769% | tail_usdt=0.00059890 | tails_map={USDT=0.00059890}
2026/04/10 09:05:43 [STATS] ticks=11375 triangles_seen=16578 cand=245 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=881, no_quote_leg_3=2466, no_quote_leg_2=3827, no_quote_leg_1=9156} | exec_rejects={profit_below_threshold=245}
2026/04/10 09:05:44 [EXEC CMP] USDT→LINK→BTC | est=-0.3730% | ideal=-0.0735% | rounded=-0.0739% | strict=-0.3742% | tails=-0.3740% | tail_usdt=0.00020044 | tails_map={USDT=0.00020044}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→FET | est=-0.6156% | ideal=-0.3168% | rounded=-0.3170% | strict=-0.6161% | tails=-0.6160% | tail_usdt=0.00014560 | tails_map={USDT=0.00014560}
2026/04/10 09:05:44 [EXEC CMP] USDT→LINK→BTC | est=-0.3650% | ideal=-0.0654% | rounded=-0.0659% | strict=-0.3662% | tails=-0.3660% | tail_usdt=0.00020050 | tails_map={USDT=0.00020050}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→FET | est=-0.5741% | ideal=-0.2752% | rounded=-0.2754% | strict=-0.5746% | tails=-0.5744% | tail_usdt=0.00020400 | tails_map={USDT=0.00020400}
2026/04/10 09:05:44 [EXEC CMP] USDT→LINK→BTC | est=-0.3728% | ideal=-0.0733% | rounded=-0.0738% | strict=-0.3741% | tails=-0.3739% | tail_usdt=0.00020038 | tails_map={USDT=0.00020038}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→FET | est=-0.5662% | ideal=-0.2673% | rounded=-0.2675% | strict=-0.5667% | tails=-0.5665% | tail_usdt=0.00016010 | tails_map={USDT=0.00016010}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→HBAR | est=-0.4459% | ideal=-0.1467% | rounded=-0.1468% | strict=-0.4463% | tails=-0.4462% | tail_usdt=0.00009965 | tails_map={USDT=0.00009965}
2026/04/10 09:05:44 [EXEC CMP] USDT→LINK→BTC | est=-0.3695% | ideal=-0.0700% | rounded=-0.0709% | strict=-0.3719% | tails=-0.3710% | tail_usdt=0.00090122 | tails_map={USDT=0.00090122}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→DOT | est=-0.3945% | ideal=-0.0951% | rounded=-0.0953% | strict=-0.3952% | tails=-0.3950% | tail_usdt=0.00018840 | tails_map={USDT=0.00018840}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→AVAX | est=-0.3772% | ideal=-0.0777% | rounded=-0.0780% | strict=-0.3790% | tails=-0.3786% | tail_usdt=0.00037200 | tails_map={USDT=0.00037200}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→DASH | est=-0.4698% | ideal=-0.1706% | rounded=-0.1800% | strict=-0.4800% | tails=-0.4795% | tail_usdt=0.00047300 | tails_map={USDT=0.00047300}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→DOT | est=-0.3869% | ideal=-0.0874% | rounded=-0.0876% | strict=-0.3875% | tails=-0.3873% | tail_usdt=0.00018070 | tails_map={USDT=0.00018070}
2026/04/10 09:05:44 [EXEC CMP] USDT→BTC→HBAR | est=-0.4459% | ideal=-0.1467% | rounded=-0.1468% | strict=-0.4463% | tails=-0.4462% | tail_usdt=0.00009965 | tails_map={USDT=0.00009965}
2026/04/10 09:05:45 [EXEC CMP] USDT→BTC→AVAX | est=-0.3772% | ideal=-0.0777% | rounded=-0.0780% | strict=-0.3790% | tails=-0.3786% | tail_usdt=0.00037200 | tails_map={USDT=0.00037200}
2026/04/10 09:05:45 [EXEC CMP] USDT→BTC→AVAX | est=-0.3772% | ideal=-0.0777% | rounded=-0.0780% | strict=-0.3790% | tails=-0.3786% | tail_usdt=0.00037200 | tails_map={USDT=0.00037200}
2026/04/10 09:05:45 [EXEC CMP] USDT→BTC→DASH | est=-0.4957% | ideal=-0.1966% | rounded=-0.2100% | strict=-0.5000% | tails=-0.4995% | tail_usdt=0.00049300 | tails_map={USDT=0.00049300}
2026/04/10 09:05:45 [EXEC CMP] USDT→LINK→BTC | est=-0.3695% | ideal=-0.0700% | rounded=-0.0709% | strict=-0.3719% | tails=-0.3710% | tail_usdt=0.00090122 | tails_map={USDT=0.00090122}
2026/04/10 09:05:45 [EXEC CMP] USDT→BTC→DOT | est=-0.4021% | ideal=-0.1027% | rounded=-0.1029% | strict=-0.4027% | tails=-0.4026% | tail_usdt=0.00009600 | tails_map={USDT=0.00009600}
2026/04/10 09:05:45 [EXEC CMP] USDT→BTC→DASH | est=-0.4698% | ideal=-0.1706% | rounded=-0.1800% | strict=-0.4800% | tails=-0.4795% | tail_usdt=0.00047300 | tails_map={USDT=0.00047300}
2026/04/10 09:05:45 [EXEC CMP] USDT→LINK→BTC | est=-0.3695% | ideal=-0.0700% | rounded=-0.0709% | strict=-0.3719% | tails=-0.3710% | tail_usdt=0.00090122 | tails_map={USDT=0.00090122}
2026/04/10 09:05:45 [EXEC CMP] USDT→LINK→BTC | est=-0.3684% | ideal=-0.0689% | rounded=-0.0695% | strict=-0.3698% | tails=-0.3696% | tail_usdt=0.00020005 | tails_map={USDT=0.00020005}
2026/04/10 09:05:47 [EXEC CMP] USDT→LINK→BTC | est=-0.3776% | ideal=-0.0782% | rounded=-0.0792% | strict=-0.3794% | tails=-0.3790% | tail_usdt=0.00040170 | tails_map={USDT=0.00040170}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→ATOM | est=-0.4243% | ideal=-0.1250% | rounded=-0.1255% | strict=-0.4253% | tails=-0.4249% | tail_usdt=0.00041660 | tails_map={USDT=0.00041660}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→AVAX | est=-0.3774% | ideal=-0.0780% | rounded=-0.0800% | strict=-0.3800% | tails=-0.3793% | tail_usdt=0.00067100 | tails_map={USDT=0.00067100}
2026/04/10 09:05:47 [EXEC CMP] USDT→LINK→BTC | est=-0.3665% | ideal=-0.0670% | rounded=-0.0683% | strict=-0.3686% | tails=-0.3678% | tail_usdt=0.00080088 | tails_map={USDT=0.00080088}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→AVAX | est=-0.3774% | ideal=-0.0780% | rounded=-0.0800% | strict=-0.3800% | tails=-0.3793% | tail_usdt=0.00067100 | tails_map={USDT=0.00067100}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→FET | est=-0.5974% | ideal=-0.2986% | rounded=-0.2990% | strict=-0.5982% | tails=-0.5977% | tail_usdt=0.00048960 | tails_map={USDT=0.00048960}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→ATOM | est=-0.4243% | ideal=-0.1250% | rounded=-0.1255% | strict=-0.4253% | tails=-0.4249% | tail_usdt=0.00041660 | tails_map={USDT=0.00041660}
2026/04/10 09:05:47 [EXEC CMP] USDT→LINK→BTC | est=-0.3665% | ideal=-0.0670% | rounded=-0.0683% | strict=-0.3686% | tails=-0.3678% | tail_usdt=0.00080088 | tails_map={USDT=0.00080088}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→HBAR | est=-0.4468% | ideal=-0.1475% | rounded=-0.1479% | strict=-0.4475% | tails=-0.4471% | tail_usdt=0.00039882 | tails_map={USDT=0.00039882}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→ATOM | est=-0.4460% | ideal=-0.1468% | rounded=-0.1473% | strict=-0.4470% | tails=-0.4466% | tail_usdt=0.00043830 | tails_map={USDT=0.00043830}
2026/04/10 09:05:47 [EXEC CMP] USDT→ATOM→ETH | est=-0.5487% | ideal=-0.2498% | rounded=-0.2508% | strict=-0.5525% | tails=-0.5523% | tail_usdt=0.00020092 | tails_map={USDT=0.00020092}
2026/04/10 09:05:47 [EXEC CMP] USDT→LINK→BTC | est=-0.3665% | ideal=-0.0670% | rounded=-0.0683% | strict=-0.3686% | tails=-0.3678% | tail_usdt=0.00080088 | tails_map={USDT=0.00080088}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→ATOM | est=-0.4352% | ideal=-0.1359% | rounded=-0.1364% | strict=-0.4361% | tails=-0.4357% | tail_usdt=0.00042740 | tails_map={USDT=0.00042740}
2026/04/10 09:05:47 [EXEC CMP] USDT→ATOM→ETH | est=-0.5487% | ideal=-0.2498% | rounded=-0.2508% | strict=-0.5525% | tails=-0.5523% | tail_usdt=0.00020092 | tails_map={USDT=0.00020092}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→HBAR | est=-0.4468% | ideal=-0.1475% | rounded=-0.1479% | strict=-0.4475% | tails=-0.4471% | tail_usdt=0.00039882 | tails_map={USDT=0.00039882}
2026/04/10 09:05:47 [EXEC CMP] USDT→LINK→BTC | est=-0.3665% | ideal=-0.0670% | rounded=-0.0683% | strict=-0.3686% | tails=-0.3678% | tail_usdt=0.00080088 | tails_map={USDT=0.00080088}
2026/04/10 09:05:47 [EXEC CMP] USDT→BTC→ATOM | est=-0.4460% | ideal=-0.1468% | rounded=-0.1473% | strict=-0.4470% | tails=-0.4466% | tail_usdt=0.00043830 | tails_map={USDT=0.00043830}
2026/04/10 09:05:47 [EXEC CMP] USDT→ATOM→ETH | est=-0.5487% | ideal=-0.2498% | rounded=-0.2508% | strict=-0.5525% | tails=-0.5523% | tail_usdt=0.00020092 | tails_map={USDT=0.00020092}
2026/04/10 09:05:48 [STATS] ticks=11636 triangles_seen=17037 cand=284 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=916, no_quote_leg_3=2550, no_quote_leg_2=3957, no_quote_leg_1=9327} | exec_rejects={profit_below_threshold=284}
2026/04/10 09:05:53 [STATS] ticks=11838 triangles_seen=17277 cand=284 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=928, no_quote_leg_3=2587, no_quote_leg_2=3976, no_quote_leg_1=9499} | exec_rejects={profit_below_threshold=284}
2026/04/10 09:05:58 [STATS] ticks=11978 triangles_seen=17447 cand=284 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=944, no_quote_leg_3=2614, no_quote_leg_2=3996, no_quote_leg_1=9606} | exec_rejects={profit_below_threshold=284}
2026/04/10 09:05:59 [REJECT] stage=scan reason=no_quote_leg_2 count=4000 tri=USDT->KCS->SOL
2026/04/10 09:06:03 [STATS] ticks=12230 triangles_seen=17745 cand=284 exec=0 | scan_rejects={max_start_zero=3, max_start_lt_100=959, no_quote_leg_3=2667, no_quote_leg_2=4019, no_quote_leg_1=9813} | exec_rejects={profit_below_threshold=284}
^C2026/04/10 09:06:06 [Main] shutting down...
