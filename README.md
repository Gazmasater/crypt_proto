func printProfitStats(events []Event) {
	fmt.Println("\n=== PROFIT ===")
	printOneStats("profit_pct", collect(events, func(e Event) float64 { return e.ProfitPct }))
	printOneStats("profit_usdt", collect(events, func(e Event) float64 { return e.ProfitUSDT }))
	printOneStats("volume_usdt", collect(events, func(e Event) float64 { return e.VolumeUSDT }))
	printOneStats("opportunity_strength", collect(events, func(e Event) float64 { return e.OpportunityStrength } }))


[{
	"resource": "/home/gaz358/myprog/crypt_proto/cmd/arb/metrics/metrics.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"severity": 8,
	"message": "missing ',' in argument list",
	"source": "syntax",
	"startLineNumber": 258,
	"startColumn": 111,
	"endLineNumber": 258,
	"endColumn": 111,
	"modelVersionId": 9,
	"origin": "extHost1"
}]
