package invariants

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateMetadataKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "valid simple", input: "category"},
		{name: "valid with dot", input: "user.role"},
		{name: "valid with colon", input: "audit:reviewer"},
		{name: "valid with hyphen", input: "compliance-tier"},
		{name: "valid with underscore", input: "internal_flag"},
		{name: "valid mixed", input: "system:user.role-v2"},
		{name: "valid digits", input: "tag123"},
		{name: "valid uppercase", input: "TAG"},
		{name: "empty", input: "", wantErr: ErrMetadataKeyEmpty},
		// The charset is identifier-style: no whitespace, no free
		// punctuation, no UTF-8 multibyte, no null bytes.
		{name: "contains null byte", input: "key\x00value", wantErr: ErrMetadataKeyInvalidChar},
		{name: "contains space", input: "user role", wantErr: ErrMetadataKeyInvalidChar},
		{name: "contains slash", input: "user/role", wantErr: ErrMetadataKeyInvalidChar},
		{name: "contains question mark", input: "user?role", wantErr: ErrMetadataKeyInvalidChar},
		{name: "contains comma", input: "a,b", wantErr: ErrMetadataKeyInvalidChar},
		{name: "contains equals", input: "a=b", wantErr: ErrMetadataKeyInvalidChar},
		{name: "contains percent", input: "100%", wantErr: ErrMetadataKeyInvalidChar},
		{name: "contains non-ASCII utf8", input: "étape", wantErr: ErrMetadataKeyInvalidChar},
		{name: "contains emoji", input: "key🚀", wantErr: ErrMetadataKeyInvalidChar},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateMetadataKey(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateMetadataString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "valid string", input: "admin"},
		{name: "empty", input: ""},
		{name: "non-ASCII accepted", input: "Léa"},
		{name: "free punctuation accepted", input: "Validé pour clôture (étape #2)"},
		{name: "contains null byte", input: "admin\x00evil", wantErr: ErrMetadataValueContainsNullByte},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateMetadataString(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
