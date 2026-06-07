// Package report renders the local Discovery Report (HTML) — the free
// Foundation artifact: inventory, configuration, OS/EOL, dependency view, and
// optimization/right-sizing recommendations. Scored analysis (Risk Score, TCO)
// is the paid Navigator layer and is NOT produced here.
package report

import "github.com/ambercloudai/amber-lens/internal/fingerprint"

// Render writes a self-contained Discovery Report HTML for the given fingerprint.
func Render(fp *fingerprint.Fingerprint, outPath string) error {
	// TODO: render the dark/amber Discovery Report (see design mockups).
	return nil
}
