package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidateAsset checks that an asset string matches the expected format:
// [A-Z][A-Z0-9]{0,16}(_[A-Z]{1,16})?(/[1-9]\d{0,2})?
// Examples: "USD", "EUR/2", "BTC/8", "CUSTOM_TOKEN/6".
//
// Precision rules are tight on purpose: the canonical volume key encodes the
// precision as a single byte, and ParseAssetPrecision relies on validation to
// have rejected anything that would lose information. Without these rules,
// "USD", "USD/0", "USD/02", and "USD/256" all collapse onto the same volume
// entry — cross-asset fund contamination in a double-entry ledger.
func ValidateAsset(asset string) error {
	if len(asset) == 0 {
		return ErrAssetInvalid
	}

	base, precisionStr, hasPrecision := strings.Cut(asset, "/")

	if !validateAssetBase(base) {
		return ErrAssetInvalid
	}

	if hasPrecision && !validateAssetPrecision(precisionStr) {
		return ErrAssetInvalid
	}

	return nil
}

// validateAssetPrecision enforces a canonical, uint8-safe precision suffix:
//   - 1 to 3 digits (max numeric value 255 fits in 3 chars).
//   - no leading zero (rejects "02" → 2 aliasing).
//   - numeric value in [AssetPrecisionMin, AssetPrecisionMax].
func validateAssetPrecision(s string) bool {
	if len(s) == 0 || len(s) > 3 {
		return false
	}

	if s[0] == '0' {
		return false
	}

	for i := range s {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	v, _ := strconv.Atoi(s)

	return v >= AssetPrecisionMin && v <= AssetPrecisionMax
}

// validateAssetBase checks the base part: [A-Z][A-Z0-9]{0,16}(_[A-Z]{1,16})?
func validateAssetBase(base string) bool {
	if len(base) == 0 {
		return false
	}

	head, tail, hasUnderscore := strings.Cut(base, "_")

	if !isUpperAlphaStart(head, 17) {
		return false
	}

	if hasUnderscore && !isUpperAlpha(tail, 1, 16) {
		return false
	}

	return true
}

func isUpperAlphaStart(s string, maxLen int) bool {
	if len(s) == 0 || len(s) > maxLen {
		return false
	}

	if s[0] < 'A' || s[0] > 'Z' {
		return false
	}

	for i := 1; i < len(s); i++ {
		c := s[i]
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}

	return true
}

func isUpperAlpha(s string, minLen, maxLen int) bool {
	if len(s) < minLen || len(s) > maxLen {
		return false
	}

	for i := range s {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}

	return true
}

// ParseAssetPrecision splits an asset string into its base name and precision.
// "USD/4" → ("USD", 4), "EUR" → ("EUR", 0).
func ParseAssetPrecision(asset string) (string, uint8) {
	base, precStr, found := strings.Cut(asset, "/")
	if !found {
		return asset, 0
	}

	prec, _ := strconv.ParseUint(precStr, 10, 8)

	return base, uint8(prec)
}

// FormatAsset reconstructs an asset string from base and precision.
// ("USD", 4) → "USD/4", ("EUR", 0) → "EUR".
func FormatAsset(base string, precision uint8) string {
	if precision == 0 {
		return base
	}

	return fmt.Sprintf("%s/%d", base, precision)
}
