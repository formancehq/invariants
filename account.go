package invariants

// isLedgerAccountAddressByte returns true if the byte is allowed in a
// ledger account address. The valid alphabet is [a-zA-Z0-9_-] plus the
// segment separator ':'. Byte-typed (not rune) because the alphabet is
// pure ASCII: any multi-byte UTF-8 sequence fails on its first byte.
func isLedgerAccountAddressByte(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == ':' || c == '_' || c == '-'
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
//
// Implemented as a single byte-wise pass to avoid the strings.Split
// allocation that the earlier two-pass implementation paid on every
// call — Ledger's FSM apply path validates every posting's source and
// destination, so the hot path calls this function O(postings) times
// per proposal.
func ValidateLedgerAccountAddress(address string) error {
	if address == "" {
		return ErrLedgerAccountAddressEmpty
	}

	if len(address) > LedgerAccountAddressMaxLength {
		return ErrLedgerAccountAddressTooLong
	}

	// prevColon tracks whether the previous byte was a colon so we can
	// reject empty segments in a single pass:
	//   - initialised to true so a leading ':' fires the check
	//   - stays true across consecutive colons so "::" fires the check
	//   - final value of true means the address ended on ':', so a
	//     trailing colon fires the check after the loop
	prevColon := true
	for i := 0; i < len(address); i++ {
		c := address[i]
		if !isLedgerAccountAddressByte(c) {
			return ErrLedgerAccountAddressInvalidChar
		}

		if c == ':' {
			if prevColon {
				return ErrLedgerAccountAddressEmptySegment
			}
			prevColon = true
		} else {
			prevColon = false
		}
	}

	if prevColon {
		return ErrLedgerAccountAddressEmptySegment
	}

	return nil
}
