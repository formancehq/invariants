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
	ErrLedgerNameRequired         = errors.New("ledger name is required")
	ErrLedgerNameContainsNullByte = errors.New("ledger name must not contain null bytes")
	ErrLedgerNameInvalidChar      = errors.New("ledger name must contain only printable ASCII (0x20–0x7E)")
	ErrLedgerNameTooLong          = fmt.Errorf("ledger name exceeds maximum length of %d bytes", LedgerNameMaxLength)
)

// Metadata sentinels.
var (
	ErrMetadataKeyEmpty              = errors.New("metadata key must not be empty")
	ErrMetadataKeyContainsNullByte   = errors.New("metadata key must not contain null bytes")
	ErrMetadataValueContainsNullByte = errors.New("metadata value must not contain null bytes")
)
