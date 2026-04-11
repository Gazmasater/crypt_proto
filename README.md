gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb/metrics$ go run .
=== BASIC ===
rows: 953413
duration_sec: 1555.35
events_per_min: 36779.39
unique_triangles: 302

=== PROFIT ===

profit_pct:
count=953413 min=-0.727733 max=0.000280 mean=-0.010907 p50=-0.006813 p90=-0.003158 p95=-0.002925 p99=-0.002526

profit_usdt:
count=953413 min=-680.458682 max=0.021974 mean=-0.982330 p50=-0.295460 p90=-0.000018 p95=-0.000004 p99=-0.000001

volume_usdt:
count=953413 min=0.000010 max=12355.685002 mean=162.214420 p50=46.569367 p90=296.471233 p95=702.027507 p99=2030.729820

opportunity_strength:
count=953413 min=-0.872632 max=0.000579 mean=-0.001097 p50=-0.000006 p90=-0.000000 p95=-0.000000 p99=-0.000000

share(profit_pct > 0): 0.0000
share(profit_pct > 0.001): 0.0000
share(profit_pct > 0.002): 0.0000

=== AGE ===

age_min_ms:
count=953413 min=0.000000 max=44.000000 mean=0.855303 p50=0.000000 p90=2.000000 p95=3.000000 p99=6.000000

age_max_ms:
count=953413 min=1.000000 max=1550628.000000 mean=48657.798824 p50=6659.000000 p90=99356.000000 p95=193251.000000 p99=946748.640000

age_spread_ms:
count=953413 min=0.000000 max=1550627.000000 mean=48656.943521 p50=6658.000000 p90=99353.800000 p95=193248.400000 p99=946747.640000

share(age_max_ms <= 180): 0.0530
share(age_max_ms <= 360): 0.0962
share(age_max_ms > 360): 0.9038

share(age_spread_ms <= 50): 0.0083
share(age_spread_ms <= 180): 0.0532
share(age_spread_ms > 180): 0.9468

=== TOP TRIANGLES BY COUNT ===
 1. BTC->ETH->EUR                    19647
 2. BTC->EUR->ETH                    19647
 3. USDT->ETH->EUR                   12648
 4. USDT->EUR->ETH                   12642
 5. USDT->BTC->EUR                   10762
 6. USDT->EUR->BTC                   10762
 7. USDT->BTC->TEL                   9813
 8. USDT->TEL->BTC                   9813
 9. USDT->ETH->TEL                   9450
10. USDT->TEL->ETH                   9447
11. USDT->XMR->BTC                   7637
12. USDT->BTC->XMR                   7625
13. USDT->XMR->ETH                   7234
14. USDT->ETH->XMR                   7200
15. BTC->TEL->ETH                    7105
16. BTC->ETH->TEL                    7103
17. USDT->BTC->ZEC                   6869
18. USDT->ZEC->BTC                   6745
19. USDT->BTC->DASH                  6685
20. USDT->AAVE->BTC                  6402

=== TOP TRIANGLES BY MEAN PROFIT ===
 1. USDT->VSYS->BTC                  count=10 mean_profit_pct=-0.002045 median=-0.002047 max=-0.002043 mean_vol=2.20 mean_age_max=1590.50 mean_age_spread=1590.50
 2. USDT->BTC->ERG                   count=668 mean_profit_pct=-0.002469 median=-0.002330 max=-0.001497 mean_vol=1.66 mean_age_max=130746.73 mean_age_spread=130745.98
 3. USDT->ETH->BTC                   count=4912 mean_profit_pct=-0.002501 median=-0.002502 max=-0.002170 mean_vol=1172.86 mean_age_max=7202.24 mean_age_spread=7200.78
 4. USDT->BTC->ETH                   count=4906 mean_profit_pct=-0.002637 median=-0.002631 max=-0.002211 mean_vol=1178.74 mean_age_max=7209.32 mean_age_spread=7208.66
 5. USDT->ETH->XRP                   count=4597 mean_profit_pct=-0.002723 median=-0.002723 max=-0.002424 mean_vol=86.98 mean_age_max=9551.41 mean_age_spread=9550.87
 6. USDT->XRP->ETH                   count=4590 mean_profit_pct=-0.002738 median=-0.002746 max=-0.002301 mean_vol=119.16 mean_age_max=9566.05 mean_age_spread=9565.26
 7. USDT->DOGE->BTC                  count=4031 mean_profit_pct=-0.002763 median=-0.002788 max=-0.001849 mean_vol=104.51 mean_age_max=9562.02 mean_age_spread=9560.54
 8. USDT->XRP->BTC                   count=5095 mean_profit_pct=-0.002834 median=-0.002858 max=-0.002349 mean_vol=1177.16 mean_age_max=13409.29 mean_age_spread=13407.51
 9. USDT->BTC->XRP                   count=5095 mean_profit_pct=-0.002904 median=-0.002908 max=-0.002263 mean_vol=320.72 mean_age_max=13408.57 mean_age_spread=13407.51
