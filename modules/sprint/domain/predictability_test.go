package sprint

import (
	"testing"
)

func TestPredictability_Empty(t *testing.T) {
	result := Predictability(nil)
	if result.Rating != "UNKNOWN" {
		t.Errorf("Predictability(nil).Rating = %q, want UNKNOWN", result.Rating)
	}
}

func TestPredictability_High(t *testing.T) {
	result := Predictability([]float64{85, 90, 88, 92})
	if result.Rating != "HIGH" {
		t.Errorf("Predictability(high).Rating = %q, want HIGH", result.Rating)
	}
	if result.Average < 88 || result.Average > 89 {
		t.Errorf("Predictability(high).Average = %v, want ~88.75", result.Average)
	}
}

func TestPredictability_Medium(t *testing.T) {
	result := Predictability([]float64{70, 85, 60, 90})
	if result.Rating != "MEDIUM" {
		t.Errorf("Predictability(medium).Rating = %q, want MEDIUM", result.Rating)
	}
}

func TestPredictability_Low(t *testing.T) {
	result := Predictability([]float64{90, 40, 85, 30, 95})
	if result.Rating != "LOW" {
		t.Errorf("Predictability(low).Rating = %q, want LOW", result.Rating)
	}
}

func TestPredictability_SingleValue(t *testing.T) {
	result := Predictability([]float64{80})
	if result.Rating != "HIGH" {
		t.Errorf("Predictability(single).Rating = %q, want HIGH", result.Rating)
	}
	if result.Average != 80 {
		t.Errorf("Predictability(single).Average = %v, want 80", result.Average)
	}
}
