package domain

// isAccountAddressChar returns true if the rune is allowed in an account address.
// Segments are [a-zA-Z0-9_-]+, separated by colons.
func isAccountAddressChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ':' || r == '_' || r == '-'
}

// ValidateAccountAddress checks that an account address contains only allowed characters
// (letters, digits, colons, underscores, hyphens) and is within length limits.
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

	return nil
}
