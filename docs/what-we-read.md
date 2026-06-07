# What the Lens Reads (Transparency)

Amber Lens is open-source so your security team can verify exactly what it does.
Every query is **read-only**. This document enumerates, per collector, the exact
metadata read and the access method — and ships in every scan as
`scan-manifest.json`.

| Source | Method (read-only) | Reads |
|---|---|---|
| VMware vSphere | vCenter REST/SOAP API | VM inventory + 72h performance, cluster/host config, licensing |
| SQL Server | DMV / catalog (T-SQL) | editions, sizes, cores, workload intensity, features |
| Oracle | v$/dictionary; AWR/ASH or Statspack | edition, options/packs usage, workload |
| OS / Host | WMI / PowerShell / SSH | OS version, EOL, right-sizing metrics |
| Dependencies | agentless flow sampling | host/process affinity for wave planning |

Note: Oracle AWR/ASH requires the Diagnostic Pack; Lens detects entitlement and
**falls back to Statspack/v$** so the scan never creates a licensing exposure.
