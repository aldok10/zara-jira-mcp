package team

// Workload represents a single team member's workload.
type Workload struct {
	Name       string
	Assigned   int
	InProgress int
	Blocked    int
	Done       int
}

// Signal returns the overload status based on workload thresholds.
func (w Workload) Signal(avgLoad float64) string {
	switch {
	case w.Assigned > int(avgLoad*2):
		return "OVERLOADED"
	case w.InProgress >= 3:
		return "HIGH WIP"
	case w.Blocked > 0 && w.InProgress == 0:
		return "STUCK"
	case w.Assigned == 0:
		return "IDLE?"
	default:
		return "OK"
	}
}

// IsOverloaded returns true if assigned exceeds 2x the average.
func (w Workload) IsOverloaded(avgLoad float64) bool {
	return float64(w.Assigned) > avgLoad*2
}

// Workloads is a collection of Workload values.
type Workloads []Workload

// AverageLoad computes the mean assigned count across members with work.
func (ws Workloads) AverageLoad() float64 {
	if len(ws) == 0 {
		return 0
	}
	var total int
	var count int
	for _, w := range ws {
		total += w.Assigned
		if w.Assigned > 0 {
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return float64(total) / float64(count)
}

// Overloaded returns members whose assigned > 2x average.
func (ws Workloads) Overloaded() []Workload {
	avg := ws.AverageLoad()
	var out []Workload
	for _, w := range ws {
		if w.IsOverloaded(avg) {
			out = append(out, w)
		}
	}
	return out
}
