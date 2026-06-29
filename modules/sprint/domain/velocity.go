package sprint

import "math"

// Velocity represents sprint velocity with trend analysis.
type Velocity float64

// Trend represents a velocity direction.
type Trend string

const (
	TrendImproving Trend = "improving"
	TrendDeclining Trend = "declining"
	TrendStable    Trend = "stable"
	TrendUnknown   Trend = "unknown"
)

// Velocities is a slice of Velocity values ordered by sprint (oldest first).
type Velocities []Velocity

// Trend computes the overall direction across the velocities.
func (v Velocities) Trend() Trend {
	if len(v) < 2 {
		return TrendUnknown
	}
	first := v[0]
	last := v[len(v)-1]
	diff := float64(last - first)

	switch {
	case diff > 0:
		return TrendImproving
	case diff < 0:
		return TrendDeclining
	default:
		return TrendStable
	}
}

// Mean returns the average velocity.
func (v Velocities) Mean() float64 {
	if len(v) == 0 {
		return 0
	}
	var sum float64
	for _, vel := range v {
		sum += float64(vel)
	}
	return sum / float64(len(v))
}

// StdDev returns the population standard deviation.
func (v Velocities) StdDev() float64 {
	if len(v) < 2 {
		return 0
	}
	mean := v.Mean()
	var sumSq float64
	for _, vel := range v {
		d := float64(vel) - mean
		sumSq += d * d
	}
	return math.Sqrt(sumSq / float64(len(v)))
}

// CV returns the coefficient of variation (StdDev/Mean) as a percentage.
// High CV (>40%) indicates velocity rollercoaster.
func (v Velocities) CV() float64 {
	mean := v.Mean()
	if mean == 0 {
		return 0
	}
	return (v.StdDev() / mean) * 100
}

// IsRollercoaster returns true if CV exceeds 40%.
func (v Velocities) IsRollercoaster() bool {
	return v.CV() > 40
}
