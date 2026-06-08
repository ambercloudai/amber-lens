// Package fingerprint defines the Workload Fingerprint — the versioned,
// anonymized, destination-neutral output of an Amber Lens scan and the shared
// contract between the free Lens and the paid Amber Navigator.
//
// Design rules (enforced by the contract, not just convention):
//   - FACTS ONLY. The fingerprint carries observed facts (inventory, config,
//     summarized performance, DB editions/options, OS/EOL, dependencies). It
//     carries NO risk score, NO TCO, NO cost — those are derived in Navigator
//     (the paid layer). This keeps the free/paid line in the data model.
//   - TOKENS ONLY. Workloads are identified by salted tokens; the
//     token->real-name map never leaves the customer network.
//   - DESTINATION-NEUTRAL. No target-cloud assumptions are encoded.
//   - SUMMARIZED, not raw. Performance is window-summarized (avg/p50/p95/peak),
//     never raw time series — compact and privacy-preserving.
//   - EXTENSIBLE. Additive fields bump the minor version; collector profiles
//     extend a workload via Extensions without breaking the v1 core.
//
// Canonical JSON Schema: pkg/schema/fingerprint.schema.json. v1 core freezes 2026-07-01.
package fingerprint

import (
	"errors"
	"fmt"
	"time"
)

// SchemaID and SchemaVersion identify the contract.
const (
	SchemaID      = "amber.fingerprint/v1"
	SchemaVersion = "1.0.0"
)

// Kind enumerates the v1 core workload types.
type Kind string

const (
	KindVM       Kind = "vm"
	KindHost     Kind = "host"
	KindCluster  Kind = "cluster"
	KindDatabase Kind = "database"
)

// Fingerprint is the top-level scan output.
type Fingerprint struct {
	Schema           string            `json:"schema"`
	SchemaVersion    string            `json:"schemaVersion"`
	ScanID           string            `json:"scanId"`
	GeneratedAt      time.Time         `json:"generatedAt"`
	Lens             Lens              `json:"lens"`
	Scan             Scan              `json:"scan"`
	Tokenization     Tokenization      `json:"tokenization"`
	Environment      *Environment      `json:"environment,omitempty"`
	Summary          Summary           `json:"summary"`
	Workloads        []Workload        `json:"workloads"`
	Dependencies     []Dependency      `json:"dependencies,omitempty"`
	CollectionIssues []CollectionIssue `json:"collectionIssues,omitempty"`
	Recommendations  []Recommendation  `json:"recommendations,omitempty"`
}

// Lens records collector provenance.
type Lens struct {
	Version  string `json:"version"`
	Mode     string `json:"mode"` // always "read-only"
	BuildSha string `json:"buildSha,omitempty"`
}

// Scan records the sampling window and statistical provenance.
type Scan struct {
	Window                     Window  `json:"window"`
	SampleIntervalSeconds      int     `json:"sampleIntervalSeconds"`
	SamplesCollected           int     `json:"samplesCollected,omitempty"`
	DailyCyclesCaptured        int     `json:"dailyCyclesCaptured,omitempty"`
	MaintenanceWindowsCaptured int     `json:"maintenanceWindowsCaptured,omitempty"`
	Confidence                 float64 `json:"confidence"`
}

// Window is the scan time window.
type Window struct {
	StartedAt     time.Time `json:"startedAt"`
	EndedAt       time.Time `json:"endedAt"`
	DurationHours float64   `json:"durationHours"`
}

// Tokenization documents how identifiers were anonymized (no secrets included).
type Tokenization struct {
	Scheme      string `json:"scheme"`
	Algorithm   string `json:"algorithm"`
	SaltID      string `json:"saltId,omitempty"`
	MapLocation string `json:"mapLocation"` // always "on-prem"
}

// Environment holds anonymized site/datacenter tokens.
type Environment struct {
	SiteToken       string `json:"siteToken,omitempty"`
	DatacenterToken string `json:"datacenterToken,omitempty"`
}

// Summary holds estate aggregates (facts).
type Summary struct {
	Counts Counts `json:"counts"`
	Totals Totals `json:"totals"`
}

// Counts are estate cardinalities.
type Counts struct {
	VMs          int `json:"vms,omitempty"`
	Hosts        int `json:"hosts,omitempty"`
	Clusters     int `json:"clusters,omitempty"`
	Databases    int `json:"databases,omitempty"`
	Dependencies int `json:"dependencies,omitempty"`
}

// Totals are estate-wide capacity sums.
type Totals struct {
	VCPU                  int     `json:"vcpu,omitempty"`
	MemoryGiB             float64 `json:"memoryGiB,omitempty"`
	StorageProvisionedGiB float64 `json:"storageProvisionedGiB,omitempty"`
	StorageUsedGiB        float64 `json:"storageUsedGiB,omitempty"`
}

