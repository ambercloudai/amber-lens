# Workload Fingerprint — v1

The **Workload Fingerprint** is the versioned, anonymized, destination-neutral
output of an Amber Lens scan. It is the **contract** between the free Lens
(which writes it) and the paid Amber Navigator (which validates it and computes
the Operational Risk Score, TCO, what-if, and Veiled Bid Lab proposals).

- Canonical schema: [`pkg/schema/fingerprint.schema.json`](../pkg/schema/fingerprint.schema.json) (JSON Schema, Draft 2020-12)
- Go types: [`internal/fingerprint`](../internal/fingerprint/fingerprint.go)
- Worked example: [`pkg/schema/examples/fingerprint.example.json`](../pkg/schema/examples/fingerprint.example.json)
- **v1 core freezes 2026-07-01.**

## Design rules (enforced by the contract)

1. **Facts + preliminary recommendations (free); no paid analysis.** The
   fingerprint contains observed facts (inventory, config, summarized
   performance, DB editions/options, OS/EOL, dependencies) **and preliminary,
   advisory, non-monetary recommendations** — the free hook. It contains **no
   validated risk score, no TCO/cost, and no reverse-bid proposals**; those are
   computed in Navigator (the paid layer). The directional `impact` on a
   recommendation (capacity, %, counts) is the free teaser; the dollar figure is
   the paid payoff.
2. **Tokens only.** Every workload is identified by a salted token (`env-role-NN`
   scheme). The token-to-real-name map **never leaves the customer network**
   (`tokenization.mapLocation: on-prem`). Only an opaque `saltId` is recorded —
   never the salt.
3. **Destination-neutral.** No target-cloud assumptions are encoded, so the same
   fingerprint serves an objective rehost / replatform / refactor / repurchase /
   retire decision.
4. **Summarized, not raw.** Performance is window-summarized
   (`avg / p50 / p95 / peak`), never raw 5-minute series — compact and
   privacy-preserving.
5. **Transparent provenance.** `lens`, `scan`, and per-DB `collection` record how
   data was read (e.g., Oracle `statspack` fallback when the Diagnostic Pack is
   unlicensed), backing the report's "what we read" promise.

## Top-level shape

| Field | Meaning |
|---|---|
| `schema` / `schemaVersion` | Contract id (`amber.fingerprint/v1`) + SemVer |
| `scanId`, `generatedAt` | Opaque scan id + timestamp |
| `lens` | Collector version + `mode: read-only` |
| `scan` | Window, sample interval, cycles, **confidence** (e.g. 0.95) |
| `tokenization` | Scheme, algorithm, `saltId`, `mapLocation: on-prem` |
| `environment` | Anonymized site/datacenter tokens |
| `summary` | Estate `counts` + capacity `totals` |
| `workloads[]` | The core inventory (see below) |
| `dependencies[]` | Tokenized relationship edges (wave planning) |
| `collectionIssues[]` | Non-fatal notes (e.g., fallbacks) |
| `recommendations[]` | Preliminary, advisory, **non-monetary** optimization/target recs (the free hook) |

## Workload model

Each workload has a `token` and a `kind` (`vm` · `host` · `cluster` ·
`database`). The sub-object matching `kind` is **required** (enforced by the
schema's `if/then` rules and by Go `Validate()`). Shared optional blocks: `os`,
`compute`, `performance`. Profile-specific additions go in `extensions` so they
never break the v1 core.

- **vm** — hypervisor, cluster token, power state, uptime (+ `os`, `compute`, `performance`)
- **host** — hypervisor + version, sockets, physical cores
- **cluster** — hypervisor/version, host count, DRS/HA/vSAN/NSX, hypervisor licensing
- **database** — engine, edition, version, size, licensable cores, HA, features,
  **options** (Oracle Partitioning / Advanced Security / Diagnostic & Tuning
  packs, etc. with `inUse`), summarized `workload`, and `collection` provenance

## Recommendations (the free hook)

Each recommendation is **preliminary, advisory, and non-monetary** and links to
the workload token(s) it concerns:

- `action` — right-size · consolidate · retire-eol · rehost · replatform · refactor · repurchase · retire · modernize
- `evidence` — the facts that drove it (e.g. `cpuPct.p95=34% over 72h`)
- `targetProfile` — the suggested optimized target (directional)
- `impact` — **ranges/counts only, never dollars** (e.g. `vcpu_reclaimable 1600-2000`)
- `status` — `preliminary` (free, from the Lens). Navigator validates and flips
  it to `validated`, then attaches TCO and bids — which are **never** stored here.

This is the funnel encoded in the data model: the Lens emits directional recs for
free; Navigator (paid) validates them and completes the proposal.

## Versioning & extensibility

- Additive, backward-compatible fields bump the **minor/patch** within v1.
- New workload *kinds* (containers, storage, cloud) arrive via collector
  profiles using `extensions`, or a future major schema — without re-locking v1.
- Core objects are `additionalProperties: false` for rigor; `extensions` is the
  sanctioned escape hatch.

## v1 freeze checklist (by 2026-07-01)

- [ ] Validate every field the Discovery Report renders maps to a schema field.
- [ ] Confirm each collector (VMware, SQL Server, Oracle, OS, deps) can populate
      its part read-only.
- [ ] Confirm no field leaks a real name or secret (tokens only).
- [ ] Confirm no derived score/TCO field has crept in.
- [ ] Confirm recommendations are non-monetary (no cost fields) and carry status=preliminary.
- [ ] Tag the schema `1.0.0` and treat changes as contract changes thereafter.
