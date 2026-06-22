package domain

const (
	// AccountAddressMaxLength caps account addresses. Beyond this size, key
	// storage cost dominates and operators almost always mean a different
	// data layout.
	AccountAddressMaxLength = 1024

	// LedgerNameMaxLength caps ledger names. The bound also matches the
	// largest fixed-width identifier block storage backends can safely
	// reserve for the name (zero-padded); validating upstream prevents
	// silent truncation, which would otherwise cause key collisions between
	// names sharing the first N bytes.
	LedgerNameMaxLength = 64

	// AssetPrecisionMin and AssetPrecisionMax bound the decimal precision
	// suffix of an asset (e.g. "/2" in "EUR/2"). The precision is encoded
	// as a single byte in the canonical volume key, so values must fit a
	// uint8 — and "0" is reserved to mean "no precision suffix" so it must
	// not appear in the suffix itself (otherwise "USD/0" would alias "USD").
	AssetPrecisionMin = 1
	AssetPrecisionMax = 255
)
