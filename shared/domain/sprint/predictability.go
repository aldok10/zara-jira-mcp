package sprint

import "math"

// PredictabilityResult holds the predictability analysis for a set of sprints.
type PredictabilityResult struct {
	Scores         []float64 // predictability scores per sprint (0-100)
	Average        float64
	StdDev         float64
	CoeffVariation float64 // lower is better
	Rating         string  // HIGH / MEDIUM / LOW
}

// Predictability computes predictability metrics from a slice of completion rates.
// Each completionRate is 0-100 percentage of items completed in a sprint.
func Predictability(completionRates []float64) PredictabilityResult {
	if len(completionRates) == 0 {
		return PredictabilityResult{Rating: "UNKNOWN"}
	}

	var sum float64
	for _, r := range completionRates {
		sum += r
	}
	avg := sum / float64(len(completionRates))

	var sumSq float64
	for _, r := range completionRates {
		d := r - avg
		sumSq += d * d
	}
	stdDev := math.Sqrt(sumSq / float64(len(completionRates)))

	cv := 0.0
	if avg > 0 {
		cv = (stdDev / avg) * 100
	}

	rating := "HIGH"
	if cv > 30 {
		rating = "LOW"
	} else if cv > 15 {
		rating = "MEDIUM"
	}

	return PredictabilityResult{
		Scores:         completionRates,
		Average:        math.Round(avg*100) / 100,
		StdDev:         math.Round(stdDev*100) / 100,
		CoeffVariation: math.Round(cv*100) / 100,
		Rating:         rating,
	}
}
