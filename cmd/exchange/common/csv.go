package common

import (
	"encoding/csv"
	"os"
)

func SaveTrianglesCSV(path string, list []Triangle) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{
		"A", "B", "C",
		"Leg1", "Leg2", "Leg3",
	})

	for _, t := range list {
		w.Write([]string{
			t.A, t.B, t.C,
			t.Leg1, t.Leg2, t.Leg3,
		})
	}

	return nil
}
