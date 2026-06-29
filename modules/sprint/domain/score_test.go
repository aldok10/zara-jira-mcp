package sprint

import (
	"testing"
)

func TestScore_Grade(t *testing.T) {
	tests := []struct {
		name  string
		score Score
		want  string
	}{
		{name: "perfect", score: 100, want: "A"},
		{name: "just healthy", score: 80, want: "A"},
		{name: "near healthy", score: 79, want: "B"},
		{name: "just fair", score: 60, want: "B"},
		{name: "near fair", score: 59, want: "C"},
		{name: "just at risk", score: 40, want: "C"},
		{name: "near at risk", score: 39, want: "D"},
		{name: "zero", score: 0, want: "D"},
		{name: "very high", score: 150, want: "A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.score.Grade(); got != tt.want {
				t.Errorf("Score(%v).Grade() = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}

func TestScore_Rating(t *testing.T) {
	tests := []struct {
		name  string
		score Score
		want  string
	}{
		{name: "healthy", score: 85, want: "Healthy"},
		{name: "fair", score: 65, want: "Fair"},
		{name: "at risk", score: 45, want: "At Risk"},
		{name: "critical", score: 20, want: "Critical"},
		{name: "boundary healthy", score: 80, want: "Healthy"},
		{name: "boundary fair", score: 60, want: "Fair"},
		{name: "boundary at risk", score: 40, want: "At Risk"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.score.Rating(); got != tt.want {
				t.Errorf("Score(%v).Rating() = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}

func TestScore_IsHealthy(t *testing.T) {
	tests := []struct {
		name  string
		score Score
		want  bool
	}{
		{name: "healthy", score: 100, want: true},
		{name: "at threshold", score: 80, want: true},
		{name: "just below", score: 79, want: false},
		{name: "zero", score: 0, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.score.IsHealthy(); got != tt.want {
				t.Errorf("Score(%v).IsHealthy() = %v, want %v", tt.score, got, tt.want)
			}
		})
	}
}

func TestWeightedScore(t *testing.T) {
	tests := []struct {
		name        string
		score       float64
		weight      float64
		totalWeight float64
		want        Score
	}{
		{name: "full weight", score: 100, weight: 10, totalWeight: 10, want: 100},
		{name: "half weight", score: 80, weight: 5, totalWeight: 10, want: 40},
		{name: "zero total", score: 100, weight: 10, totalWeight: 0, want: 0},
		{name: "quarter weight", score: 60, weight: 2, totalWeight: 8, want: 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WeightedScore(tt.score, tt.weight, tt.totalWeight); got != tt.want {
				t.Errorf("WeightedScore(%v, %v, %v) = %v, want %v", tt.score, tt.weight, tt.totalWeight, got, tt.want)
			}
		})
	}
}
