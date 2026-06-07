// Package sampling drives the 72-hour, 5-minute-interval capture that yields
// ~95% statistical confidence (three daily cycles + maintenance windows)
// without a 30-day agent rollout.
package sampling

import "time"

const (
	Interval = 5 * time.Minute
	Window   = 72 * time.Hour
)

// TODO: scheduler that invokes performance collectors on Interval over Window.
