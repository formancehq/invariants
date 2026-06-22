package invariants

import "strings"

// isLedgerAccountAddressChar returns true if the rune is allowed in a ledger
// account address segment. Segments are [a-zA-Z0-9_-]+, joined by colons.
func isLedgerAccountAddressChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ':' || r == '_' || r == '-'
}

// ValidateLedgerAccountAddress checks that a ledger account address matches
// `[a-zA-Z0-9_-]+(:[a-zA-Z0-9_-]+)*` and stays within
// LedgerAccountAddressMaxLength.
//
// The colon is the conventional hierarchy separator: segments must be
// non-empty (no leading colon, no trailing colon, no consecutive colons).
// "users::alice", ":world", "users:" are all rejected — they would
// collapse onto distinct storage keys without a meaningful hierarchy and
// confuse downstream tooling.
func ValidateLedgerAccountAddress(address string) error {
	if address == "" {
		return ErrLedgerAccountAddressEmpty
	}

	if len(address) > LedgerAccountAddressMaxLength {
		return ErrLedgerAccountAddressTooLong
	}

	for _, r := range address {
		if !isLedgerAccountAddressChar(r) {
			return ErrLedgerAccountAddressInvalidChar
		}
	}

	for _, segment := range strings.Split(address, ":") {
		if segment == "" {
			return ErrLedgerAccountAddressEmptySegment
		}
	}

	return nil
}
