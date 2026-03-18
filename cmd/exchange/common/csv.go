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

	if err := w.Write([]string{
		"A", "B", "C",
		"Leg1", "Step1", "MinQty1", "MinNotional1",
		"Leg2", "Step2", "MinQty2", "MinNotional2",
		"Leg3", "Step3", "MinQty3", "MinNotional3",
	}); err != nil {
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
