// Package deps collects Dependencies metadata. Read-only, agentless.
// Source method: Agentless flow sampling + vCenter — host/process affinity for waves
package deps

import "context"

// Collector reads Dependencies metadata.
type Collector struct{}

func (Collector) Name() string { return "deps" }

// Collect performs a read-only read of Dependencies metadata.
func (Collector) Collect(ctx context.Context) (any, error) {
	// TODO: implement read-only collection. See docs/fingerprint-schema.md.
	return nil, nil
}
