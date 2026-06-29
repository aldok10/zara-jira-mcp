package sprint

import "math"

// Score represents a standardized 0-100 score with grading.
type Score float64

const (
	ScoreHealthy Score = 80
	ScoreFair    Score = 60
	ScoreAtRisk  Score = 40
	ScoreMin     Score = 0
	ScoreMax     Score = 100
)

// Grade returns the letter grade for the score.
func (s Score) Grade() string {
	switch {
	case s >= ScoreHealthy:
		return "A"
	case s >= ScoreFair:
		return "B"
	case s >= ScoreAtRisk:
		return "C"
	default:
		return "D"
	}
}

// Rating returns a human-readable rating.
func (s Score) Rating() string {
	switch {
	case s >= ScoreHealthy:
		return "Healthy"
	case s >= ScoreFair:
		return "Fair"
	case s >= ScoreAtRisk:
		return "At Risk"
	default:
		return "Critical"
	}
}

// IsHealthy returns true if score meets the healthy threshold.
func (s Score) IsHealthy() bool {
	return s >= ScoreHealthy
}

// WeightedScore computes a weighted contribution: score * (weight / totalWeight).
func WeightedScore(score, weight, totalWeight float64) Score {
	if totalWeight == 0 {
		return 0
	}
	return Score(math.Round(score * weight / totalWeight))
}
