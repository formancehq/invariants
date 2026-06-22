package domain

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
		{name: "empty", input: "", wantErr: ErrLedgerNameRequired},
		{name: "contains null byte", input: "ledger\x00evil", wantErr: ErrLedgerNameContainsNullByte},
		{name: "null byte only", input: "\x00", wantErr: ErrLedgerNameContainsNullByte},
		// Names that flow through `x-next-cursor` gRPC trailers must survive
		// the HTTP/2-header value charset. Anything outside printable ASCII
		// (newlines, CR, multibyte UTF-8) would either be stripped or fail
		// the stream — so the validator rejects them up-front instead of
		// admitting a name we cannot paginate.
		{name: "contains newline", input: "ledger\nevil", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains carriage return", input: "ledger\revil", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains tab", input: "ledger\tevil", wantErr: ErrLedgerNameInvalidChar},
		{name: "contains DEL", input: "ledger\x7Fevil", wantErr: ErrLedgerNameInvalidChar},
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

