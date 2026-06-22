package domain

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
		{name: "valid key", input: "category"},
		{name: "valid with dots", input: "user.role"},
		{name: "empty", input: "", wantErr: ErrMetadataKeyEmpty},
		{name: "contains null byte", input: "key\x00value", wantErr: ErrMetadataKeyContainsNullByte},
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
