package domain

import "strings"

// isPrintableASCII reports whether every byte of s is in the printable ASCII
// range 0x20–0x7E (space through tilde). The bound matches the safe-value
// subset accepted by gRPC metadata headers (HTTP/2 fields strip CR/LF, reject
// control bytes, and have no defined encoding for high bytes), so any
// identifier we plan to round-trip through `x-next-cursor` trailers must
// satisfy this predicate.
func isPrintableASCII(s string) bool {
	for i := range len(s) {
		b := s[i]
		if b < 0x20 || b > 0x7E {
			return false
		}
	}

	return true
}

// ValidateLedgerName checks that a ledger name is safe for use in key encoding
// AND for transport through gRPC metadata trailers (paginated list cursors
// are derived from the ledger name). Null bytes would corrupt null-terminated
// key layouts; control or high bytes would break the `x-next-cursor` resume
// token. Length is capped to keep keys reasonable.
func ValidateLedgerName(name string) error {
	if name == "" {
		return ErrLedgerNameRequired
	}

	if strings.ContainsRune(name, 0) {
		return ErrLedgerNameContainsNullByte
	}

	if !isPrintableASCII(name) {
		return ErrLedgerNameInvalidChar
	}

	if len(name) > LedgerNameMaxLength {
		return ErrLedgerNameTooLong
	}

	return nil
}

// ValidateSigningKeyID mirrors ValidateLedgerName: the key ID lands in the
// `x-next-cursor` trailer of `signing keys list` pagination, so it must be
// safe for HTTP/2 header values (printable ASCII, bounded length). Parent
// key IDs go through the same rule so revoke/cascade traversals cannot
// smuggle in an unsafe identifier either.
func ValidateSigningKeyID(id string) error {
	if id == "" {
		return ErrSigningKeyIDRequired
	}

	if !isPrintableASCII(id) {
		return ErrSigningKeyIDInvalidChar
	}

	if len(id) > SigningKeyIDMaxLength {
		return ErrSigningKeyIDTooLong
	}

	return nil
}