// PerfStat is a window-summarized metric (no raw time series).
type PerfStat struct {
	Avg  float64  `json:"avg"`
	P50  *float64 `json:"p50,omitempty"`
	P95  *float64 `json:"p95,omitempty"`
	Peak float64  `json:"peak"`
	Unit string   `json:"unit"`
}

// Performance bundles per-workload summarized metrics.
type Performance struct {
	CPUPct             *PerfStat `json:"cpuPct,omitempty"`
	MemoryPct          *PerfStat `json:"memoryPct,omitempty"`
	DiskIops           *PerfStat `json:"diskIops,omitempty"`
	DiskThroughputMBps *PerfStat `json:"diskThroughputMBps,omitempty"`
	NetworkMbps        *PerfStat `json:"networkMbps,omitempty"`
}

// OS describes the operating system and its licensing facts.
type OS struct {
	Family        string       `json:"family"`
	Distribution  string       `json:"distribution,omitempty"`
	Version       string       `json:"version,omitempty"`
	EndOfLife     *bool        `json:"endOfLife,omitempty"`
	EndOfLifeDate string       `json:"endOfLifeDate,omitempty"`
	Licensing     *OSLicensing `json:"licensing,omitempty"`
}

// OSLicensing holds OS licensing facts (e.g., Windows Server).
type OSLicensing struct {
	Product          string `json:"product,omitempty"`
	Edition          string `json:"edition,omitempty"`
	LicensableCores  int    `json:"licensableCores,omitempty"`
	SoftwareAssurance string `json:"softwareAssurance,omitempty"` // yes | no | unknown
}

// Compute holds allocated capacity for a vm/host.
type Compute struct {
	VCPU                  int     `json:"vcpu,omitempty"`
	MemoryGiB             float64 `json:"memoryGiB,omitempty"`
	ProvisionedStorageGiB float64 `json:"provisionedStorageGiB,omitempty"`
	UsedStorageGiB        float64 `json:"usedStorageGiB,omitempty"`
}

// VM holds VM-specific facts.
type VM struct {
	Hypervisor   string  `json:"hypervisor"`
	ClusterToken string  `json:"clusterToken,omitempty"`
	PowerState   string  `json:"powerState,omitempty"`
	UptimeHours  float64 `json:"uptimeHours,omitempty"`
}

// Host holds hypervisor-host facts.
type Host struct {
	Hypervisor        string `json:"hypervisor,omitempty"`
	HypervisorVersion string `json:"hypervisorVersion,omitempty"`
	ClusterToken      string `json:"clusterToken,omitempty"`
	CPUSockets        int    `json:"cpuSockets,omitempty"`
	PhysicalCores     int    `json:"physicalCores,omitempty"`
}

// Cluster holds cluster facts including hypervisor licensing.
type Cluster struct {
	Hypervisor  string              `json:"hypervisor,omitempty"`
	Version     string              `json:"version,omitempty"`
	HostCount   int                 `json:"hostCount,omitempty"`
	DRSEnabled  *bool               `json:"drsEnabled,omitempty"`
	HAEnabled   *bool               `json:"haEnabled,omitempty"`
	VSANEnabled *bool               `json:"vsanEnabled,omitempty"`
	NSXEnabled  *bool               `json:"nsxEnabled,omitempty"`
	Licensing   []ClusterLicense    `json:"licensing,omitempty"`
}

// ClusterLicense is one hypervisor license fact.
type ClusterLicense struct {
	Product  string `json:"product,omitempty"`
	Edition  string `json:"edition,omitempty"`
	CPUCount int    `json:"cpuCount,omitempty"`
}

// Database holds DB engine facts incl. licensing-relevant options.
type Database struct {
	Engine          string         `json:"engine"`
	Edition         string         `json:"edition,omitempty"`
	Version         string         `json:"version,omitempty"`
	SizeGiB         float64        `json:"sizeGiB,omitempty"`
	InstanceCount   int            `json:"instanceCount,omitempty"`
	LicensableCores int            `json:"licensableCores,omitempty"`
	HA              *DatabaseHA    `json:"ha,omitempty"`
	Features        []string       `json:"features,omitempty"`
	Options         []DBOption     `json:"options,omitempty"`
	Workload        *Performance   `json:"workload,omitempty"`
	Collection      *DBCollection  `json:"collection,omitempty"`
}

// DatabaseHA describes high-availability topology.
type DatabaseHA struct {
	Type     string `json:"type,omitempty"`
	Replicas int    `json:"replicas,omitempty"`
}

// DBOption is a licensing-relevant option/pack and whether it is in use.
type DBOption struct {
	Name  string `json:"name"`
	InUse bool   `json:"inUse"`
}

