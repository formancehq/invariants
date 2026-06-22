package invariants

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAsset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "simple", input: "USD"},
		{name: "with precision", input: "EUR/2"},
		{name: "long precision", input: "BTC/8"},
		{name: "max base length", input: "ABCDEFGHIJKLMNOPQ"},
		{name: "single char", input: "A"},
		{name: "alphanumeric base", input: "USD2"},
		{name: "empty", input: "", wantErr: ErrAssetInvalid},
		{name: "lowercase", input: "usd", wantErr: ErrAssetInvalid},
		{name: "starts with digit", input: "1USD", wantErr: ErrAssetInvalid},
		{name: "contains hyphen", input: "US-D", wantErr: ErrAssetInvalid},
		{name: "contains space", input: "US D", wantErr: ErrAssetInvalid},
		{name: "base too long", input: "ABCDEFGHIJKLMNOPQR", wantErr: ErrAssetInvalid},
		{name: "precision too long", input: "USD/1234567", wantErr: ErrAssetInvalid},
		// Underscores were accepted in earlier revisions; they are now rejected
		// outright (no more CUSTOM_TOKEN / USD_EUR forms).
		{name: "underscore inside base", input: "CUSTOM_TOKEN", wantErr: ErrAssetInvalid},
		{name: "underscore base with precision", input: "CUSTOM_TOKEN/6", wantErr: ErrAssetInvalid},
		{name: "underscore suffix uppercase", input: "USD_EUR", wantErr: ErrAssetInvalid},
		{name: "underscore suffix lowercase", input: "USD_eur", wantErr: ErrAssetInvalid},
		{name: "trailing underscore", input: "USD_", wantErr: ErrAssetInvalid},
		{name: "leading underscore", input: "_USD", wantErr: ErrAssetInvalid},
		{name: "double slash", input: "USD//2", wantErr: ErrAssetInvalid},
		{name: "trailing slash", input: "USD/", wantErr: ErrAssetInvalid},
		// Precision must fit in uint8 (the volume-key byte) and use a single
		// canonical form per (base, precision) pair.
		{name: "precision overflows uint8", input: "USD/256", wantErr: ErrAssetInvalid},
		{name: "precision way over uint8", input: "USD/999999", wantErr: ErrAssetInvalid},
		{name: "precision boundary 255", input: "USD/255"},
		{name: "precision min 1", input: "USD/1"},
		{name: "precision zero aliases bare base", input: "USD/0", wantErr: ErrAssetInvalid},
		{name: "precision leading zero", input: "USD/02", wantErr: ErrAssetInvalid},
		{name: "precision multi leading zero", input: "USD/007", wantErr: ErrAssetInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateAsset(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidateAsset_CanonicalRoundTrip pins the contract used by the
// volume-key encoding: every input ValidateAsset accepts must survive the
// ParseAssetPrecision → FormatAsset round trip unchanged. If two valid inputs
// collapsed onto the same (base, precision) pair, the canonical form returned
// by FormatAsset would not match one of them and consensus-deterministic
// asset aliasing would already exist in production.
func TestValidateAsset_CanonicalRoundTrip(t *testing.T) {
	t.Parallel()

	valid := []string{
		"USD",
		"EUR/2",
		"BTC/8",
		"USD/1",
		"USD/255",
		"A",
		"USD2",
		"ABCDEFGHIJKLMNOPQ",
	}

	for _, asset := range valid {
		t.Run(asset, func(t *testing.T) {
			t.Parallel()

			require.NoError(t, ValidateAsset(asset))

			base, precision := ParseAssetPrecision(asset)
			require.Equal(t, asset, FormatAsset(base, precision),
				"every accepted asset must round-trip through Parse/Format unchanged")
		})
	}
}
