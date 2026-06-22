# formancehq/domain

Formance-wide business domain contracts. This package centralises the invariants — formats,
validations, bounds — that the Formance platform shares across services (ledger, payments,
wallets, reconciliation, etc.) and SDKs.

The package keeps a minimal API surface (primitive Go types in, `error` out) so it can be
imported from any Go tool — server, CLI, operator, SDK — and so the same rules can later be
exported as a declarative spec for non-Go clients.

## Contract reference

### Account address

Identifies an account in any Formance service.

- **Format**: `[a-zA-Z0-9_:-]+`
- **Max length**: 1024 bytes (`AccountAddressMaxLength`)
- **Separator**: `:` (colon) used to nest segments (e.g. `users:alice:checking`)
- **Allowed characters**: letters (A-Z, a-z), digits (0-9), colon (`:`), underscore (`_`), hyphen (`-`)
- **Forbidden**: dots (`.`), slashes (`/`), whitespace, null bytes, any other punctuation
- **Examples**: `world`, `users:alice:checking`, `platform:fees`, `treasury-2026`

### Asset

Identifies a unit of value (currency, token) carried by a posting.

- **Format**: `[A-Z][A-Z0-9]{0,16}(/[1-9][0-9]{0,2})?`
- **Base**: uppercase letter followed by 0–16 uppercase letters or digits — max 17 characters total
- **Precision suffix**: optional `/N` where `N ∈ [1, 255]` (`AssetPrecisionMin`, `AssetPrecisionMax`)
- **No underscores**: codes like `CUSTOM_TOKEN` or `USD_EUR` are rejected; the base must be a single uppercase code
- **Canonical form**: every accepted asset round-trips through `ParseAssetPrecision` → `FormatAsset` unchanged. Aliases like `USD/0`, `USD/02`, `USD/256` are rejected so the same volume never has two key encodings.
- **Examples**: `USD`, `EUR/2`, `BTC/8`, `USDC2/6`, `ABCDEFGHIJKLMNOPQ`

### Ledger name

Identifies a ledger / tenant namespace.

- **Format**: printable ASCII only (bytes `0x20`–`0x7E`)
- **Max length**: 64 bytes (`LedgerNameMaxLength`)
- **Forbidden**: null bytes, CR, LF, tab, DEL (`0x7F`), any non-ASCII (UTF-8 multibyte)
- **Rationale**: ledger names land in `x-next-cursor` gRPC trailers used by paginated list endpoints. They must survive HTTP/2 metadata round-trips. The 64-byte bound matches the fixed-width block reserved for the ledger name in canonical key prefixes (zero-padded); going beyond would cause silent truncation collisions.
- **Examples**: `default`, `my-ledger-123`, `ledger.prod`

### Signing key ID

Identifies an Ed25519 signing key used for request signing.

- **Format**: printable ASCII only (bytes `0x20`–`0x7E`)
- **Max length**: 256 bytes (`SigningKeyIDMaxLength`)
- **Forbidden**: same as ledger names — control bytes, null bytes, non-ASCII
- **Rationale**: same trailer-envelope concern (paginated `signing keys list` cursor)
- **Examples**: `admin-key-1`, `kms.prod.2026`, `team/treasury/v2`

### Metadata key

Identifies a metadata entry attached to an account, transaction, or ledger.

- **Format**: any UTF-8 string except null bytes
- **Required**: must be non-empty
- **Forbidden**: null byte (`0x00`) — would corrupt null-terminated key encodings used by the metadata read index
- **Note**: there is no length cap at this layer; storage backends may enforce one (e.g. Pebble key size).
- **Examples**: `category`, `user.role`, `audit:reviewer`, `compliance-tier`

### Metadata string value

A string-typed metadata payload (other typed variants — int, uint, bool, null — are accepted unchanged at the wire layer).

- **Format**: any UTF-8 string except null bytes
- **Allowed**: empty string
- **Forbidden**: null byte (`0x00`)
- **Examples**: `admin`, `tier-3`, `2026-06-22T10:00:00Z`

## Stability

`v0.x` is treated as evolving but already covers the invariants used by `ledger-v3-poc`.
Once a stable `v1` is tagged, any change to a sentinel identity or to a validation rule
becomes a breaking change subject to SemVer.

## Usage

```go
import domain "github.com/formancehq/domain"

if err := domain.ValidateLedgerName(name); err != nil {
    // err matches one of domain.ErrLedgerName* sentinels via errors.Is
}
```

## Contribution

Owners: `@formancehq/backend`. Open a PR; CI runs `go test` + `golangci-lint`.