// DBCollection records how DB metadata was read (transparency/provenance).
type DBCollection struct {
	Method                 []string `json:"method,omitempty"`
	DiagnosticPackLicensed *bool    `json:"diagnosticPackLicensed,omitempty"`
	FallbackUsed           *bool    `json:"fallbackUsed,omitempty"`
}

// Workload is one tokenized source workload. The sub-object matching Kind is required.
type Workload struct {
	Token       string         `json:"token"`
	Kind        Kind           `json:"kind"`
	OS          *OS            `json:"os,omitempty"`
	Compute     *Compute       `json:"compute,omitempty"`
	Performance *Performance   `json:"performance,omitempty"`
	VM          *VM            `json:"vm,omitempty"`
	Host        *Host          `json:"host,omitempty"`
	Cluster     *Cluster       `json:"cluster,omitempty"`
	Database    *Database      `json:"database,omitempty"`
	Extensions  map[string]any `json:"extensions,omitempty"`
}

// Dependency is a relationship edge between two tokens.
type Dependency struct {
	FromToken    string `json:"fromToken"`
	ToToken      string `json:"toToken"`
	Kind         string `json:"kind"`
	Port         int    `json:"port,omitempty"`
	Protocol     string `json:"protocol,omitempty"`
	Observations int    `json:"observations,omitempty"`
}

// CollectionIssue records a non-fatal collection note (e.g., a fallback).
type CollectionIssue struct {
	Code           string   `json:"code"`
	Severity       string   `json:"severity"`
	Message        string   `json:"message"`
	AffectedTokens []string `json:"affectedTokens,omitempty"`
}

// Recommendation is a preliminary, advisory, NON-monetary recommendation — the
// free hook. Navigator later validates it (Status -> "validated") and attaches
// TCO and reverse bids, which are NOT stored in the fingerprint.
type Recommendation struct {
	ID            string         `json:"id"`
	Action        string         `json:"action"`
	AppliesTo     []string       `json:"appliesTo"`
	Rationale     string         `json:"rationale,omitempty"`
	Evidence      []string       `json:"evidence,omitempty"`
	TargetProfile map[string]any `json:"targetProfile,omitempty"`
	Impact        *RecImpact     `json:"impact,omitempty"`
	Confidence    string         `json:"confidence,omitempty"`
	Status        string         `json:"status"`
	GeneratedBy   string         `json:"generatedBy,omitempty"`
}

// RecImpact is a directional, non-monetary impact estimate — ranges/counts only,
// never cost (TCO is the paid layer).
type RecImpact struct {
	Metric string   `json:"metric"`
	Low    *float64 `json:"low,omitempty"`
	High   *float64 `json:"high,omitempty"`
	Unit   string   `json:"unit,omitempty"`
}

// New returns an empty Fingerprint stamped with the current schema id/version.
func New() *Fingerprint {
	return &Fingerprint{Schema: SchemaID, SchemaVersion: SchemaVersion}
}

// Validate checks structural invariants the JSON Schema also enforces, so Go
// callers fail fast. It is intentionally dependency-free.
func (f *Fingerprint) Validate() error {
	if f.Schema != SchemaID {
		return fmt.Errorf("fingerprint: unexpected schema %q (want %q)", f.Schema, SchemaID)
	}
	if f.Tokenization.MapLocation != "on-prem" {
		return errors.New("fingerprint: tokenization.mapLocation must be on-prem")
	}
	if f.Lens.Mode != "read-only" {
		return errors.New("fingerprint: lens.mode must be read-only")
	}
	for i, w := range f.Workloads {
		if w.Token == "" {
			return fmt.Errorf("workloads[%d]: empty token", i)
		}
		switch w.Kind {
		case KindVM:
			if w.VM == nil {
				return fmt.Errorf("workloads[%d] (%s): kind vm requires vm block", i, w.Token)
			}
		case KindHost:
			if w.Host == nil {
				return fmt.Errorf("workloads[%d] (%s): kind host requires host block", i, w.Token)
			}
		case KindCluster:
			if w.Cluster == nil {
				return fmt.Errorf("workloads[%d] (%s): kind cluster requires cluster block", i, w.Token)
			}
		case KindDatabase:
			if w.Database == nil {
				return fmt.Errorf("workloads[%d] (%s): kind database requires database block", i, w.Token)
			}
		default:
			return fmt.Errorf("workloads[%d] (%s): unknown kind %q", i, w.Token, w.Kind)
		}
	}
	for i, r := range f.Recommendations {
		if r.ID == "" || r.Action == "" {
			return fmt.Errorf("recommendations[%d]: id and action are required", i)
		}
		if len(r.AppliesTo) == 0 {
			return fmt.Errorf("recommendations[%d] (%s): appliesTo must reference at least one token", i, r.ID)
		}
		if r.Status != "preliminary" && r.Status != "validated" {
			return fmt.Errorf("recommendations[%d] (%s): status must be preliminary or validated", i, r.ID)
		}
	}
	return nil
}
