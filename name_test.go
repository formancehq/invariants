package invariants

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateLedgerName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "valid simple name", input: "default"},
		{name: "valid with hyphens", input: "my-ledger-123"},
		{name: "valid with dots", input: "ledger.prod"},
		{name: "valid with colons", input: "tenant:prod"},
		{name: "valid with underscores", input: "ledger_test"},
		{name: "valid mixed", input: "tenant:prod.eu-west-1"},
		{name: "empty", input: "", wantErr: ErrLedgerNameRequired},
		// The charset is identifier-style — no whitespace, no free
		// punctuation, no UTF-8 multibyte, no null bytes. Same rule as
		// metadata keys, since both can land in HTTP/2 metadata trailers
		// and benefit from being safe to render in CLI / logs / URLs.
		{name: "contains null byte", input: "ledger\x00evil", wantErr: ErrLedgerNameInvalidChar},
		{name: "null byte only", input: "\x00", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains newline", input: "ledger\nevil", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains carriage return", input: "ledger\revil", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains tab", input: "ledger\tevil", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains space", input: "my ledger", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains slash", input: "tenant/prod", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains question mark", input: "ledger?", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains non-ASCII utf8", input: "ledgér", wantErr: ErrLedgerNameInvalidChar},
		{name: "too long", input: strings.Repeat("a", LedgerNameMaxLength+1), wantErr: ErrLedgerNameTooLong},
		{name: "max length", input: strings.Repeat("a", LedgerNameMaxLength)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateLedgerName(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
