// Package collectors defines the contract every source-workload collector
// implements. ALL collectors MUST be strictly read-only and agentless.
package collectors

import "context"

// Collector reads metadata from one source-workload type.
type Collector interface {
	// Name is the collector identifier, e.g. "vmware", "sqlserver".
	Name() string
	// Collect performs a read-only metadata read. It must never write to,
	// or install anything on, the target.
	Collect(ctx context.Context) (any, error)
}
