// Package osmeta collects OS / Host metadata. Read-only, agentless.
// Source method: WMI / PowerShell / SSH (agentless) — version, EOL, right-sizing
package osmeta

import "context"

// Collector reads OS / Host metadata.
type Collector struct{}

func (Collector) Name() string { return "osmeta" }

// Collect performs a read-only read of OS / Host metadata.
func (Collector) Collect(ctx context.Context) (any, error) {
	// TODO: implement read-only collection. See docs/fingerprint-schema.md.
	return nil, nil
}
