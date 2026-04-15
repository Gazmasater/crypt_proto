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






2026/04/15 15:03:12 [KuCoin] started with 2 WS
2026/04/15 15:03:12 [Main] KuCoinCollector started
2026/04/15 15:03:12 [Calculator] indexed 206 symbols
2026/04/15 15:03:13 [KuCoin WS 0] connected
2026/04/15 15:03:13 [KuCoin WS 1] connected
2026/04/15 15:03:23 [KuCoin WS 1] bootstrap progress 10/80 last=RSR-USDT
2026/04/15 15:03:23 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/15 15:03:24 [KuCoin WS 1] bootstrap progress 20/80 last=SNX-ETH
2026/04/15 15:03:24 [KuCoin WS 1] bootstrap progress 30/80 last=TEL-BTC
2026/04/15 15:03:24 [KuCoin WS 1] bootstrap progress 40/80 last=TRX-ETH
2026/04/15 15:03:25 [KuCoin WS 1] bootstrap progress 50/80 last=VSYS-USDT
2026/04/15 15:03:25 [KuCoin WS 1] bootstrap progress 60/80 last=XDC-USDT
2026/04/15 15:03:25 [KuCoin WS 1] bootstrap progress 70/80 last=XRP-KCS
2026/04/15 15:03:26 [KuCoin WS 1] bootstrap progress 80/80 last=ZIL-ETH
2026/04/15 15:03:26 [KuCoin WS 1] bootstrap finished 80/80 in 2.734650349s
2026/04/15 15:03:29 [KuCoin WS 0] bootstrap progress 10/126 last=ALGO-ETH
2026/04/15 15:03:29 [KuCoin WS 0] bootstrap progress 20/126 last=AVA-ETH
2026/04/15 15:03:29 [KuCoin WS 0] bootstrap progress 30/126 last=BDX-USDT
2026/04/15 15:03:30 [KuCoin WS 0] bootstrap progress 40/126 last=CKB-USDT
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.46357360
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26749593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26749593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26749593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.44111962
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26749593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24672595
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24672595
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.80074467
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24672595
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.80074467
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24939595
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26952593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26952593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26952593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.23856596
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24799595
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26860593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26860593
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.23764596
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24732595
2026/04/15 15:03:30 [KuCoin WS 0] bootstrap progress 50/126 last=DASH-ETH
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.23789596
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24569595
2026/04/15 15:03:30 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.23898596
2026/04/15 15:03:30 [KuCoin WS 0] bootstrap progress 60/126 last=ENJ-ETH
2026/04/15 15:03:31 [KuCoin WS 0] bootstrap progress 70/126 last=ETH-USDT
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27469592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24373595
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25194595
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27644592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27644592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.25126141
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.25126141
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24548595
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25219594
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.25126141
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.25126141
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [KuCoin WS 0] bootstrap progress 80/126 last=ICP-BTC
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25166595
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.27797178
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27495592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27495592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27495592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.27797178
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.27797178
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.27797178
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27495592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.65664256 final=71.17272185 profit=-0.48392071 pct=-0.675333%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.27797178
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27495592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.47043207
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27495592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.28263592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.47043207
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.47043207
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.28263592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.66507871 final=71.17272185 profit=-0.49235686 pct=-0.687025%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.51526718
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.51526718
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AAVE->BTC reason=small_depth_volume maxStart=6.51526718
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->AAVE reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.28263592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.75106306
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.75106306
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.72014807
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=103.22219266 final=102.67540201 profit=-0.54679065 pct=-0.529722%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.50360272
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=103.22219266 final=102.67540201 profit=-0.54679065 pct=-0.529722%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27495592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.75106306
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.75106306
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24399595
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.50360272
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=100.81299056 final=100.34187015 profit=-0.47112041 pct=-0.467321%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27592592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27592592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=7.12398621
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.12398621
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27592592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=7.07710761
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.07710761
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27592592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=7.12709919
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.12709919
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27592592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=7.07710761
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.86682655
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27592592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.86682655
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.86682655
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AVAX->BTC reason=small_depth_volume maxStart=6.86682655
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.50360272
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.95688402
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.95688402
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AVAX->BTC reason=small_depth_volume maxStart=6.95688402
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27592592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.95688402
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.95688402
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.50360272
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.95688402
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AVAX->BTC reason=small_depth_volume maxStart=6.91788405
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27592592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.91788405
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.91788405
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.91788405
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.46814304
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=104.56633658 final=104.07547442 profit=-0.49086216 pct=-0.469427%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.93493131
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.66435762 final=71.17272185 profit=-0.49163577 pct=-0.686026%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.46814304
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.66435762 final=71.17272185 profit=-0.49163577 pct=-0.686026%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24496595
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.46814304
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.66435762 final=71.17272185 profit=-0.49163577 pct=-0.686026%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25285594
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->HBAR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->HBAR->BTC reason=small_depth_volume maxStart=6.93493131
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->HBAR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->HBAR->BTC reason=small_depth_volume maxStart=6.93493131
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->HBAR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->HBAR->BTC reason=small_depth_volume maxStart=6.93493131
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->HBAR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->HBAR->BTC reason=small_depth_volume maxStart=6.93493131
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.46814304
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.66435762 final=71.17272185 profit=-0.49163577 pct=-0.686026%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.46814304
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.66435762 final=71.17272185 profit=-0.49163577 pct=-0.686026%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AVAX->BTC reason=small_depth_volume maxStart=6.93413131
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AVAX->BTC reason=small_depth_volume maxStart=7.05277630
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24464774
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AVAX->BTC reason=small_depth_volume maxStart=6.79121476
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24464774
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->HBAR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->HBAR->BTC reason=small_depth_volume maxStart=6.79121476
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.79121476
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.79121476
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.79121476
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.79121476
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25192595
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.46814304
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.44606419
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.50360272
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.66435762 final=71.17272185 profit=-0.49163577 pct=-0.686026%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.62008142
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.50360272
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.62008142
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->FET->BTC reason=small_depth_volume maxStart=6.62008142
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.50360272
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.62008142
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50764987 final=87.26853346 profit=-0.23911641 pct=-0.273252%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.62008142
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.65847461
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->ETH reason=depth_non_positive maxStart=55.53255093 final=55.22287275 profit=-0.30967818 pct=-0.557652%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->AVAX->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27499592
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->HBAR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->FET->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->HBAR->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->HBAR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->HBAR->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [KuCoin WS 0] bootstrap progress 90/126 last=IOTX-ETH
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.80847447
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->ONT reason=small_depth_volume maxStart=9.51348186
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->ETH reason=depth_non_positive maxStart=55.53255093 final=55.22287275 profit=-0.30967818 pct=-0.557652%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->FET->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->DASH->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ETH->ONT reason=small_depth_volume maxStart=9.51348186
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->ETH reason=depth_non_positive maxStart=55.53255093 final=55.22287275 profit=-0.30967818 pct=-0.557652%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->DASH->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->ETH reason=depth_non_positive maxStart=55.53255093 final=55.22287275 profit=-0.30967818 pct=-0.557652%
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->DASH->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:31 [DEPTH REJECT] tri=USDT->DASH->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.92463018
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=7.07660175
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.88969926
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.88969926
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.88969926
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27448592
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.28216592
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->DASH->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.28216592
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.28216592
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->DASH->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->A->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [KuCoin WS 0] bootstrap progress 100/126 last=KNC-BTC
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->HYPE->KCS reason=depth_non_positive maxStart=55.94994664 final=55.15238244 profit=-0.79756420 pct=-1.425496%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.97268712
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24738595
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26255593
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->A->BTC reason=small_depth_volume maxStart=7.16236436
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25572594
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=7.16236436
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->STX->BTC reason=small_depth_volume maxStart=7.16236436
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=7.16236436
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.16236436
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.05257910
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.05257910
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->ONT reason=small_depth_volume maxStart=9.58695075
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->ETH reason=depth_non_positive maxStart=55.53255093 final=55.22287275 profit=-0.30967818 pct=-0.557652%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->WIN reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->WIN->BTC reason=small_depth_volume maxStart=7.05989724
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=72.12744000 final=70.88663532 profit=-1.24080468 pct=-1.720295%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=7.05989724
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->ONT reason=small_depth_volume maxStart=9.56277141
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=7.05989724
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->ETH reason=depth_non_positive maxStart=55.53255093 final=55.22287275 profit=-0.30967818 pct=-0.557652%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26255593
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=56.25081000 final=55.28312608 profit=-0.96768392 pct=-1.720302%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=7.05989724
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=7.05989724
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.05989724
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24113596
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25292594
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->WIN reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->ONT reason=small_depth_volume maxStart=9.56277141
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->ONT reason=small_depth_volume maxStart=9.56277141
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->ETH reason=depth_non_positive maxStart=93.75379499 final=93.11637552 profit=-0.63741947 pct=-0.679887%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->ONT reason=small_depth_volume maxStart=9.56277141
2026/04/15 15:03:32 [KuCoin WS 0] bootstrap progress 110/126 last=LTC-ETH
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->WIN reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->WIN reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=56.25081000 final=55.28312608 profit=-0.96768392 pct=-1.720302%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->WIN reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->WIN reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=56.25081000 final=55.28312608 profit=-0.96768392 pct=-1.720302%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=56.25081000 final=55.28312608 profit=-0.96768392 pct=-1.720302%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->ONT reason=small_depth_volume maxStart=9.56277141
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ZEC reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ZEC->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=56.25081000 final=55.28312608 profit=-0.96768392 pct=-1.720302%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=91.85213020 final=90.24663577 profit=-1.60549443 pct=-1.747912%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.02142079
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=91.85213020 final=90.24663577 profit=-1.60549443 pct=-1.747912%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=91.85213020 final=90.24663577 profit=-1.60549443 pct=-1.747912%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=91.85213020 final=90.26018109 profit=-1.59194911 pct=-1.733165%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.21239233
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=91.85213020 final=90.22935381 profit=-1.62277639 pct=-1.766727%
2026/04/15 15:03:32 [KuCoin WS 0] bootstrap progress 120/126 last=NEO-BTC
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=91.85213020 final=90.22935381 profit=-1.62277639 pct=-1.766727%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=91.85213020 final=90.22935381 profit=-1.62277639 pct=-1.766727%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27453592
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=91.85213020 final=90.22935381 profit=-1.62277639 pct=-1.766727%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.21239233
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=93.03267718 final=92.64112158 profit=-0.39155560 pct=-0.420880%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=93.82104383 final=92.16274360 profit=-1.65830023 pct=-1.767514%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=93.82104383 final=93.34047583 profit=-0.48056799 pct=-0.512218%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=99.80848212 final=99.40713474 profit=-0.40134738 pct=-0.402118%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=100.56145000 final=98.78242321 profit=-1.77902679 pct=-1.769094%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=100.56145000 final=100.10704519 profit=-0.45440481 pct=-0.451868%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=102.82945113 final=101.01029658 profit=-1.81915455 pct=-1.769099%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=108.36328856 final=107.80861743 profit=-0.55467112 pct=-0.511863%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=107.30875486 final=105.41017171 profit=-1.89858316 pct=-1.769271%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=107.30875486 final=106.87516197 profit=-0.43359289 pct=-0.404061%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=112.02117177 final=110.04054237 profit=-1.98062940 pct=-1.768085%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=109.88343491 final=109.44105284 profit=-0.44238207 pct=-0.402592%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=109.88343491 final=107.94267931 profit=-1.94075560 pct=-1.766195%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=106.03527458 final=105.47419300 profit=-0.56108158 pct=-0.529146%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.76792960
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=106.03527458 final=105.74556302 profit=-0.28971156 pct=-0.273222%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=106.03527458 final=104.16288263 profit=-1.87239195 pct=-1.765820%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=106.03527458 final=105.47419300 profit=-0.56108158 pct=-0.529146%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=106.03527458 final=104.16288263 profit=-1.87239195 pct=-1.765820%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=106.03527458 final=105.47419300 profit=-0.56108158 pct=-0.529146%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=113.28125597 final=111.27992334 profit=-2.00133263 pct=-1.766694%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.76792960
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=113.28125597 final=112.70800222 profit=-0.57325375 pct=-0.506045%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=113.28125597 final=112.97189905 profit=-0.30935691 pct=-0.273087%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->ETH reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=113.28125597 final=111.27992334 profit=-2.00133263 pct=-1.766694%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->BTC reason=small_depth_volume maxStart=7.16232852
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.76792960
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=113.28125597 final=112.70800222 profit=-0.57325375 pct=-0.506045%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=113.28125597 final=112.98707896 profit=-0.29417701 pct=-0.259687%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.76792960
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50636634 final=87.28235969 profit=-0.22400665 pct=-0.255989%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=113.76598610 final=111.75680593 profit=-2.00918017 pct=-1.766064%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=113.76598610 final=111.75680593 profit=-2.00918017 pct=-1.766064%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=6.65091108
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=43.33295867
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=113.76598610 final=111.75680593 profit=-2.00918017 pct=-1.766064%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=113.76598610 final=111.75680593 profit=-2.00918017 pct=-1.766064%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=7.33931327
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=43.06195893
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24743595
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=7.33931327
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=43.33295867
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=113.76598610 final=113.17484272 profit=-0.59114338 pct=-0.519613%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=7.33931327
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=43.46195855
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25040595
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=7.30488562
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=42.19695976
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25253594
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.76792960
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50636634 final=87.28235969 profit=-0.22400665 pct=-0.255989%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27505592
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.67207450 final=71.17180831 profit=-0.50026619 pct=-0.697993%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=113.76588610 final=113.17484272 profit=-0.59104338 pct=-0.519526%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=113.76588610 final=111.73625468 profit=-2.02963143 pct=-1.784042%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=107.55541683 final=105.64070490 profit=-1.91471193 pct=-1.780210%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=107.55541683 final=105.64070490 profit=-1.91471193 pct=-1.780210%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=107.55541683 final=105.64070490 profit=-1.91471193 pct=-1.780210%
2026/04/15 15:03:32 [KuCoin WS 0] bootstrap progress 126/126 last=ONE-BTC
2026/04/15 15:03:32 [KuCoin WS 0] bootstrap finished 126/126 in 4.234559201s
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=6.70068514
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=49.05395322
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=6.42961407
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=depth_non_positive maxStart=51.35095103 final=50.65861701 profit=-0.69233401 pct=-1.348240%
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=6.42961407
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=42.62995934
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=6.42961407
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=42.62995934
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=6.42961407
2026/04/15 15:03:32 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=42.62995934
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=6.42961407
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=41.74796019
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=107.11909084 final=106.64094396 profit=-0.47814689 pct=-0.446369%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->XMR reason=small_depth_volume maxStart=6.42961407
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->XRP reason=small_depth_volume maxStart=7.76792960
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=107.11909084 final=106.64094396 profit=-0.47814689 pct=-0.446369%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->ETH reason=small_depth_volume maxStart=41.74796019
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XRP->ETH reason=depth_non_positive maxStart=87.50636634 final=87.28235969 profit=-0.22400665 pct=-0.255989%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=107.11909084 final=106.64094396 profit=-0.47814689 pct=-0.446369%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BRL->ETH reason=depth_non_positive maxStart=100.53240582 final=99.87374170 profit=-0.65866412 pct=-0.655176%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->BRL reason=depth_non_positive maxStart=91.69355484 final=88.12944000 profit=-3.56411484 pct=-3.886985%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.09425240
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.09425240
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26375593
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=100.53240582 final=100.10704519 profit=-0.42536064 pct=-0.423108%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=100.53240582 final=98.74155423 profit=-1.79085160 pct=-1.781367%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->EUR reason=small_depth_volume maxStart=32.42346908
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->EUR->ETH reason=depth_non_positive maxStart=71.67135333 final=71.17180831 profit=-0.49954502 pct=-0.696994%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BRL->ETH reason=depth_non_positive maxStart=100.53240582 final=99.87374170 profit=-0.65866412 pct=-0.655176%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ETH->BRL reason=depth_non_positive maxStart=91.69355484 final=88.13123856 profit=-3.56231628 pct=-3.885024%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27658592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=100.15356339 final=98.40059018 profit=-1.75297321 pct=-1.750285%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XLM reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XLM->BTC reason=small_depth_volume maxStart=7.09425240
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XLM reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XLM->BTC reason=small_depth_volume maxStart=7.09425240
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=100.15356339 final=98.40059018 profit=-1.75297321 pct=-1.750285%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XLM reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XLM->BTC reason=small_depth_volume maxStart=7.09425240
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ENJ->ETH reason=depth_non_positive maxStart=100.15356339 final=98.36696085 profit=-1.78660254 pct=-1.783863%
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24948595
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25245594
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26318593
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XLM reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XLM->BTC reason=small_depth_volume maxStart=7.13761782
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27475592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=7.13761782
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24765595
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=7.13761782
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25062595
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26117594
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->ONT reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->ONT->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26275593
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->ICP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->ICP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27558592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27558592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27558592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->WAN->BTC reason=small_depth_volume maxStart=6.99076627
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.99076627
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.99076627
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27558592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24848595
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24848595
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.25145595
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26364593
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->WAN->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27745592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27530592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27530592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XRP reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XRP->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.87541056
2026/04/15 15:03:33 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24820595
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24820595
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24910017
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->BNB reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27870592
2026/04/15 15:03:33 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27870592
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->BNB reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27596592
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.53706354
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->BNB reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27596592
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->XMR reason=depth_max_start depthMaxStart=0.00000000
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->XMR->BTC reason=small_depth_volume maxStart=6.53706354
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27596592
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27596592
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.24886595
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26009594
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.26391593
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27614592
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->EUR reason=small_depth_volume maxStart=4.27614592
2026/04/15 15:03:34 [DEPTH REJECT] tri=USDT->BTC->BNB reason=depth_max_start depthMaxStart=0.00000000
^C2026/04/15 15:03:34 [Main] shutting down...

