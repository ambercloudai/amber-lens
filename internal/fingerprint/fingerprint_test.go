package fingerprint

import (
	"encoding/json"
	"os"
	"testing"
)

// TestExampleRoundTrip ensures the canonical example unmarshals into the Go
// types and satisfies the structural invariants. This keeps the schema, the
// example, and the Go contract in lockstep (and runs in CI).
func TestExampleRoundTrip(t *testing.T) {
	b, err := os.ReadFile("../../pkg/schema/examples/fingerprint.example.json")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	var fp Fingerprint
	if err := json.Unmarshal(b, &fp); err != nil {
		t.Fatalf("unmarshal example: %v", err)
	}
	if err := fp.Validate(); err != nil {
		t.Fatalf("validate example: %v", err)
	}
	if fp.Schema != SchemaID {
		t.Errorf("schema = %q, want %q", fp.Schema, SchemaID)
	}
	if len(fp.Workloads) == 0 {
		t.Fatal("expected workloads in example")
	}
	// Spot-check the Oracle options fact survived the round trip.
	var sawOracleOption bool
	for _, w := range fp.Workloads {
		if w.Kind == KindDatabase && w.Database != nil && w.Database.Engine == "oracle" {
			for _, o := range w.Database.Options {
				if o.Name == "Partitioning" && o.InUse {
					sawOracleOption = true
				}
			}
		}
	}
	if !sawOracleOption {
		t.Error("expected Oracle Partitioning option in example")
	}
}

func TestValidateRejectsBadKind(t *testing.T) {
	fp := New()
	fp.Lens = Lens{Version: "x", Mode: "read-only"}
	fp.Tokenization = Tokenization{MapLocation: "on-prem"}
	fp.Workloads = []Workload{{Token: "t1", Kind: KindDatabase}} // missing Database block
	if err := fp.Validate(); err == nil {
		t.Fatal("expected validation error for database kind without database block")
	}
}
