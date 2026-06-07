// Package vmware collects VMware vSphere metadata. Read-only, agentless.
// Source method: vCenter REST/SOAP API (read-only) — inventory + 72h performance
package vmware

import "context"

// Collector reads VMware vSphere metadata.
type Collector struct{}

func (Collector) Name() string { return "vmware" }

// Collect performs a read-only read of VMware vSphere metadata.
func (Collector) Collect(ctx context.Context) (any, error) {
	// TODO: implement read-only collection. See docs/fingerprint-schema.md.
	return nil, nil
}
