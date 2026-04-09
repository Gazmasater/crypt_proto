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
