package sprint

import (
	"testing"
)

func TestVelocities_Trend(t *testing.T) {
	tests := []struct {
		name string
		v    Velocities
		want Trend
	}{
		{name: "empty", v: Velocities{}, want: TrendUnknown},
		{name: "single", v: Velocities{10}, want: TrendUnknown},
		{name: "improving", v: Velocities{10, 20, 30}, want: TrendImproving},
		{name: "declining", v: Velocities{30, 20, 10}, want: TrendDeclining},
		{name: "stable", v: Velocities{20, 20, 20}, want: TrendStable},
		{name: "two improving", v: Velocities{10, 15}, want: TrendImproving},
		{name: "two declining", v: Velocities{15, 10}, want: TrendDeclining},
		{name: "fluctuating net improving", v: Velocities{10, 5, 20}, want: TrendImproving},
		{name: "fluctuating net declining", v: Velocities{20, 5, 10}, want: TrendDeclining},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Trend(); got != tt.want {
				t.Errorf("Velocities(%v).Trend() = %q, want %q", tt.v, got, tt.want)
			}
		})
	}
}

func TestVelocities_Mean(t *testing.T) {
	tests := []struct {
		name string
		v    Velocities
		want float64
	}{
		{name: "empty", v: Velocities{}, want: 0},
		{name: "single", v: Velocities{10}, want: 10},
		{name: "multiple", v: Velocities{10, 20, 30}, want: 20},
		{name: "with fractions", v: Velocities{15, 25}, want: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Mean(); got != tt.want {
				t.Errorf("Velocities(%v).Mean() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVelocities_StdDev(t *testing.T) {
	tests := []struct {
		name string
		v    Velocities
		want float64
	}{
		{name: "empty", v: Velocities{}, want: 0},
		{name: "single", v: Velocities{10}, want: 0},
		{name: "identical", v: Velocities{10, 10, 10}, want: 0},
		{name: "varying", v: Velocities{10, 20}, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.StdDev(); got != tt.want {
				t.Errorf("Velocities(%v).StdDev() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVelocities_CV(t *testing.T) {
	tests := []struct {
		name string
		v    Velocities
		want float64
	}{
		{name: "empty", v: Velocities{}, want: 0},
		{name: "single", v: Velocities{10}, want: 0},
		{name: "identical", v: Velocities{10, 10, 10}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.CV(); got != tt.want {
				t.Errorf("Velocities(%v).CV() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVelocities_IsRollercoaster(t *testing.T) {
	tests := []struct {
		name string
		v    Velocities
		want bool
	}{
		{name: "empty", v: Velocities{}, want: false},
		{name: "stable", v: Velocities{10, 10, 10}, want: false},
		{name: "small variation", v: Velocities{20, 22, 19}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsRollercoaster(); got != tt.want {
				t.Errorf("Velocities(%v).IsRollercoaster() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}
