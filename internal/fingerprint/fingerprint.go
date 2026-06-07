// Package fingerprint defines the Workload Fingerprint — the versioned,
// destination-neutral output schema of a scan. The v1 core (compute + SQL
// Server + Oracle + OS + dependencies) is the artifact to FREEZE by July 1.
// Only salted, deterministic tokens appear here; raw names stay on-prem.
package fingerprint

// SchemaVersion is the fingerprint schema version. Bump on breaking changes.
const SchemaVersion = "v1"

// Fingerprint is the anonymized, structured result of a discovery scan.
type Fingerprint struct {
	Schema     string   `json:"schema"`      // == SchemaVersion
	ScanID     string   `json:"scanId"`      // opaque
	WindowHrs  int      `json:"windowHours"` // e.g. 72
	Confidence float64  `json:"confidence"`  // e.g. 0.95
	Workloads  []Item   `json:"workloads"`
}

// Item is one tokenized source workload (token, not real name).
type Item struct {
	Token    string         `json:"token"`    // e.g. prod-cl-01
	Kind     string         `json:"kind"`     // vm | sqlserver | oracle | host | ...
	Metadata map[string]any `json:"metadata"` // collector-specific, schema-governed
}

// Build assembles a Fingerprint from collected, tokenized findings.
func Build( /* TODO: tokenized findings */ ) (*Fingerprint, error) {
	// TODO: validate against pkg/schema/fingerprint.schema.json before return.
	return &Fingerprint{Schema: SchemaVersion}, nil
}
