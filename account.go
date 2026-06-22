package domain

import "strings"

// isAccountAddressChar returns true if the rune is allowed in an account
// address segment. Segments are [a-zA-Z0-9_-]+, joined by colons.
func isAccountAddressChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ':' || r == '_' || r == '-'
}

// ValidateAccountAddress checks that an account address matches
// `[a-zA-Z0-9_-]+(:[a-zA-Z0-9_-]+)*` and stays within
// AccountAddressMaxLength.
//
// The colon is the conventional hierarchy separator: segments must be
// non-empty (no leading colon, no trailing colon, no consecutive colons).
// "users::alice", ":world", "users:" are all rejected — they would
// collapse onto distinct storage keys without a meaningful hierarchy and
// confuse downstream tooling.
func ValidateAccountAddress(address string) error {
	if address == "" {
		return ErrAccountAddressEmpty
	}

	if len(address) > AccountAddressMaxLength {
		return ErrAccountAddressTooLong
	}

	for _, r := range address {
		if !isAccountAddressChar(r) {
			return ErrAccountAddressInvalidChar
		}
	}

	for _, segment := range strings.Split(address, ":") {
		if segment == "" {
			return ErrAccountAddressEmptySegment
		}
	}

	return nil
}
