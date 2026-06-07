// Package sqlserver collects SQL Server metadata. Read-only, agentless.
// Source method: DMV / catalog views (read-only T-SQL) — editions, workload, features
package sqlserver

import "context"

// Collector reads SQL Server metadata.
type Collector struct{}

func (Collector) Name() string { return "sqlserver" }

// Collect performs a read-only read of SQL Server metadata.
func (Collector) Collect(ctx context.Context) (any, error) {
	// TODO: implement read-only collection. See docs/fingerprint-schema.md.
	return nil, nil
}
