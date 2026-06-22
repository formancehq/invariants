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

func TestValidateNumscriptName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "valid simple", input: "transfer"},
		{name: "valid with hyphens", input: "user-transfer-v2"},
		{name: "valid with dots", input: "transfers.v1"},
		{name: "empty", input: "", wantErr: ErrNumscriptNameRequired},
		// Same charset rationale as ledger names: numscript names become
		// resume cursors on the `numscripts list` stream, so they must be
		// safe for gRPC metadata trailers.
		{name: "contains newline", input: "trans\nfer", wantErr: ErrNumscriptNameInvalidChar},
		{name: "contains null byte", input: "trans\x00fer", wantErr: ErrNumscriptNameInvalidChar},
		{name: "contains non-ASCII utf8", input: "transférer", wantErr: ErrNumscriptNameInvalidChar},
		{name: "too long", input: strings.Repeat("a", NumscriptNameMaxLength+1), wantErr: ErrNumscriptNameTooLong},
		{name: "max length", input: strings.Repeat("a", NumscriptNameMaxLength)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateNumscriptName(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateSigningKeyID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "valid simple", input: "admin-key-1"},
		{name: "valid with dots", input: "kms.prod.2026"},
		{name: "valid with slashes", input: "team/treasury/v2"},
		{name: "empty", input: "", wantErr: ErrSigningKeyIDRequired},
		// Key IDs end up in the `x-next-cursor` trailer for the
		// `signing keys list` stream — same HTTP/2-header charset constraint
		// as ledger and numscript names.
		{name: "contains newline", input: "key\nid", wantErr: ErrSigningKeyIDInvalidChar},
		{name: "contains null byte", input: "key\x00id", wantErr: ErrSigningKeyIDInvalidChar},
		{name: "contains non-ASCII utf8", input: "clé", wantErr: ErrSigningKeyIDInvalidChar},
		{name: "too long", input: strings.Repeat("a", SigningKeyIDMaxLength+1), wantErr: ErrSigningKeyIDTooLong},
		{name: "max length", input: strings.Repeat("a", SigningKeyIDMaxLength)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateSigningKeyID(tt.input)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
