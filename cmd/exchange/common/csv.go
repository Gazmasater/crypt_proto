package common

import (
	"encoding/csv"
	"os"
	"strconv"
)

func SaveTrianglesCSV(path string, list []Triangle) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"A", "B", "C",
		"Leg1", "Step1", "MinQty1", "MinNotional1",
		"Leg2", "Step2", "MinQty2", "MinNotional2",
		"Leg3", "Step3", "MinQty3", "MinNotional3",

		"Leg1Symbol", "Leg1Side", "Leg1Base", "Leg1Quote", "Leg1QtyStep", "Leg1QuoteStep", "Leg1PriceStep", "Leg1MinQty", "Leg1MinQuote", "Leg1MinNotional",
		"Leg2Symbol", "Leg2Side", "Leg2Base", "Leg2Quote", "Leg2QtyStep", "Leg2QuoteStep", "Leg2PriceStep", "Leg2MinQty", "Leg2MinQuote", "Leg2MinNotional",
		"Leg3Symbol", "Leg3Side", "Leg3Base", "Leg3Quote", "Leg3QtyStep", "Leg3QuoteStep", "Leg3PriceStep", "Leg3MinQty", "Leg3MinQuote", "Leg3MinNotional",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, t := range list {
		row := []string{
			t.A,
			t.B,
			t.C,

			t.Leg1,
			ff(t.Step1),
			ff(t.MinQty1),
			ff(t.MinNotional1),

			t.Leg2,
			ff(t.Step2),
			ff(t.MinQty2),
			ff(t.MinNotional2),

			t.Leg3,
			ff(t.Step3),
			ff(t.MinQty3),
			ff(t.MinNotional3),

			t.Leg1Symbol,
			t.Leg1Side,
			t.Leg1Base,
			t.Leg1Quote,
			ff(t.Leg1QtyStep),
			ff(t.Leg1QuoteStep),
			ff(t.Leg1PriceStep),
			ff(t.Leg1MinQty),
			ff(t.Leg1MinQuote),
			ff(t.Leg1MinNotional),

			t.Leg2Symbol,
			t.Leg2Side,
			t.Leg2Base,
			t.Leg2Quote,
			ff(t.Leg2QtyStep),
			ff(t.Leg2QuoteStep),
			ff(t.Leg2PriceStep),
			ff(t.Leg2MinQty),
			ff(t.Leg2MinQuote),
			ff(t.Leg2MinNotional),

			t.Leg3Symbol,
			t.Leg3Side,
			t.Leg3Base,
			t.Leg3Quote,
			ff(t.Leg3QtyStep),
			ff(t.Leg3QuoteStep),
			ff(t.Leg3PriceStep),
			ff(t.Leg3MinQty),
			ff(t.Leg3MinQuote),
			ff(t.Leg3MinNotional),
		}

		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

func ff(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}
