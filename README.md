2026/04/10 05:53:02 pprof on http://localhost:6060/debug/pprof/
2026/04/10 05:53:04 [KuCoin WS 0] connected
2026/04/10 05:53:04 [KuCoin] started with 1 WS
2026/04/10 05:53:04 [Main] KuCoinCollector started
2026/04/10 05:53:04 [CALC] started | triangles indexed=56 | minVolume=10.00 USDT | minProfit=0.0000% | quoteAgeMax=400ms | logMode=debug
2026/04/10 05:53:04 [REJECT] stage=scan reason=no_quote_leg_2 count=1 tri=USDT->LYX->ETH
2026/04/10 05:53:05 [REJECT] stage=scan reason=no_quote_leg_1 count=1 tri=USDT->VET->ETH
2026/04/10 05:53:07 [REJECT] stage=scan reason=no_quote_leg_1 count=10 tri=USDT->BTC->IOTX
2026/04/10 05:53:07 [REJECT] stage=scan reason=no_quote_leg_2 count=10 tri=USDT->LINK->BTC
2026/04/10 05:53:07 [REJECT] stage=scan reason=no_quote_leg_3 count=1 tri=USDT->XRP->KCS
2026/04/10 05:53:08 [REJECT] stage=scan reason=no_quote_leg_3 count=10 tri=USDT->LYX->ETH
2026/04/10 05:53:08 [REJECT] stage=scan reason=max_start_lt_10 count=1 tri=USDT->LYX->ETH
2026/04/10 05:53:09 [STATS] ticks=75 triangles_seen=102 cand=0 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=1, no_quote_leg_3=15, no_quote_leg_2=29, no_quote_leg_1=57} | exec_rejects={none}
2026/04/10 05:53:10 [EXEC CMP] USDT→LINK→BTC | est=-0.3864% | ideal=-0.0869% | rounded=-0.0987% | final=-0.4069%
2026/04/10 05:53:10 [REJECT] stage=exec reason=profit_below_threshold count=1 tri=USDT->LINK->BTC
2026/04/10 05:53:10 [EXEC CMP] USDT→LINK→BTC | est=-0.3931% | ideal=-0.0936% | rounded=-0.1058% | final=-0.4141%
2026/04/10 05:53:10 [EXEC CMP] USDT→LINK→BTC | est=-0.3931% | ideal=-0.0936% | rounded=-0.1058% | final=-0.4141%
2026/04/10 05:53:11 [REJECT] stage=scan reason=no_quote_leg_1 count=100 tri=USDT->BTC->ATOM
2026/04/10 05:53:12 [EXEC CMP] USDT→LINK→BTC | est=-0.3810% | ideal=-0.0816% | rounded=-0.0911% | final=-0.3994%
2026/04/10 05:53:12 [REJECT] stage=scan reason=no_quote_leg_2 count=100 tri=USDT->BTC->TEL
2026/04/10 05:53:12 [EXEC CMP] USDT→BTC→DASH | est=-0.4699% | ideal=-0.1707% | rounded=-0.2000% | final=-0.6000%
2026/04/10 05:53:12 [EXEC CMP] USDT→LINK→BTC | est=-0.3643% | ideal=-0.0648% | rounded=-0.0767% | final=-0.3851%
2026/04/10 05:53:12 [EXEC CMP] USDT→LINK→BTC | est=-0.3884% | ideal=-0.0890% | rounded=-0.0983% | final=-0.4066%
2026/04/10 05:53:12 [EXEC CMP] USDT→BTC→STX | est=-1.2937% | ideal=-0.9970% | rounded=-1.0020% | final=-1.3000%
2026/04/10 05:53:12 [EXEC CMP] USDT→LINK→BTC | est=-0.3951% | ideal=-0.0957% | rounded=-0.0983% | final=-0.4066%
2026/04/10 05:53:12 [EXEC CMP] USDT→BTC→ICP | est=-0.8399% | ideal=-0.5419% | rounded=-0.5500% | final=-0.8500%
2026/04/10 05:53:12 [REJECT] stage=exec reason=profit_below_threshold count=10 tri=USDT->BTC->ICP
2026/04/10 05:53:12 [EXEC CMP] USDT→LINK→BTC | est=-0.3871% | ideal=-0.0876% | rounded=-0.0911% | final=-0.3994%
2026/04/10 05:53:12 [EXEC CMP] USDT→LINK→BTC | est=-0.3721% | ideal=-0.0726% | rounded=-0.0761% | final=-0.3845%
2026/04/10 05:53:12 [EXEC CMP] USDT→BTC→ICP | est=-0.8548% | ideal=-0.5568% | rounded=-0.5700% | final=-0.8700%
2026/04/10 05:53:12 [EXEC CMP] USDT→BTC→STX | est=-1.3085% | ideal=-1.0119% | rounded=-1.0160% | final=-1.3140%
2026/04/10 05:53:12 [EXEC CMP] USDT→BTC→STX | est=-0.9902% | ideal=-0.6925% | rounded=-0.6970% | final=-0.9970%
2026/04/10 05:53:12 [EXEC CMP] USDT→BTC→ATOM | est=-0.4742% | ideal=-0.1751% | rounded=-0.1800% | final=-0.4830%
2026/04/10 05:53:12 [EXEC CMP] USDT→BTC→ETH | est=-0.3350% | ideal=-0.0355% | rounded=-0.0396% | final=-0.3406%
2026/04/10 05:53:14 [STATS] ticks=217 triangles_seen=404 cand=17 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=2, no_quote_leg_3=61, no_quote_leg_1=157, no_quote_leg_2=167} | exec_rejects={profit_below_threshold=17}
2026/04/10 05:53:14 [EXEC CMP] USDT→ETC→ETH | est=-0.5463% | ideal=-0.2474% | rounded=-0.2477% | final=-0.5661%
2026/04/10 05:53:14 [EXEC CMP] USDT→ETC→ETH | est=-0.5205% | ideal=-0.2215% | rounded=-0.2258% | final=-0.5442%
2026/04/10 05:53:14 [EXEC CMP] USDT→ETC→ETH | est=-0.5010% | ideal=-0.2019% | rounded=-0.2062% | final=-0.5246%
2026/04/10 05:53:14 [EXEC CMP] USDT→XRP→ETH | est=-0.3342% | ideal=-0.0346% | rounded=-0.0355% | final=-0.3366%
2026/04/10 05:53:14 [EXEC CMP] USDT→XRP→KCS | est=-0.4486% | ideal=-0.1493% | rounded=-0.1580% | final=-0.4520%
2026/04/10 05:53:14 [EXEC CMP] USDT→XRP→ETH | est=-0.3179% | ideal=-0.0183% | rounded=-0.0202% | final=-0.3213%
2026/04/10 05:53:14 [EXEC CMP] USDT→XRP→KCS | est=-0.4486% | ideal=-0.1493% | rounded=-0.1580% | final=-0.4520%
2026/04/10 05:53:15 [EXEC CMP] USDT→LYX→ETH | est=-1.6434% | ideal=-1.3478% | rounded=-1.3491% | final=-1.6467%
2026/04/10 05:53:15 [EXEC CMP] USDT→ETC→ETH | est=-0.4756% | ideal=-0.1764% | rounded=-0.1848% | final=-0.5033%
2026/04/10 05:53:15 [EXEC CMP] USDT→ETC→ETH | est=-0.4498% | ideal=-0.1505% | rounded=-0.1629% | final=-0.4595%
2026/04/10 05:53:15 [EXEC CMP] USDT→LYX→ETH | est=-1.8059% | ideal=-1.5107% | rounded=-1.5110% | final=-1.8063%
2026/04/10 05:53:15 [EXEC CMP] USDT→ETC→ETH | est=-0.4498% | ideal=-0.1505% | rounded=-0.1629% | final=-0.4595%
2026/04/10 05:53:16 [EXEC CMP] USDT→ETC→ETH | est=-0.4308% | ideal=-0.1314% | rounded=-0.1479% | final=-0.4445%
2026/04/10 05:53:16 [EXEC CMP] USDT→ETC→ETH | est=-0.4308% | ideal=-0.1314% | rounded=-0.1479% | final=-0.4445%
2026/04/10 05:53:16 [REJECT] stage=scan reason=max_start_lt_10 count=10 tri=USDT->LYX->ETH
2026/04/10 05:53:17 [EXEC CMP] USDT→SHIB→DOGE | est=-0.9126% | ideal=-0.6147% | rounded=-0.6149% | final=-0.9127%
2026/04/10 05:53:17 [EXEC CMP] USDT→LINK→BTC | est=-0.3493% | ideal=-0.0497% | rounded=-0.0550% | final=-0.3634%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→FET | est=-0.5690% | ideal=-0.2701% | rounded=-0.2750% | final=-0.5750%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→STX | est=-1.3384% | ideal=-1.0418% | rounded=-1.0470% | final=-1.3440%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→NEAR | est=-0.3544% | ideal=-0.0549% | rounded=-0.0610% | final=-0.3610%
2026/04/10 05:53:17 [REJECT] stage=scan reason=no_quote_leg_3 count=100 tri=USDT->BTC->XLM
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→DASH | est=-0.4918% | ideal=-0.1927% | rounded=-0.3000% | final=-0.6000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ADA | est=-0.5787% | ideal=-0.2798% | rounded=-0.2880% | final=-0.5890%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→FET | est=-0.5690% | ideal=-0.2701% | rounded=-0.2750% | final=-0.5750%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→AVAX | est=-0.3780% | ideal=-0.0785% | rounded=-0.0900% | final=-0.4000%
2026/04/10 05:53:17 [EXEC CMP] USDT→LINK→BTC | est=-0.3733% | ideal=-0.0738% | rounded=-0.0766% | final=-0.3849%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→DASH | est=-0.4918% | ideal=-0.1927% | rounded=-0.3000% | final=-0.6000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→AVAX | est=-0.3994% | ideal=-0.1000% | rounded=-0.1100% | final=-0.4200%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ATOM | est=-0.4979% | ideal=-0.1988% | rounded=-0.2060% | final=-0.5070%
2026/04/10 05:53:17 [REJECT] stage=scan reason=max_start_zero count=1 tri=USDT->BTC->KRL
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→AVAX | est=-0.3994% | ideal=-0.1000% | rounded=-0.1100% | final=-0.4200%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→DASH | est=-0.4733% | ideal=-0.1741% | rounded=-0.2000% | final=-0.6000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→NEO | est=-1.0170% | ideal=-0.7195% | rounded=-0.7240% | final=-1.0223%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ATOM | est=-0.4979% | ideal=-0.1988% | rounded=-0.2060% | final=-0.5070%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→NEO | est=-1.0170% | ideal=-0.7195% | rounded=-0.7240% | final=-1.0223%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ADA | est=-0.6182% | ideal=-0.3194% | rounded=-0.3270% | final=-0.6290%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→AVAX | est=-0.3887% | ideal=-0.0892% | rounded=-0.1000% | final=-0.4100%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ATOM | est=-0.5767% | ideal=-0.2779% | rounded=-0.2830% | final=-0.5860%
2026/04/10 05:53:17 [EXEC CMP] USDT→SHIB→DOGE | est=-0.8958% | ideal=-0.5979% | rounded=-0.5979% | final=-0.8960%
2026/04/10 05:53:17 [EXEC CMP] USDT→LYX→ETH | est=-1.6578% | ideal=-1.3621% | rounded=-1.3628% | final=-1.6602%
2026/04/10 05:53:17 [EXEC CMP] USDT→XRP→ETH | est=-0.3372% | ideal=-0.0376% | rounded=-0.0396% | final=-0.3406%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ETH | est=-0.3289% | ideal=-0.0293% | rounded=-0.0352% | final=-0.3362%
2026/04/10 05:53:17 [EXEC CMP] USDT→MANA→ETH | est=-0.6813% | ideal=-0.3827% | rounded=-0.3830% | final=-0.6836%
2026/04/10 05:53:17 [EXEC CMP] USDT→ATOM→ETH | est=-0.5047% | ideal=-0.2056% | rounded=-0.2277% | final=-0.5241%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→DASH | est=-0.4733% | ideal=-0.1741% | rounded=-0.2000% | final=-0.6000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→AVAX | est=-0.3810% | ideal=-0.0815% | rounded=-0.0900% | final=-0.4000%
2026/04/10 05:53:17 [EXEC CMP] USDT→XRP→KCS | est=-0.4352% | ideal=-0.1359% | rounded=-0.1410% | final=-0.4430%
2026/04/10 05:53:17 [EXEC CMP] USDT→DOT→KCS | est=-0.6237% | ideal=-0.3250% | rounded=-0.3260% | final=-0.6280%
2026/04/10 05:53:17 [EXEC CMP] USDT→KCS→SOL | est=-0.5982% | ideal=-0.2994% | rounded=-0.4000% | final=-2.2000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ICP | est=-0.8902% | ideal=-0.5923% | rounded=-0.6000% | final=-0.9000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→NEO | est=-0.7619% | ideal=-0.4636% | rounded=-0.4681% | final=-0.7672%
2026/04/10 05:53:17 [EXEC CMP] USDT→LINK→BTC | est=-0.3733% | ideal=-0.0738% | rounded=-0.0766% | final=-0.3849%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→DASH | est=-0.4918% | ideal=-0.1927% | rounded=-0.3000% | final=-0.6000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→FET | est=-0.5690% | ideal=-0.2701% | rounded=-0.2750% | final=-0.5750%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→NEO | est=-0.7619% | ideal=-0.4636% | rounded=-0.4681% | final=-0.7672%
2026/04/10 05:53:17 [EXEC CMP] USDT→XRP→KCS | est=-0.4477% | ideal=-0.1485% | rounded=-0.1500% | final=-0.4520%
2026/04/10 05:53:17 [EXEC CMP] USDT→LINK→BTC | est=-0.3768% | ideal=-0.0773% | rounded=-0.0801% | final=-0.3884%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ICP | est=-0.8868% | ideal=-0.5888% | rounded=-0.6000% | final=-0.9000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ATOM | est=-0.5733% | ideal=-0.2744% | rounded=-0.2760% | final=-0.5790%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→FET | est=-0.5655% | ideal=-0.2666% | rounded=-0.2680% | final=-0.5680%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ADA | est=-0.6147% | ideal=-0.3160% | rounded=-0.3270% | final=-0.6290%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→DASH | est=-0.4883% | ideal=-0.1892% | rounded=-0.3000% | final=-0.6000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→AVAX | est=-0.3775% | ideal=-0.0781% | rounded=-0.0900% | final=-0.3900%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→NEO | est=-0.7585% | ideal=-0.4602% | rounded=-0.4610% | final=-0.7600%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ETH | est=-0.3255% | ideal=-0.0258% | rounded=-0.0287% | final=-0.3297%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→STX | est=-1.3350% | ideal=-1.0384% | rounded=-1.0400% | final=-1.3370%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→STX | est=-1.0167% | ideal=-0.7192% | rounded=-0.7210% | final=-1.0200%
2026/04/10 05:53:17 [EXEC CMP] USDT→KCS→SOL | est=-0.5982% | ideal=-0.2994% | rounded=-0.4000% | final=-2.2000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ATOM | est=-0.5733% | ideal=-0.2744% | rounded=-0.2760% | final=-0.5790%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→NEO | est=-0.7513% | ideal=-0.4530% | rounded=-0.4538% | final=-0.7529%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ICP | est=-0.8868% | ideal=-0.5888% | rounded=-0.6000% | final=-0.9000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→DASH | est=-0.4883% | ideal=-0.1892% | rounded=-0.3000% | final=-0.6000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→AVAX | est=-0.3882% | ideal=-0.0888% | rounded=-0.1000% | final=-0.4000%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→ATOM | est=-0.4550% | ideal=-0.1557% | rounded=-0.1590% | final=-0.4590%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→AVAX | est=-0.3959% | ideal=-0.0965% | rounded=-0.1000% | final=-0.4100%
2026/04/10 05:53:17 [EXEC CMP] USDT→XRP→KCS | est=-0.4189% | ideal=-0.1195% | rounded=-0.1240% | final=-0.4260%
2026/04/10 05:53:17 [EXEC CMP] USDT→XRP→ETH | est=-0.3533% | ideal=-0.0538% | rounded=-0.0549% | final=-0.3559%
2026/04/10 05:53:17 [EXEC CMP] USDT→XRP→ETH | est=-0.3371% | ideal=-0.0375% | rounded=-0.0396% | final=-0.3406%
2026/04/10 05:53:17 [EXEC CMP] USDT→BTC→DASH | est=-0.4883% | ideal=-0.1892% | rounded=-0.3000% | final=-0.6000%
2026/04/10 05:53:17 [EXEC CMP] USDT→XRP→KCS | est=-0.4204% | ideal=-0.1210% | rounded=-0.1240% | final=-0.4260%
2026/04/10 05:53:17 [EXEC CMP] USDT→XRP→ETH | est=-0.3386% | ideal=-0.0390% | rounded=-0.0396% | final=-0.3406%
2026/04/10 05:53:19 [STATS] ticks=500 triangles_seen=853 cand=95 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=2, max_start_lt_10=37, no_quote_leg_3=134, no_quote_leg_2=259, no_quote_leg_1=326} | exec_rejects={profit_below_threshold=95}
2026/04/10 05:53:21 [EXEC CMP] USDT→XRP→ETH | est=-0.3227% | ideal=-0.0231% | rounded=-0.0247% | final=-0.3256%
2026/04/10 05:53:21 [EXEC CMP] USDT→LYX→ETH | est=-1.5338% | ideal=-1.2379% | rounded=-1.2385% | final=-1.5361%
2026/04/10 05:53:21 [EXEC CMP] USDT→LYX→ETH | est=-2.6169% | ideal=-2.3242% | rounded=-2.3255% | final=-2.6198%
2026/04/10 05:53:21 [EXEC CMP] USDT→LYX→ETH | est=-2.9383% | ideal=-2.6465% | rounded=-2.6470% | final=-2.9410%
2026/04/10 05:53:24 [STATS] ticks=633 triangles_seen=1017 cand=99 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=2, max_start_lt_10=39, no_quote_leg_3=161, no_quote_leg_2=282, no_quote_leg_1=434} | exec_rejects={profit_below_threshold=99}
2026/04/10 05:53:25 [EXEC CMP] USDT→BTC→DASH | est=-0.4326% | ideal=-0.1333% | rounded=-0.2000% | final=-0.5000%
2026/04/10 05:53:25 [REJECT] stage=exec reason=profit_below_threshold count=100 tri=USDT->BTC->DASH
2026/04/10 05:53:25 [EXEC CMP] USDT→BTC→STX | est=-1.0613% | ideal=-0.7639% | rounded=-0.7650% | final=-1.0650%
2026/04/10 05:53:28 [EXEC CMP] USDT→KCS→SOL | est=-0.5681% | ideal=-0.2692% | rounded=-0.4000% | final=-2.2000%
2026/04/10 05:53:28 [EXEC CMP] USDT→KCS→SOL | est=-0.5681% | ideal=-0.2692% | rounded=-0.4000% | final=-2.2000%
2026/04/10 05:53:29 [STATS] ticks=753 triangles_seen=1195 cand=103 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=2, max_start_lt_10=42, no_quote_leg_3=184, no_quote_leg_2=323, no_quote_leg_1=541} | exec_rejects={profit_below_threshold=103}
2026/04/10 05:53:30 [EXEC CMP] USDT→LYX→ETH | est=-2.5456% | ideal=-2.2526% | rounded=-2.2544% | final=-2.5466%
2026/04/10 05:53:34 [STATS] ticks=818 triangles_seen=1290 cand=104 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=2, max_start_lt_10=43, no_quote_leg_3=187, no_quote_leg_2=330, no_quote_leg_1=624} | exec_rejects={profit_below_threshold=104}
2026/04/10 05:53:35 [EXEC CMP] USDT→LINK→BTC | est=-0.3787% | ideal=-0.0793% | rounded=-0.0874% | final=-0.3885%
2026/04/10 05:53:35 [EXEC CMP] USDT→BTC→DASH | est=-0.4480% | ideal=-0.1488% | rounded=-0.2000% | final=-0.5000%
2026/04/10 05:53:35 [EXEC CMP] USDT→BTC→DASH | est=-0.4295% | ideal=-0.1302% | rounded=-0.2000% | final=-0.5000%
2026/04/10 05:53:35 [EXEC CMP] USDT→LINK→BTC | est=-0.3787% | ideal=-0.0793% | rounded=-0.0874% | final=-0.3885%
2026/04/10 05:53:35 [EXEC CMP] USDT→LINK→BTC | est=-0.3786% | ideal=-0.0791% | rounded=-0.0873% | final=-0.3884%
2026/04/10 05:53:35 [EXEC CMP] USDT→BTC→DASH | est=-0.4295% | ideal=-0.1302% | rounded=-0.2000% | final=-0.5000%
2026/04/10 05:53:35 [EXEC CMP] USDT→BTC→DASH | est=-0.4295% | ideal=-0.1302% | rounded=-0.2000% | final=-0.5000%
2026/04/10 05:53:35 [EXEC CMP] USDT→LINK→BTC | est=-0.3706% | ideal=-0.0711% | rounded=-0.0729% | final=-0.3812%
2026/04/10 05:53:35 [EXEC CMP] USDT→LINK→BTC | est=-0.3694% | ideal=-0.0699% | rounded=-0.0729% | final=-0.3812%
2026/04/10 05:53:35 [EXEC CMP] USDT→LINK→BTC | est=-0.3775% | ideal=-0.0780% | rounded=-0.0873% | final=-0.3884%
2026/04/10 05:53:35 [EXEC CMP] USDT→LINK→BTC | est=-0.3819% | ideal=-0.0825% | rounded=-0.0944% | final=-0.4028%
2026/04/10 05:53:37 [EXEC CMP] USDT→BTC→ICP | est=-0.9012% | ideal=-0.6033% | rounded=-0.6100% | final=-0.9200%
2026/04/10 05:53:37 [EXEC CMP] USDT→BTC→ATOM | est=-0.4583% | ideal=-0.1591% | rounded=-0.1660% | final=-0.4670%
2026/04/10 05:53:37 [EXEC CMP] USDT→LINK→BTC | est=-0.3730% | ideal=-0.0735% | rounded=-0.0839% | final=-0.3851%
2026/04/10 05:53:37 [EXEC CMP] USDT→BTC→FET | est=-0.5982% | ideal=-0.2994% | rounded=-0.3050% | final=-0.6050%
2026/04/10 05:53:37 [EXEC CMP] USDT→BTC→ICP | est=-0.8733% | ideal=-0.5753% | rounded=-0.5900% | final=-0.8900%
2026/04/10 05:53:37 [EXEC CMP] USDT→BTC→FET | est=-0.5982% | ideal=-0.2994% | rounded=-0.3050% | final=-0.6050%
2026/04/10 05:53:37 [EXEC CMP] USDT→BTC→ICP | est=-0.8733% | ideal=-0.5753% | rounded=-0.5900% | final=-0.8900%
2026/04/10 05:53:37 [EXEC CMP] USDT→BTC→ICP | est=-0.8453% | ideal=-0.5473% | rounded=-0.5600% | final=-0.8600%
2026/04/10 05:53:39 [STATS] ticks=923 triangles_seen=1533 cand=123 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=4, max_start_lt_10=43, no_quote_leg_3=215, no_quote_leg_2=450, no_quote_leg_1=698} | exec_rejects={profit_below_threshold=123}
2026/04/10 05:53:42 [EXEC CMP] USDT→XRP→ETH | est=-0.3473% | ideal=-0.0477% | rounded=-0.0487% | final=-0.3497%
2026/04/10 05:53:42 [EXEC CMP] USDT→XRP→ETH | est=-0.3473% | ideal=-0.0477% | rounded=-0.0487% | final=-0.3497%
2026/04/10 05:53:44 [STATS] ticks=1011 triangles_seen=1669 cand=125 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=4, max_start_lt_10=43, no_quote_leg_3=245, no_quote_leg_2=470, no_quote_leg_1=782} | exec_rejects={profit_below_threshold=125}
2026/04/10 05:53:46 [EXEC CMP] USDT→DOT→KCS | est=-0.6391% | ideal=-0.3404% | rounded=-0.3430% | final=-0.6450%
2026/04/10 05:53:46 [EXEC CMP] USDT→DOT→KCS | est=-0.6622% | ideal=-0.3636% | rounded=-0.3680% | final=-0.6700%
2026/04/10 05:53:46 [EXEC CMP] USDT→DOT→KCS | est=-0.6622% | ideal=-0.3636% | rounded=-0.3680% | final=-0.6700%
2026/04/10 05:53:46 [EXEC CMP] USDT→DOT→KCS | est=-0.5972% | ideal=-0.2984% | rounded=-0.3010% | final=-0.6030%
2026/04/10 05:53:46 [EXEC CMP] USDT→DOT→KCS | est=-0.5819% | ideal=-0.2830% | rounded=-0.2840% | final=-0.5860%
2026/04/10 05:53:49 [STATS] ticks=1108 triangles_seen=1812 cand=130 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=4, max_start_lt_10=43, no_quote_leg_3=268, no_quote_leg_2=492, no_quote_leg_1=875} | exec_rejects={profit_below_threshold=130}
2026/04/10 05:53:49 [EXEC CMP] USDT→BTC→DASH | est=-0.4603% | ideal=-0.1611% | rounded=-0.2000% | final=-0.5000%
2026/04/10 05:53:54 [STATS] ticks=1166 triangles_seen=1934 cand=131 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=4, max_start_lt_10=43, no_quote_leg_3=286, no_quote_leg_2=540, no_quote_leg_1=930} | exec_rejects={profit_below_threshold=131}
2026/04/10 05:53:59 [STATS] ticks=1241 triangles_seen=2020 cand=131 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=4, max_start_lt_10=43, no_quote_leg_3=305, no_quote_leg_2=553, no_quote_leg_1=984} | exec_rejects={profit_below_threshold=131}
2026/04/10 05:54:00 [REJECT] stage=scan reason=no_quote_leg_1 count=1000 tri=USDT->BTC->DASH
2026/04/10 05:54:02 [EXEC CMP] USDT→BTC→DOT | est=-0.3982% | ideal=-0.0988% | rounded=-0.1000% | final=-0.4020%
2026/04/10 05:54:02 [EXEC CMP] USDT→LINK→BTC | est=-0.3556% | ideal=-0.0561% | rounded=-0.0579% | final=-0.3663%
2026/04/10 05:54:02 [EXEC CMP] USDT→LINK→BTC | est=-0.3556% | ideal=-0.0561% | rounded=-0.0579% | final=-0.3663%
2026/04/10 05:54:02 [EXEC CMP] USDT→BTC→DASH | est=-0.4525% | ideal=-0.1532% | rounded=-0.2000% | final=-0.5000%
2026/04/10 05:54:04 [STATS] ticks=1349 triangles_seen=2209 cand=135 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=4, max_start_lt_10=43, no_quote_leg_3=338, no_quote_leg_2=595, no_quote_leg_1=1094} | exec_rejects={profit_below_threshold=135}
2026/04/10 05:54:06 [EXEC CMP] USDT→BTC→DASH | est=-0.4785% | ideal=-0.1794% | rounded=-0.3000% | final=-0.5000%
2026/04/10 05:54:06 [EXEC CMP] USDT→BTC→FET | est=-0.5851% | ideal=-0.2863% | rounded=-0.2940% | final=-0.5940%
2026/04/10 05:54:06 [EXEC CMP] USDT→LINK→BTC | est=-0.3605% | ideal=-0.0610% | rounded=-0.0648% | final=-0.3732%
2026/04/10 05:54:06 [EXEC CMP] USDT→BTC→ICP | est=-0.8891% | ideal=-0.5912% | rounded=-0.6000% | final=-0.9100%
2026/04/10 05:54:06 [EXEC CMP] USDT→LINK→BTC | est=-0.3685% | ideal=-0.0690% | rounded=-0.0720% | final=-0.3804%
2026/04/10 05:54:07 [EXEC CMP] USDT→BTC→DASH | est=-0.4785% | ideal=-0.1794% | rounded=-0.3000% | final=-0.5000%
2026/04/10 05:54:07 [EXEC CMP] USDT→BTC→DASH | est=-0.5300% | ideal=-0.2310% | rounded=-0.3000% | final=-0.6000%
2026/04/10 05:54:09 [STATS] ticks=1481 triangles_seen=2420 cand=142 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=45, no_quote_leg_3=369, no_quote_leg_2=646, no_quote_leg_1=1212} | exec_rejects={profit_below_threshold=142}
2026/04/10 05:54:14 [STATS] ticks=1535 triangles_seen=2489 cand=142 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=45, no_quote_leg_3=375, no_quote_leg_2=659, no_quote_leg_1=1262} | exec_rejects={profit_below_threshold=142}
2026/04/10 05:54:19 [STATS] ticks=1619 triangles_seen=2592 cand=142 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=45, no_quote_leg_3=392, no_quote_leg_2=685, no_quote_leg_1=1322} | exec_rejects={profit_below_threshold=142}
2026/04/10 05:54:24 [STATS] ticks=1681 triangles_seen=2667 cand=142 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=45, no_quote_leg_3=396, no_quote_leg_2=703, no_quote_leg_1=1375} | exec_rejects={profit_below_threshold=142}
2026/04/10 05:54:29 [STATS] ticks=1781 triangles_seen=2804 cand=142 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=45, no_quote_leg_3=410, no_quote_leg_2=736, no_quote_leg_1=1465} | exec_rejects={profit_below_threshold=142}
2026/04/10 05:54:29 [EXEC CMP] USDT→SHIB→DOGE | est=-0.6867% | ideal=-0.3881% | rounded=-0.3882% | final=-0.6869%
2026/04/10 05:54:29 [EXEC CMP] USDT→SHIB→DOGE | est=-0.6698% | ideal=-0.3712% | rounded=-0.3713% | final=-0.6700%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3577% | ideal=-0.0582% | rounded=-0.0719% | final=-0.3802%
2026/04/10 05:54:32 [EXEC CMP] USDT→SHIB→DOGE | est=-0.6759% | ideal=-0.3774% | rounded=-0.3775% | final=-0.6761%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→FET | est=-0.4855% | ideal=-0.1864% | rounded=-0.1940% | final=-0.4940%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→FET | est=-0.4561% | ideal=-0.1569% | rounded=-0.1640% | final=-0.4650%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3577% | ideal=-0.0582% | rounded=-0.0719% | final=-0.3802%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3677% | ideal=-0.0682% | rounded=-0.0791% | final=-0.3874%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→EWT | est=-0.7031% | ideal=-0.4046% | rounded=-0.4502% | final=-0.7782%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3711% | ideal=-0.0716% | rounded=-0.0791% | final=-0.3874%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→ICP | est=-0.7101% | ideal=-0.4117% | rounded=-0.4300% | final=-0.7300%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3582% | ideal=-0.0587% | rounded=-0.0661% | final=-0.3745%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→ICP | est=-0.7230% | ideal=-0.4246% | rounded=-0.4300% | final=-0.7300%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→EWT | est=-0.7159% | ideal=-0.4175% | rounded=-0.4502% | final=-0.7782%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→FET | est=-0.4690% | ideal=-0.1698% | rounded=-0.1720% | final=-0.4720%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→ATOM | est=-0.4595% | ideal=-0.1603% | rounded=-0.1640% | final=-0.4650%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3582% | ideal=-0.0587% | rounded=-0.0661% | final=-0.3745%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→FET | est=-0.4690% | ideal=-0.1698% | rounded=-0.1720% | final=-0.4720%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→FET | est=-0.4280% | ideal=-0.1287% | rounded=-0.1310% | final=-0.4310%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→ICP | est=-0.9188% | ideal=-0.6210% | rounded=-0.6300% | final=-0.9300%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3488% | ideal=-0.0492% | rounded=-0.0567% | final=-0.3651%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→ICP | est=-0.9282% | ideal=-0.6304% | rounded=-0.6400% | final=-0.9400%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→EWT | est=-0.7253% | ideal=-0.4269% | rounded=-0.4502% | final=-0.7782%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→ATOM | est=-0.4689% | ideal=-0.1697% | rounded=-0.1780% | final=-0.4790%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→FET | est=-0.4374% | ideal=-0.1381% | rounded=-0.1450% | final=-0.4460%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3610% | ideal=-0.0615% | rounded=-0.0711% | final=-0.3795%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→ATOM | est=-0.4689% | ideal=-0.1697% | rounded=-0.1780% | final=-0.4790%
2026/04/10 05:54:32 [EXEC CMP] USDT→XRP→KCS | est=-0.4551% | ideal=-0.1558% | rounded=-0.1630% | final=-0.4650%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3771% | ideal=-0.0776% | rounded=-0.0927% | final=-0.3939%
2026/04/10 05:54:32 [EXEC CMP] USDT→BTC→ICP | est=-0.9003% | ideal=-0.6024% | rounded=-0.6200% | final=-0.9200%
2026/04/10 05:54:32 [EXEC CMP] USDT→ATOM→ETH | est=-0.4811% | ideal=-0.1819% | rounded=-0.1930% | final=-0.5114%
2026/04/10 05:54:32 [EXEC CMP] USDT→XRP→KCS | est=-0.4551% | ideal=-0.1558% | rounded=-0.1630% | final=-0.4650%
2026/04/10 05:54:32 [EXEC CMP] USDT→SHIB→DOGE | est=-0.6989% | ideal=-0.4004% | rounded=-0.4005% | final=-0.6991%
2026/04/10 05:54:32 [EXEC CMP] USDT→LINK→BTC | est=-0.3691% | ideal=-0.0696% | rounded=-0.0783% | final=-0.3867%
2026/04/10 05:54:33 [EXEC CMP] USDT→SHIB→DOGE | est=-0.6821% | ideal=-0.3835% | rounded=-0.3836% | final=-0.6822%
2026/04/10 05:54:33 [EXEC CMP] USDT→DOT→KCS | est=-0.9214% | ideal=-0.6235% | rounded=-0.6260% | final=-0.9280%
2026/04/10 05:54:33 [EXEC CMP] USDT→XRP→KCS | est=-0.4551% | ideal=-0.1558% | rounded=-0.1630% | final=-0.4650%
2026/04/10 05:54:33 [EXEC CMP] USDT→LINK→BTC | est=-0.3657% | ideal=-0.0662% | rounded=-0.0711% | final=-0.3795%
2026/04/10 05:54:33 [EXEC CMP] USDT→LINK→BTC | est=-0.3577% | ideal=-0.0582% | rounded=-0.0639% | final=-0.3723%
2026/04/10 05:54:33 [EXEC CMP] USDT→SHIB→DOGE | est=-0.6989% | ideal=-0.4004% | rounded=-0.4005% | final=-0.6991%
2026/04/10 05:54:33 [EXEC CMP] USDT→DOT→KCS | est=-0.9331% | ideal=-0.6354% | rounded=-0.6380% | final=-0.9400%
2026/04/10 05:54:33 [EXEC CMP] USDT→SHIB→DOGE | est=-0.6821% | ideal=-0.3835% | rounded=-0.3836% | final=-0.6822%
2026/04/10 05:54:33 [EXEC CMP] USDT→DOT→KCS | est=-0.8487% | ideal=-0.5507% | rounded=-0.5590% | final=-0.8610%
2026/04/10 05:54:34 [STATS] ticks=1948 triangles_seen=3131 cand=185 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=47, no_quote_leg_3=473, no_quote_leg_2=883, no_quote_leg_1=1537} | exec_rejects={profit_below_threshold=185}
2026/04/10 05:54:38 [EXEC CMP] USDT→BTC→SNX | est=-1.8506% | ideal=-1.5556% | rounded=-1.5690% | final=-1.8740%
2026/04/10 05:54:38 [EXEC CMP] USDT→BTC→DOT | est=-0.4440% | ideal=-0.1447% | rounded=-0.1470% | final=-0.4480%
2026/04/10 05:54:38 [EXEC CMP] USDT→DOT→KCS | est=-0.9290% | ideal=-0.6312% | rounded=-0.6340% | final=-0.9360%
2026/04/10 05:54:39 [EXEC CMP] USDT→DOT→KCS | est=-0.9408% | ideal=-0.6430% | rounded=-0.6460% | final=-0.9480%
2026/04/10 05:54:39 [STATS] ticks=2046 triangles_seen=3341 cand=189 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=49, no_quote_leg_3=495, no_quote_leg_2=971, no_quote_leg_1=1631} | exec_rejects={profit_below_threshold=189}
2026/04/10 05:54:40 [EXEC CMP] USDT→DOT→KCS | est=-0.7914% | ideal=-0.4932% | rounded=-0.5000% | final=-0.8020%
2026/04/10 05:54:43 [REJECT] stage=scan reason=no_quote_leg_2 count=1000 tri=USDT->BTC->XDC
2026/04/10 05:54:43 [EXEC CMP] USDT→BTC→KCS | est=-0.4163% | ideal=-0.1170% | rounded=-0.1240% | final=-0.4260%
2026/04/10 05:54:44 [EXEC CMP] USDT→BTC→KCS | est=-0.4163% | ideal=-0.1170% | rounded=-0.1240% | final=-0.4260%
2026/04/10 05:54:44 [EXEC CMP] USDT→BTC→KCS | est=-0.4163% | ideal=-0.1170% | rounded=-0.1240% | final=-0.4260%
2026/04/10 05:54:44 [STATS] ticks=2186 triangles_seen=3577 cand=193 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=56, no_quote_leg_3=522, no_quote_leg_2=1030, no_quote_leg_1=1770} | exec_rejects={profit_below_threshold=193}
2026/04/10 05:54:44 [EXEC CMP] USDT→XRP→KCS | est=-0.4635% | ideal=-0.1643% | rounded=-0.1710% | final=-0.4730%
2026/04/10 05:54:44 [EXEC CMP] USDT→XRP→KCS | est=-0.4769% | ideal=-0.1777% | rounded=-0.1830% | final=-0.4850%
2026/04/10 05:54:45 [EXEC CMP] USDT→DOT→KCS | est=-0.7720% | ideal=-0.4737% | rounded=-0.4800% | final=-0.7820%
2026/04/10 05:54:45 [EXEC CMP] USDT→DOT→KCS | est=-0.7720% | ideal=-0.4737% | rounded=-0.4800% | final=-0.7820%
2026/04/10 05:54:45 [EXEC CMP] USDT→DOT→KCS | est=-0.7566% | ideal=-0.4583% | rounded=-0.4630% | final=-0.7650%
2026/04/10 05:54:46 [EXEC CMP] USDT→DOT→KCS | est=-0.7684% | ideal=-0.4701% | rounded=-0.4740% | final=-0.7770%
2026/04/10 05:54:46 [EXEC CMP] USDT→DOT→KCS | est=-0.8032% | ideal=-0.5050% | rounded=-0.5120% | final=-0.8140%
2026/04/10 05:54:46 [EXEC CMP] USDT→XRP→KCS | est=-0.4783% | ideal=-0.1792% | rounded=-0.1830% | final=-0.4850%
2026/04/10 05:54:46 [EXEC CMP] USDT→XRP→KCS | est=-0.4665% | ideal=-0.1673% | rounded=-0.1710% | final=-0.4730%
2026/04/10 05:54:46 [EXEC CMP] USDT→DOT→KCS | est=-0.7914% | ideal=-0.4932% | rounded=-0.5000% | final=-0.8020%
2026/04/10 05:54:46 [EXEC CMP] USDT→XRP→ETH | est=-0.3274% | ideal=-0.0277% | rounded=-0.0290% | final=-0.3301%
2026/04/10 05:54:46 [EXEC CMP] USDT→DOT→KCS | est=-0.7914% | ideal=-0.4932% | rounded=-0.5000% | final=-0.8020%
2026/04/10 05:54:46 [EXEC CMP] USDT→XRP→ETH | est=-0.3274% | ideal=-0.0277% | rounded=-0.0290% | final=-0.3301%
2026/04/10 05:54:46 [EXEC CMP] USDT→XRP→KCS | est=-0.4665% | ideal=-0.1673% | rounded=-0.1710% | final=-0.4730%
2026/04/10 05:54:46 [EXEC CMP] USDT→DOT→KCS | est=-0.7914% | ideal=-0.4932% | rounded=-0.5000% | final=-0.8020%
2026/04/10 05:54:46 [EXEC CMP] USDT→XRP→ETH | est=-0.3274% | ideal=-0.0277% | rounded=-0.0290% | final=-0.3301%
2026/04/10 05:54:46 [EXEC CMP] USDT→XRP→KCS | est=-0.4783% | ideal=-0.1792% | rounded=-0.1830% | final=-0.4850%
2026/04/10 05:54:46 [EXEC CMP] USDT→DOT→KCS | est=-0.8032% | ideal=-0.5050% | rounded=-0.5120% | final=-0.8140%
2026/04/10 05:54:46 [EXEC CMP] USDT→DOT→KCS | est=-0.8682% | ideal=-0.5702% | rounded=-0.5790% | final=-0.8720%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.8682% | ideal=-0.5702% | rounded=-0.5790% | final=-0.8720%
2026/04/10 05:54:47 [EXEC CMP] USDT→XRP→KCS | est=-0.4665% | ideal=-0.1673% | rounded=-0.1710% | final=-0.4730%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.8564% | ideal=-0.5584% | rounded=-0.5670% | final=-0.8610%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.8641% | ideal=-0.5661% | rounded=-0.5670% | final=-0.8690%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.8759% | ideal=-0.5779% | rounded=-0.5790% | final=-0.8810%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.8759% | ideal=-0.5779% | rounded=-0.5790% | final=-0.8810%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.8109% | ideal=-0.5127% | rounded=-0.5200% | final=-0.8220%
2026/04/10 05:54:47 [EXEC CMP] USDT→XRP→KCS | est=-0.4872% | ideal=-0.1881% | rounded=-0.1920% | final=-0.4940%
2026/04/10 05:54:47 [EXEC CMP] USDT→XRP→KCS | est=-0.4754% | ideal=-0.1762% | rounded=-0.1800% | final=-0.4820%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.7991% | ideal=-0.5009% | rounded=-0.5080% | final=-0.8100%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.7991% | ideal=-0.5009% | rounded=-0.5080% | final=-0.8100%
2026/04/10 05:54:47 [EXEC CMP] USDT→XRP→KCS | est=-0.4754% | ideal=-0.1762% | rounded=-0.1800% | final=-0.4820%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.7914% | ideal=-0.4932% | rounded=-0.5000% | final=-0.8020%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.7265% | ideal=-0.4281% | rounded=-0.4320% | final=-0.7340%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.6615% | ideal=-0.3629% | rounded=-0.3650% | final=-0.6670%
2026/04/10 05:54:47 [EXEC CMP] USDT→XRP→KCS | est=-0.4754% | ideal=-0.1762% | rounded=-0.1800% | final=-0.4820%
2026/04/10 05:54:47 [EXEC CMP] USDT→DOT→KCS | est=-0.6733% | ideal=-0.3747% | rounded=-0.3770% | final=-0.6790%
2026/04/10 05:54:49 [STATS] ticks=2378 triangles_seen=3845 cand=229 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=57, no_quote_leg_3=536, no_quote_leg_2=1072, no_quote_leg_1=1945} | exec_rejects={profit_below_threshold=229}
2026/04/10 05:54:49 [EXEC CMP] USDT→DOT→KCS | est=-0.5888% | ideal=-0.2900% | rounded=-0.2980% | final=-0.6000%
2026/04/10 05:54:50 [REJECT] stage=scan reason=no_quote_leg_1 count=2000 tri=USDT->BTC->RSR
2026/04/10 05:54:51 [EXEC CMP] USDT→DOT→KCS | est=-0.5965% | ideal=-0.2977% | rounded=-0.3060% | final=-0.6080%
2026/04/10 05:54:51 [EXEC CMP] USDT→LINK→BTC | est=-0.3762% | ideal=-0.0768% | rounded=-0.0831% | final=-0.3915%
2026/04/10 05:54:51 [EXEC CMP] USDT→DOT→KCS | est=-0.6083% | ideal=-0.3096% | rounded=-0.3180% | final=-0.6200%
2026/04/10 05:54:52 [EXEC CMP] USDT→LINK→BTC | est=-0.3762% | ideal=-0.0768% | rounded=-0.0831% | final=-0.3915%
2026/04/10 05:54:52 [EXEC CMP] USDT→LINK→BTC | est=-0.3682% | ideal=-0.0687% | rounded=-0.0759% | final=-0.3771%
2026/04/10 05:54:52 [EXEC CMP] USDT→DOT→KCS | est=-0.5965% | ideal=-0.2977% | rounded=-0.3060% | final=-0.6080%
2026/04/10 05:54:52 [EXEC CMP] USDT→LINK→BTC | est=-0.3660% | ideal=-0.0665% | rounded=-0.0759% | final=-0.3771%
2026/04/10 05:54:52 [EXEC CMP] USDT→BTC→SNX | est=-1.8544% | ideal=-1.5594% | rounded=-1.5690% | final=-1.9030%
2026/04/10 05:54:54 [STATS] ticks=2515 triangles_seen=4122 cand=238 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_zero=6, max_start_lt_10=64, no_quote_leg_3=567, no_quote_leg_2=1176, no_quote_leg_1=2071} | exec_rejects={profit_below_threshold=238}
2026/04/10 05:54:54 [EXEC CMP] USDT→DOT→KCS | est=-0.6733% | ideal=-0.3747% | rounded=-0.3770% | final=-0.6790%
2026/04/10 05:54:54 [EXEC CMP] USDT→DOT→KCS | est=-0.6615% | ideal=-0.3629% | rounded=-0.3650% | final=-0.6670%
^C2026/04/10 05:54:54 [Main] shutting down...
2026/04/10 05:54:54 [KuCoin WS 0] read error: read tcp 192.168.1.66:40060->13.33.235.88:443: use of closed network connection
2026/04/10 05:54:54 [Main] exited
gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ 



