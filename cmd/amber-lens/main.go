// Command amber-lens is the Amber Lens collector: a local-first, read-only,
// agentless discovery scanner. It samples source workloads, emits a versioned
// Workload Fingerprint, and renders a local Discovery Report.
//
// Principles (non-negotiable): read-only, agentless, local-first, two-pass
// tokenization (mapping keys never leave the network), zero telemetry,
// destination-neutral. Apache-2.0.
package main

import (
	"fmt"
	"os"
)

const version = "0.0.0-scaffold"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "amber-lens:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	// TODO: CLI subcommands — scan | report | version.
	// TODO: load collectors, run 5-min sampling over 72h, tokenize, build
	//       fingerprint, render Discovery Report.
	fmt.Printf("amber-lens %s — scaffold. No telemetry. Read-only.\n", version)
	return nil
}
