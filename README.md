if state.ProfitPct < f.cfg.MinProfitPct {
    if f.cfg.LogMode == LogDebug {
        log.Printf(
            "[EXEC REJECT] %s→%s→%s | real=%.4f%% | min=%.4f%%",
            cand.Triangle.A,
            cand.Triangle.B,
            cand.Triangle.C,
            state.ProfitPct*100,
            f.cfg.MinProfitPct*100,
        )
    }

    return ExecutableOpportunity{}, "profit_below_threshold", false
}


gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
2026/04/09 11:39:14 pprof on http://localhost:6060/debug/pprof/
2026/04/09 11:39:15 [KuCoin WS 0] connected
2026/04/09 11:39:16 [KuCoin WS 1] connected
2026/04/09 11:39:16 [KuCoin] started with 2 WS
2026/04/09 11:39:16 [Main] KuCoinCollector started
2026/04/09 11:39:16 [CALC] started | triangles indexed=121 | minVolume=10.00 USDT | minProfit=0.0000% | quoteAgeMax=400ms | logMode=debug
2026/04/09 11:39:16 [REJECT] stage=scan reason=no_quote_leg_1 count=1 tri=USDT->BTC->ETH
2026/04/09 11:39:16 [REJECT] stage=scan reason=no_quote_leg_2 count=1 tri=USDT->KCS->ETH
2026/04/09 11:39:16 [REJECT] stage=scan reason=no_quote_leg_1 count=10 tri=USDT->BTC->EUR
2026/04/09 11:39:17 [REJECT] stage=scan reason=no_quote_leg_2 count=10 tri=USDT->FET->ETH
2026/04/09 11:39:19 [REJECT] stage=scan reason=no_quote_leg_1 count=100 tri=USDT->ETH->DASH
2026/04/09 11:39:20 [REJECT] stage=scan reason=no_quote_leg_3 count=1 tri=USDT->ETC->ETH
2026/04/09 11:39:20 [EXEC REJECT] USDT→ETC→ETH | real=-0.3832% | min=0.0000%
2026/04/09 11:39:20 [REJECT] stage=exec reason=profit_below_threshold count=1 tri=USDT->ETC->ETH
2026/04/09 11:39:20 [EXEC REJECT] USDT→ETC→ETH | real=-0.3832% | min=0.0000%
2026/04/09 11:39:20 [EXEC REJECT] USDT→FET→ETH | real=-0.4377% | min=0.0000%
2026/04/09 11:39:20 [EXEC REJECT] USDT→ETC→ETH | real=-0.3832% | min=0.0000%
2026/04/09 11:39:20 [EXEC REJECT] USDT→ETC→ETH | real=-0.3832% | min=0.0000%
2026/04/09 11:39:20 [EXEC REJECT] USDT→ETC→ETH | real=-0.3832% | min=0.0000%
2026/04/09 11:39:21 [EXEC REJECT] USDT→ETC→ETH | real=-0.3832% | min=0.0000%
2026/04/09 11:39:21 [EXEC REJECT] USDT→ETC→ETH | real=-0.3832% | min=0.0000%
2026/04/09 11:39:21 [STATS] ticks=179 triangles_seen=314 cand=8 exec=0 pos=0 neg=0 logged=0 | scan_rejects={no_quote_leg_3=8, no_quote_leg_2=99, no_quote_leg_1=199} | exec_rejects={profit_below_threshold=8}
2026/04/09 11:39:21 [REJECT] stage=scan reason=no_quote_leg_2 count=100 tri=USDT->KLV->TRX
2026/04/09 11:39:22 [REJECT] stage=scan reason=no_quote_leg_3 count=10 tri=USDT->DOT->KCS
2026/04/09 11:39:25 [REJECT] stage=scan reason=max_start_lt_10 count=1 tri=USDT->KLV->TRX
2026/04/09 11:39:26 [EXEC REJECT] USDT→EUR→ETH | real=-3.9575% | min=0.0000%
2026/04/09 11:39:26 [EXEC REJECT] USDT→ETH→ZIL | real=-1.1326% | min=0.0000%
2026/04/09 11:39:26 [REJECT] stage=exec reason=profit_below_threshold count=10 tri=USDT->ETH->ZIL
2026/04/09 11:39:26 [EXEC REJECT] USDT→EUR→ETH | real=-3.9518% | min=0.0000%
2026/04/09 11:39:26 [EXEC REJECT] USDT→ETH→ZIL | real=-1.1326% | min=0.0000%
2026/04/09 11:39:26 [EXEC REJECT] USDT→ETH→ZIL | real=-1.1326% | min=0.0000%
2026/04/09 11:39:26 [STATS] ticks=528 triangles_seen=828 cand=13 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=7, no_quote_leg_3=53, no_quote_leg_2=182, no_quote_leg_1=573} | exec_rejects={profit_below_threshold=13}
2026/04/09 11:39:26 [EXEC REJECT] USDT→ETC→ETH | real=-0.3718% | min=0.0000%
2026/04/09 11:39:30 [REJECT] stage=scan reason=max_start_lt_10 count=10 tri=USDT->BTC->WAN
2026/04/09 11:39:30 [EXEC REJECT] USDT→BTC→ICP | real=-0.9200% | min=0.0000%
2026/04/09 11:39:30 [EXEC REJECT] USDT→BTC→WAN | real=-4.5966% | min=0.0000%
2026/04/09 11:39:30 [EXEC REJECT] USDT→BTC→WAN | real=-4.5966% | min=0.0000%
2026/04/09 11:39:30 [EXEC REJECT] USDT→BTC→ICP | real=-0.8900% | min=0.0000%
2026/04/09 11:39:30 [EXEC REJECT] USDT→BTC→ZEC | real=-0.5800% | min=0.0000%
2026/04/09 11:39:31 [REJECT] stage=scan reason=no_quote_leg_1 count=1000 tri=USDT->BTC->VSYS
2026/04/09 11:39:31 [STATS] ticks=1019 triangles_seen=1494 cand=19 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=14, no_quote_leg_3=80, no_quote_leg_2=295, no_quote_leg_1=1086} | exec_rejects={profit_below_threshold=19}
2026/04/09 11:39:32 [EXEC REJECT] USDT→FET→ETH | real=-0.4226% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→ENJ | real=-2.9388% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→FET→ETH | real=-0.4314% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETC→ETH | real=-0.3550% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→DAG | real=-26.8640% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→DAG | real=-26.8556% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETC→ETH | real=-0.3550% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→FET→ETH | real=-0.4314% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→DAG | real=-26.8524% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→ENJ | real=-2.9388% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→DASH | real=-0.6000% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETC→ETH | real=-0.3550% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→DASH | real=-0.6000% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→FET→ETH | real=-0.4314% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→DAG | real=-26.8608% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETC→ETH | real=-0.3550% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETH→DAG | real=-26.8524% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETC→ETH | real=-0.3550% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETC→ETH | real=-0.3550% | min=0.0000%
2026/04/09 11:39:32 [EXEC REJECT] USDT→ETC→ETH | real=-0.3550% | min=0.0000%
2026/04/09 11:39:34 [REJECT] stage=scan reason=no_quote_leg_3 count=100 tri=USDT->ETC->ETH
2026/04/09 11:39:36 [STATS] ticks=1507 triangles_seen=2193 cand=39 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=16, no_quote_leg_3=106, no_quote_leg_2=360, no_quote_leg_1=1672} | exec_rejects={profit_below_threshold=39}
2026/04/09 11:39:39 [REJECT] stage=scan reason=no_quote_leg_1 count=2000 tri=USDT->ETH->XMR
2026/04/09 11:39:41 [STATS] ticks=1976 triangles_seen=2802 cand=39 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=16, no_quote_leg_3=122, no_quote_leg_2=384, no_quote_leg_1=2241} | exec_rejects={profit_below_threshold=39}
2026/04/09 11:39:43 [EXEC REJECT] USDT→FET→ETH | real=-0.4719% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→KNC | real=-2.5800% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→KNC | real=-2.5800% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→KNC | real=-2.7310% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→ENJ | real=-3.1001% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→DAG | real=-26.8524% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→KNC | real=-2.4280% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→ENJ | real=-3.1001% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→KNC | real=-2.4280% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→DAG | real=-26.8608% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→ENJ | real=-3.0679% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→BTC→BRL | real=-1.6610% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→ENJ | real=-3.1302% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→KNC | real=-2.2830% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→KNC | real=-2.1300% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→BTC→KNC | real=-1.4010% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→DAG | real=-26.8656% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→ETH→DAG | real=-26.8572% | min=0.0000%
2026/04/09 11:39:43 [EXEC REJECT] USDT→BTC→BRL | real=-1.6610% | min=0.0000%
2026/04/09 11:39:44 [EXEC REJECT] USDT→BTC→WAN | real=-4.5966% | min=0.0000%
2026/04/09 11:39:46 [STATS] ticks=2476 triangles_seen=3568 cand=59 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=19, no_quote_leg_3=169, no_quote_leg_2=549, no_quote_leg_1=2772} | exec_rejects={profit_below_threshold=59}
2026/04/09 11:39:48 [REJECT] stage=scan reason=no_quote_leg_1 count=3000 tri=USDT->ETH->XYO
2026/04/09 11:39:49 [EXEC REJECT] USDT→BTC→WAN | real=-4.5966% | min=0.0000%
2026/04/09 11:39:49 [EXEC REJECT] USDT→BTC→DASH | real=-0.5000% | min=0.0000%
2026/04/09 11:39:49 [EXEC REJECT] USDT→BTC→BNB | real=-0.9100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→WAN | real=-4.4973% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→RUNE | real=-0.5900% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→RUNE | real=-0.5750% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→WAN | real=-4.4838% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→BNB | real=-0.9100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→ETH→ENJ | real=-3.2269% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→XRP→ETH | real=-0.3578% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→COTI | real=-0.7297% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→AVAX | real=-0.4200% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→FET→ETH | real=-0.3840% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→FET | real=-0.6190% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→ETH→XMR | real=-6.9000% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→BNB | real=-0.9100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→FET→ETH | real=-1.0934% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→HBAR | real=-0.4452% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→AVAX | real=-0.4100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→ETH→XMR | real=-6.9000% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→BNB | real=-0.9100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→HBAR | real=-0.4452% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→FET→ETH | real=-0.7660% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→VET | real=-0.9650% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→VET | real=-0.8270% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→ETH→XMR | real=-6.9000% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→BNB | real=-0.9100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→AVAX | real=-0.4100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→WAN | real=-4.4838% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→NEO | real=-1.0788% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→FET→ETH | real=-0.7485% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→LINK→BTC | real=-0.3872% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→BCH | real=-0.8000% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→AVAX | real=-0.4100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→WAN | real=-4.4973% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→BNB | real=-0.9100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→BNB | real=-0.9100% | min=0.0000%
2026/04/09 11:39:50 [EXEC REJECT] USDT→BTC→AVAX | real=-0.4200% | min=0.0000%
2026/04/09 11:39:51 [EXEC REJECT] USDT→LINK→BTC | real=-0.3872% | min=0.0000%
2026/04/09 11:39:51 [EXEC REJECT] USDT→BTC→ETC | real=-1.4956% | min=0.0000%
2026/04/09 11:39:51 [STATS] ticks=3037 triangles_seen=4504 cand=99 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=20, no_quote_leg_3=243, no_quote_leg_2=850, no_quote_leg_1=3292} | exec_rejects={profit_below_threshold=99}
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→HBAR | real=-0.4411% | min=0.0000%
2026/04/09 11:39:55 [REJECT] stage=exec reason=profit_below_threshold count=100 tri=USDT->BTC->HBAR
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→FET | real=-0.5840% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→BCH | real=-0.8000% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→DASH | real=-0.5000% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→INJ | real=-0.4100% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→LINK→BTC | real=-0.4045% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→BCH | real=-0.8000% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→LINK→BTC | real=-0.4045% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→HBAR | real=-0.4411% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→DASH | real=-0.5000% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→ICP | real=-0.9500% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→LINK→BTC | real=-0.3974% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→DASH | real=-0.6000% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→PAXG | real=-5.8000% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→CHZ | real=-0.7609% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→DASH | real=-0.6000% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→AVAX | real=-0.3900% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→LINK→BTC | real=-0.3903% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→ICP | real=-0.9200% | min=0.0000%
2026/04/09 11:39:55 [EXEC REJECT] USDT→BTC→DASH | real=-0.7000% | min=0.0000%
2026/04/09 11:39:56 [STATS] ticks=3521 triangles_seen=5186 cand=119 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=20, no_quote_leg_3=274, no_quote_leg_2=956, no_quote_leg_1=3817} | exec_rejects={profit_below_threshold=119}
2026/04/09 11:39:57 [EXEC REJECT] USDT→ETC→ETH | real=-0.4455% | min=0.0000%
2026/04/09 11:39:57 [EXEC REJECT] USDT→ETH→DAG | real=-26.8492% | min=0.0000%
2026/04/09 11:39:57 [EXEC REJECT] USDT→XLM→ETH | real=-0.3800% | min=0.0000%
2026/04/09 11:39:57 [EXEC REJECT] USDT→ETH→DASH | real=-0.6000% | min=0.0000%
2026/04/09 11:39:58 [REJECT] stage=scan reason=no_quote_leg_1 count=4000 tri=USDT->ETH->XYO
2026/04/09 11:40:00 [REJECT] stage=scan reason=no_quote_leg_2 count=1000 tri=USDT->VET->ETH
2026/04/09 11:40:01 [STATS] ticks=4041 triangles_seen=5873 cand=123 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=22, no_quote_leg_3=309, no_quote_leg_2=1008, no_quote_leg_1=4411} | exec_rejects={profit_below_threshold=123}
2026/04/09 11:40:06 [STATS] ticks=4436 triangles_seen=6368 cand=123 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=22, no_quote_leg_3=316, no_quote_leg_2=1021, no_quote_leg_1=4886} | exec_rejects={profit_below_threshold=123}
2026/04/09 11:40:07 [REJECT] stage=scan reason=no_quote_leg_1 count=5000 tri=USDT->BTC->XMR
2026/04/09 11:40:07 [EXEC REJECT] USDT→LINK→BTC | real=-0.3759% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→XTZ | real=-0.9210% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→FET | real=-0.6010% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→BCH | real=-1.2000% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→WAN | real=-5.3114% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→XMR | real=-6.9000% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→LINK→BTC | real=-0.3759% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→FET | real=-0.6010% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→FET | real=-0.6600% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→ZEC | real=-0.6400% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→LINK→BTC | real=-0.3759% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→BCH | real=-1.2000% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→AVAX | real=-0.4200% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→ZEC | real=-0.6500% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→FET→ETH | real=-0.4721% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→ETH→DAG | real=-26.8492% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→ETH→XMR | real=-6.9000% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→XRP→ETH | real=-0.3171% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→ETH→DAG | real=-26.8576% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→FET→ETH | real=-0.4612% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→AVAX | real=-0.4300% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→BCH | real=-1.2000% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→XRP→ETH | real=-0.3171% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→ETH→DAG | real=-26.8576% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→XRP→ETH | real=-0.3272% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→FET→ETH | real=-0.4712% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→ETH→DAG | real=-26.8512% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→ETH→XMR | real=-6.9000% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→BTC→BCH | real=-1.2000% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→ETH→BRL | real=-3.1210% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→FET→ETH | real=-0.4625% | min=0.0000%
2026/04/09 11:40:07 [EXEC REJECT] USDT→ETH→BRL | real=-3.1210% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→BTC→BCH | real=-1.3000% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→BRL | real=-3.1190% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→BRL | real=-3.1210% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→DAG | real=-26.8262% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→DAG | real=-26.8492% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→DASH | real=-0.6000% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→ENJ | real=-3.1969% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→DAG | real=-26.8460% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→ENJ | real=-3.1969% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→DASH | real=-0.6000% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→DASH | real=-0.5000% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→DAG | real=-26.8377% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→ENJ | real=-3.1324% | min=0.0000%
2026/04/09 11:40:08 [EXEC REJECT] USDT→ETH→ENJ | real=-3.1969% | min=0.0000%
2026/04/09 11:40:09 [EXEC REJECT] USDT→ETC→ETH | real=-0.3991% | min=0.0000%
2026/04/09 11:40:09 [EXEC REJECT] USDT→ETH→ENJ | real=-3.1324% | min=0.0000%
2026/04/09 11:40:09 [EXEC REJECT] USDT→ETH→TEL | real=-2.2276% | min=0.0000%
2026/04/09 11:40:09 [EXEC REJECT] USDT→ETH→BRL | real=-3.0990% | min=0.0000%
2026/04/09 11:40:10 [EXEC REJECT] USDT→ETH→BRL | real=-3.1010% | min=0.0000%
2026/04/09 11:40:11 [STATS] ticks=5036 triangles_seen=7370 cand=174 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=28, no_quote_leg_3=405, no_quote_leg_2=1241, no_quote_leg_1=5522} | exec_rejects={profit_below_threshold=174}
2026/04/09 11:40:16 [REJECT] stage=scan reason=no_quote_leg_1 count=6000 tri=USDT->BTC->KAS
2026/04/09 11:40:16 [STATS] ticks=5493 triangles_seen=7993 cand=174 exec=0 pos=0 neg=0 logged=0 | scan_rejects={max_start_lt_10=28, no_quote_leg_3=444, no_quote_leg_2=1288, no_quote_leg_1=6059} | exec_rejects={profit_below_threshold=174}
