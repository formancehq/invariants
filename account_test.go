package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAccountAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "valid simple", input: "world"},
		{name: "valid with colon segments", input: "users:alice:checking"},
		{name: "valid uppercase", input: "USD"},
		{name: "valid mixed", input: "platform:fees"},
		{name: "valid digits", input: "user123"},
		{name: "empty", input: "", wantErr: ErrAccountAddressEmpty},
		{name: "contains null byte", input: "account\x00evil", wantErr: ErrAccountAddressInvalidChar},
		{name: "contains space", input: "my account", wantErr: ErrAccountAddressInvalidChar},
		{name: "valid hyphen", input: "my-account"},
		{name: "valid underscore", input: "my_account"},
		{name: "contains dot", input: "my.account", wantErr: ErrAccountAddressInvalidChar},
		{name: "contains slash", input: "a/b", wantErr: ErrAccountAddressInvalidChar},
		{name: "too long", input: strings.Repeat("a", AccountAddressMaxLength+1), wantErr: ErrAccountAddressTooLong},
		{name: "max length", input: strings.Repeat("a", AccountAddressMaxLength)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateAccountAddress(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
