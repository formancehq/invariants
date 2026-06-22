package invariants

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateLedgerAccountAddress(t *testing.T) {
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
		{name: "empty", input: "", wantErr: ErrLedgerAccountAddressEmpty},
		{name: "contains null byte", input: "account\x00evil", wantErr: ErrLedgerAccountAddressInvalidChar},
		{name: "contains space", input: "my account", wantErr: ErrLedgerAccountAddressInvalidChar},
		{name: "valid hyphen", input: "my-account"},
		{name: "valid underscore", input: "my_account"},
		{name: "contains dot", input: "my.account", wantErr: ErrLedgerAccountAddressInvalidChar},
		{name: "contains slash", input: "a/b", wantErr: ErrLedgerAccountAddressInvalidChar},
		// Segments separated by `:` must be non-empty. No leading, trailing,
		// or consecutive colons.
		{name: "leading colon", input: ":world", wantErr: ErrLedgerAccountAddressEmptySegment},
		{name: "trailing colon", input: "users:", wantErr: ErrLedgerAccountAddressEmptySegment},
		{name: "consecutive colons", input: "users::alice", wantErr: ErrLedgerAccountAddressEmptySegment},
		{name: "only colon", input: ":", wantErr: ErrLedgerAccountAddressEmptySegment},
		{name: "only colons", input: ":::", wantErr: ErrLedgerAccountAddressEmptySegment},
		{name: "too long", input: strings.Repeat("a", LedgerAccountAddressMaxLength+1), wantErr: ErrLedgerAccountAddressTooLong},
		{name: "max length", input: strings.Repeat("a", LedgerAccountAddressMaxLength)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateLedgerAccountAddress(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
