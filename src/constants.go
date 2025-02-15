package src

import (
	"sort"

	"gonum.org/v1/gonum/stat"
)

type StatOperation string

const (
	NotSet   StatOperation = ""
	OpMin    StatOperation = "MIN"
	OpMax    StatOperation = "MAX"
	OpMean   StatOperation = "MEAN"
	OpMedian StatOperation = "MEDIAN"
	OpP95    StatOperation = "P95"
	OpP99    StatOperation = "P99"
	OpP90    StatOperation = "P90"
)

type StatFunc func([]float64) float64

// StatsFuncs maps each StatOperation to a corresponding function
var StatsFuncs = map[StatOperation]StatFunc{
	OpMin: func(data []float64) float64 {
		if len(data) == 0 {
			return 0
		}
		min := data[0]
		for _, v := range data {
			if v < min {
				min = v
			}
		}
		return min
	},
	OpMax: func(data []float64) float64 {
		if len(data) == 0 {
			return 0
		}
		max := data[0]
		for _, v := range data {
			if v > max {
				max = v
			}
		}
		return max
	},
	OpMean:   func(data []float64) float64 { return stat.Mean(data, nil) },
	OpMedian: func(data []float64) float64 { return stat.Quantile(0.5, stat.Empirical, sortedCopy(data), nil) },
	OpP90:    func(data []float64) float64 { return stat.Quantile(0.90, stat.Empirical, sortedCopy(data), nil) },
	OpP95:    func(data []float64) float64 { return stat.Quantile(0.95, stat.Empirical, sortedCopy(data), nil) },
	OpP99:    func(data []float64) float64 { return stat.Quantile(0.99, stat.Empirical, sortedCopy(data), nil) },
}

// sortedCopy returns a sorted copy of the input slice to avoid modifying the original data
func sortedCopy(data []float64) []float64 {
	sorted := append([]float64(nil), data...)
	sort.Float64s(sorted)
	return sorted
}
