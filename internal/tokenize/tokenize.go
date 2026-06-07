// Package tokenize implements two-pass tokenization. Pass 1 builds a local
// token<->real-name map that NEVER leaves the network. Pass 2 emits only
// salted, deterministic tokens into the fingerprint. The map is the
// customer's; Amber never receives it.
package tokenize

// Tokenizer maps real identifiers to stable, salted tokens (local only).
type Tokenizer struct{ /* TODO: salt, local map persistence */ }

// Token returns the deterministic token for a real name. The mapping is
// retained ONLY in the local map file under the customer's control.
func (t *Tokenizer) Token(realName string) string {
	// TODO: HMAC(salt, realName) -> short stable token (env-role-NN scheme).
	return ""
}
