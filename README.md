Найди текущий код:

if state.ProfitPct < f.cfg.MinProfitPct {
    return ExecutableOpportunity{}, "profit_below_threshold", false
}

И замени на:

if state.ProfitPct < f.cfg.MinProfitPct {
    return ExecutableOpportunity{},
        fmt.Sprintf(
            "profit_below_threshold(real=%.4f%%, min=%.4f%%)",
            state.ProfitPct*100,
            f.cfg.MinProfitPct*100,
        ),
        false
}


