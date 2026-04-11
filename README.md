func printProfitStats(events []Event) {
	fmt.Println("\n=== PROFIT ===")
	printOneStats("profit_pct", collect(events, func(e Event) float64 { return e.ProfitPct }))
	printOneStats("profit_usdt", collect(events, func(e Event) float64 { return e.ProfitUSDT }))
	printOneStats("volume_usdt", collect(events, func(e Event) float64 { return e.VolumeUSDT }))
	printOneStats("opportunity_strength", collect(events, func(e Event) float64 { return e.OpportunityStrength }))
}
