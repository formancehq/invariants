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

- **Format**: `[a-zA-Z0-9_-]+(:[a-zA-Z0-9_-]+)*` — segments joined by colons
- **Max length**: 256 bytes (`AccountAddressMaxLength`)
- **Separator**: `:` (colon) used to nest segments (e.g. `users:alice:checking`)
- **Allowed characters per segment**: letters (A-Z, a-z), digits (0-9), underscore (`_`), hyphen (`-`)
- **Segments**: must be non-empty — no leading, trailing, or consecutive colons (`:users`, `users:`, `users::alice` all rejected)
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

- **Format**: `[a-zA-Z0-9._:-]+` (same identifier charset as metadata keys)
- **Max length**: 64 bytes (`LedgerNameMaxLength`)
- **Allowed characters**: letters (A-Z, a-z), digits (0-9), dot (`.`), colon (`:`), underscore (`_`), hyphen (`-`)
- **Forbidden**: whitespace, free punctuation, UTF-8 multibyte, null bytes
- **Rationale**: ledger names land in `x-next-cursor` gRPC trailers used by paginated list endpoints and appear in logs / CLI output / URLs, so they need the same identifier-safe charset as metadata keys. The 64-byte bound also matches the largest fixed-width identifier block storage backends can safely reserve for the name (zero-padded); going beyond would risk silent truncation collisions in any key layout that allocates a constant slot for the ledger name.
- **Examples**: `default`, `my-ledger-123`, `ledger.prod`, `tenant:prod.eu-west-1`

### Metadata key

Identifies a metadata entry attached to an account, transaction, or ledger.

- **Format**: `[a-zA-Z0-9._:-]+`
- **Required**: must be non-empty
- **Allowed characters**: letters (A-Z, a-z), digits (0-9), dot (`.`), colon (`:`), underscore (`_`), hyphen (`-`)
- **Forbidden**: whitespace, free punctuation (`/`, `?`, `=`, `,`, `%`, `<`, `>`, …), UTF-8 multibyte, null bytes
- **Rationale**: keys are identifiers — they appear in logs, CLI output, URLs and (potentially) HTTP/2 metadata trailers. The narrow charset matches the spirit of Kubernetes labels, OpenTelemetry attribute keys, and gRPC metadata header names without committing to any single platform's exact rule.
- **Note**: there is no length cap at this layer; storage backends may enforce one.
- **Examples**: `category`, `user.role`, `audit:reviewer`, `compliance-tier`, `system:user.role-v2`

### Metadata string value

A string-typed metadata payload (other typed variants — int, uint, bool, null — are accepted unchanged at the wire layer).

- **Format**: any UTF-8 string except null bytes
- **Allowed**: empty string
- **Forbidden**: null byte (`0x00`)
- **Examples**: `admin`, `tier-3`, `2026-06-22T10:00:00Z`

## Stability

`v0.x` is treated as evolving. Once a stable `v1` is tagged, any change to a sentinel
identity or to a validation rule becomes a breaking change subject to SemVer.

## Usage

```go
import domain "github.com/formancehq/domain"

if err := domain.ValidateLedgerName(name); err != nil {
    // err matches one of domain.ErrLedgerName* sentinels via errors.Is
}
```

## Contribution

Owners: `@formancehq/backend`. Open a PR; CI runs `go test` + `golangci-lint`.
