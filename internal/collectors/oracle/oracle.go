// Package oracle collects Oracle metadata. Read-only, agentless.
// Source method: v$/dictionary + AWR/ASH if Diagnostic Pack else Statspack fallback
package oracle

import "context"

// Collector reads Oracle metadata.
type Collector struct{}

func (Collector) Name() string { return "oracle" }

// Collect performs a read-only read of Oracle metadata.
func (Collector) Collect(ctx context.Context) (any, error) {
	// TODO: implement read-only collection. See docs/fingerprint-schema.md.
	return nil, nil
}