10. USDT->BTC->KCS                   count=2843 mean_profit_pct=-0.002949 median=-0.002885 max=-0.002265 mean_vol=52.52 mean_age_max=35315.18 mean_age_spread=35314.03
11. USDT->KCS->DOGE                  count=1560 mean_profit_pct=-0.002985 median=-0.002977 max=-0.002328 mean_vol=35.18 mean_age_max=19281.54 mean_age_spread=19281.22
12. USDT->TRX->ETH                   count=9 mean_profit_pct=-0.002986 median=-0.002986 max=-0.002986 mean_vol=44.99 mean_age_max=5187.78 mean_age_spread=5187.11
13. USDT->DOGE->KCS                  count=1413 mean_profit_pct=-0.002987 median=-0.002963 max=-0.002406 mean_vol=29.80 mean_age_max=18709.12 mean_age_spread=18708.88
14. USDT->TRX->BTC                   count=3789 mean_profit_pct=-0.003023 median=-0.003016 max=-0.002494 mean_vol=104.95 mean_age_max=9928.33 mean_age_spread=9926.48
15. USDT->BTC->TRX                   count=4305 mean_profit_pct=-0.003065 median=-0.003081 max=-0.002646 mean_vol=83.95 mean_age_max=10411.20 mean_age_spread=10410.12
16. USDT->LINK->BTC                  count=4535 mean_profit_pct=-0.003073 median=-0.003075 max=-0.002568 mean_vol=229.95 mean_age_max=4303.44 mean_age_spread=4301.87
17. USDT->BTC->LINK                  count=4537 mean_profit_pct=-0.003075 median=-0.003073 max=-0.002465 mean_vol=282.53 mean_age_max=4301.01 mean_age_spread=4300.18
18. USDT->WAVES->BTC                 count=2712 mean_profit_pct=-0.003098 median=-0.003203 max=-0.002415 mean_vol=3.46 mean_age_max=279588.20 mean_age_spread=279585.22
19. USDT->DOT->BTC                   count=4673 mean_profit_pct=-0.003133 median=-0.003099 max=-0.002231 mean_vol=136.50 mean_age_max=3665.23 mean_age_spread=3663.94
20. USDT->KCS->BTC                   count=2843 mean_profit_pct=-0.003153 median=-0.003205 max=-0.002651 mean_vol=72.85 mean_age_max=35316.38 mean_age_spread=35314.03

=== TOP ASSETS ===

B:
 1. BTC                  263980
 2. ETH                  190557
 3. EUR                  43051
 4. KCS                  31583
 5. TEL                  26365
 6. XMR                  17791
 7. DASH                 15455
 8. XRP                  14494
 9. LTC                  12242
10. ONT                  11230
11. TRAC                 10279
12. TRX                  9669
13. XLM                  9318
14. ATOM                 9175
15. ETC                  9094
16. FET                  8320
17. DOGE                 8193
18. DOT                  8109
19. BNB                  7577
20. ADA                  7233

C:
 1. BTC                  263592
 2. ETH                  179816
 3. EUR                  43057
 4. KCS                  28669
 5. TEL                  26366
 6. XMR                  17739
 7. DASH                 16065
 8. TRX                  14757
 9. XRP                  14501
10. LTC                  13529
11. ONT                  10774
12. TRAC                 10278
13. XLM                  9356
14. ATOM                 9291
15. ETC                  9101
16. DOGE                 8328
17. DOT                  8117
18. BNB                  7937
19. XDC                  7845
20. A                    7409

leg1_symbol:
 1. BTC-USDT             263980
 2. ETH-USDT             129124
 3. ETH-BTC              61433
 4. KCS-USDT             25723
 5. USDT-EUR             23404
 6. BTC-EUR              19647
 7. TEL-USDT             19260
 8. XMR-USDT             14871
 9. XRP-USDT             12375
10. DASH-USDT            11515
11. LTC-USDT             10854
12. ONT-USDT             9189
13. TRX-USDT             8997
14. TRAC-USDT            8460
15. ATOM-USDT            7657
16. XLM-USDT             7656
17. DOGE-USDT            7598
18. ETC-USDT             7178
19. TEL-BTC              7105
20. FET-USDT             7098

leg2_symbol:
 1. ETH-EUR              64584
 2. TEL-ETH              33105
 3. BTC-EUR              21524
 4. XMR-ETH              20268
 5. TEL-BTC              19626
 6. DASH-ETH             19267
 7. XMR-BTC              15262
 8. ZEC-BTC              13614
 9. ONT-ETH              13521
10. AAVE-BTC             12798
11. DASH-BTC             12253
12. XRP-ETH              12189
13. TRAC-ETH             11515
14. LTC-ETH              11154
15. LYX-ETH              11101
16. CRO-BTC              10634
17. ETC-ETH              10419
18. DAG-ETH              10377
19. XLM-ETH              10286
20. XRP-BTC              10190

leg3_symbol:
 1. BTC-USDT             263592
 2. ETH-USDT             121082
 3. ETH-BTC              58734
 4. USDT-EUR             23410
 5. KCS-USDT             23108
 6. BTC-EUR              19647
 7. TEL-USDT             19263
 8. XMR-USDT             14825
 9. TRX-USDT             12704
10. DASH-USDT            12675
11. XRP-USDT             12382
12. LTC-USDT             11975
13. TRAC-USDT            8462
14. ONT-USDT             8065
15. DOGE-USDT            7733
16. ATOM-USDT            7662
17. XLM-USDT             7662
18. ETC-USDT             7187
19. TEL-BTC              7103
20. ZEC-USDT             6869

=== CLEAN SUBSET ===
clean_rows: 1
clean_ratio: 0.0000

reports saved to ./arb_reports
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_pr

