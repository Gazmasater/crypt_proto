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






2026/04/13 17:10:14 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:10:24 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:10:34 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:10:44 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:10:54 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:11:04 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:11:15 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:11:25 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:11:35 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:11:45 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:11:55 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:12:05 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:12:15 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:12:25 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:12:35 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:12:45 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:12:55 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:13:05 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:13:15 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:13:25 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:13:35 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:13:45 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:13:55 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:14:05 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:14:15 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:14:25 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=DOGE-USDT localSeq=11689119590 seqStart=11689119594 seqEnd=11689119594
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=BNB-USDT localSeq=11274368401 seqStart=11274368403 seqEnd=11274368403
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=BTC-USDT localSeq=31698806564 seqStart=31698806572 seqEnd=31698806572
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=HYPE-USDT localSeq=4421157145 seqStart=4421157147 seqEnd=4421157147
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=ETC-USDT localSeq=3786010257 seqStart=3786010259 seqEnd=3786010259
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=AVAX-USDT localSeq=9259288270 seqStart=9259288275 seqEnd=9259288275
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=DASH-USDT localSeq=1552271896 seqStart=1552271898 seqEnd=1552271898
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=ETH-USDT localSeq=20731626194 seqStart=20731626198 seqEnd=20731626198
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=BNB-KCS localSeq=1207921797 seqStart=1207921799 seqEnd=1207921799
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=HBAR-USDT localSeq=5313266473 seqStart=5313266475 seqEnd=5313266475
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=LINK-USDT localSeq=13127188922 seqStart=13127188924 seqEnd=13127188924
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=NEAR-USDT localSeq=6888748993 seqStart=6888748996 seqEnd=6888748996
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=KCS-BTC localSeq=3182288135 seqStart=3182288137 seqEnd=3182288137
2026/04/13 17:14:34 [KuCoin WS 0] gap detected symbol=FET-USDT localSeq=3633429941 seqStart=3633429943 seqEnd=3633429943
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=DASH-USDT
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=ETC-USDT
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=HYPE-USDT
2026/04/13 17:14:35 [Calculator] summary checked=0 written=0 profitable=0 best_pct=0.0000% best_usdt=0.000000 best_tri=
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=LINK-USDT
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=BNB-USDT
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=NEAR-USDT
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=BNB-KCS
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=FET-USDT
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=KCS-BTC
2026/04/13 17:14:35 [KuCoin WS 0] resync complete symbol=HBAR-USDT
2026/04/13 17:14:36 [KuCoin WS 0] resync complete symbol=ETH-USDT
2026/04/13 17:14:38 [KuCoin WS 0] resync complete symbol=AVAX-USDT
2026/04/13 17:14:38 [KuCoin WS 0] resync complete symbol=DOGE-USDT
2026/04/13 17:14:39 [KuCoin WS 0] resync complete symbol=BTC-USDT
2026/04/13 17:14:40 [KuCoin WS 0] gap detected symbol=OGN-BTC localSeq=191521559 seqStart=191521561 seqEnd=191521561
2026/04/13 17:14:40 [KuCoin WS 0] resync complete symbol=OGN-BTC
