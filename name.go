package invariants

// ValidateLedgerName checks that a ledger name matches the Formance
// identifier charset [a-zA-Z0-9._:-] and stays within LedgerNameMaxLength.
//
// The charset is the same as for metadata keys — both are identifiers, both
// land in HTTP/2 metadata trailers (paginated list cursors are derived from
// the ledger name), and both benefit from being safe to render in CLI
// output, logs and URLs without escaping.
//
// The length cap also matches the largest fixed-width identifier block
// storage backends can safely reserve for the name (zero-padded); going
// beyond would risk silent truncation collisions in any key layout that
// allocates a constant slot for the ledger name.
func ValidateLedgerName(name string) error {
	if name == "" {
		return ErrLedgerNameRequired
	}

	for _, r := range name {
		if !isIdentifierChar(r) {
			return ErrLedgerNameInvalidChar
		}
	}

	if len(name) > LedgerNameMaxLength {
		return ErrLedgerNameTooLong
	}

	return nil
}
