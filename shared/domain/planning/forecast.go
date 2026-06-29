package planning

import (
	"math"
	"math/rand"
	"sort"
)

// ForecastResult holds Monte Carlo simulation results.
type ForecastResult struct {
	Simulations int
	Throughput  []float64
	Remaining   int
	Percentiles map[int]float64 // e.g. {50: 5.2, 70: 7.1, 85: 9.3, 95: 12.0}
	MeanSprints float64
	MinSprints  int
	MaxSprints  int
}

// Forecast runs Monte Carlo simulations to predict completion sprints.
// throughput: historical completed items per sprint.
// remaining: number of items remaining.
// simulations: number of simulations to run (default 10000 if 0).
//
//nolint:gosec // Monte Carlo uses math/rand, not crypto-grade randomness
func Forecast(throughput []float64, remaining int, simulations int) ForecastResult {
	if len(throughput) == 0 || remaining <= 0 {
		return ForecastResult{Remaining: remaining}
	}
	if simulations <= 0 {
		simulations = 10000
	}

	results := make([]int, simulations)
	for i := range results {
		items := remaining
		sprints := 0
		for items > 0 {
			idx := rand.Intn(len(throughput))
			items -= int(math.Round(throughput[idx]))
			sprints++
			if sprints > 1000 {
				break // safety valve
			}
		}
		results[i] = sprints
	}

	sort.Ints(results)

	var sum int
	for _, r := range results {
		sum += r
	}

	percentiles := map[int]float64{
		50: percentile(results, 0.50),
		70: percentile(results, 0.70),
		85: percentile(results, 0.85),
		95: percentile(results, 0.95),
	}

	return ForecastResult{
		Simulations: simulations,
		Throughput:  throughput,
		Remaining:   remaining,
		Percentiles: percentiles,
		MeanSprints: math.Round(float64(sum)/float64(simulations)*100) / 100,
		MinSprints:  results[0],
		MaxSprints:  results[len(results)-1],
	}
}

func percentile(sorted []int, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	pos := p * float64(len(sorted)-1)
	low := int(math.Floor(pos))
	high := int(math.Ceil(pos))
	if low == high {
		return float64(sorted[low])
	}
	return float64(sorted[low])*(1-(pos-float64(low))) + float64(sorted[high])*(pos-float64(low))
}
