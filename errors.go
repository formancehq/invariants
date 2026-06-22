package domain

import (
	"errors"
	"fmt"
)

// Account address sentinels.
var (
	ErrAccountAddressEmpty       = errors.New("account address must not be empty")
	ErrAccountAddressInvalidChar = errors.New("account address must contain only letters, digits, colons, underscores, and hyphens")
	ErrAccountAddressTooLong     = fmt.Errorf("account address exceeds maximum length of %d bytes", AccountAddressMaxLength)
)

// Asset sentinels.
var (
	ErrAssetInvalid = errors.New("asset must match [A-Z][A-Z0-9]{0,16}(/[1-9][0-9]{0,2})? with precision in [1, 255]")
)

// Ledger / numscript / signing-key name sentinels.
var (
	ErrLedgerNameRequired    = errors.New("ledger name is required")
	ErrLedgerNameInvalidChar = errors.New("ledger name must match [a-zA-Z0-9._:-]")
	ErrLedgerNameTooLong     = fmt.Errorf("ledger name exceeds maximum length of %d bytes", LedgerNameMaxLength)
)

// Metadata sentinels.
var (
	ErrMetadataKeyEmpty              = errors.New("metadata key must not be empty")
	ErrMetadataKeyInvalidChar        = errors.New("metadata key must match [a-zA-Z0-9._:-]")
	ErrMetadataValueContainsNullByte = errors.New("metadata value must not contain null bytes")
)
