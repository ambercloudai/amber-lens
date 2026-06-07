# Architecture

Amber Lens is a single Go binary deployed in the customer's environment.

```
collectors → tokenize (two-pass) → fingerprint (v1) → report (HTML)
                 ▲ keys stay local        ▲ tokens only
          5-min sampling over 72h (sampling pkg)
```

- **Local-first, cloud-last.** Collection and obfuscation happen on-prem; only
  salted, deterministic tokens leave the network.
- **Read-only & agentless.** vCenter API, SQL Server DMVs, Oracle v$/AWR
  (Statspack fallback), WMI/SSH. No agents, no writes.
- **Destination-neutral.** The fingerprint serves an objective
  rehost/replatform/refactor/repurchase/retire decision regardless of target.
- **Free vs paid.** The Lens + Discovery Report are free/open-source. Scored
  analysis (Operational Risk Score, TCO, what-if, Veiled Bid Lab) is the paid
  Navigator layer and lives elsewhere.
