# formancehq/invariants

Formance product invariants. This package centralises the formats, validations, and bounds
that are shared across Formance services, SDKs, operators, and tools.

The package keeps a minimal API surface (primitive Go types in, `error` out) so it can be
imported from any Go tool - server, CLI, operator, SDK - and so the same rules can later be
exported as a declarative spec or test-vector suite for non-Go clients.

## Scope

This module is a source of truth for **shared product contracts**, not a general helper
library. A rule belongs here only when consumers outside the owning service need the exact
same invariant to validate input or generate compatible data.

The API names should make the owning product concept explicit. For example, a ledger account
address is exposed as `ValidateLedgerAccountAddress`, not `ValidateAccountAddress`, because
`account` means different things in Ledger, Payments, and Banking Bridge.

Rules that are specific to one service should stay in that service. General Go helpers should
stay in `go-libs`. Organization, stack, and membership identifiers are out of scope unless the
architecture explicitly decides to promote them to shared product invariants.

## Contract reference

### Ledger account address

Identifies an account in a Formance ledger.

- **Format**: `[a-zA-Z0-9_-]+(:[a-zA-Z0-9_-]+)*` — segments joined by colons
- **Max length**: 256 bytes (`LedgerAccountAddressMaxLength`)
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

`v0.x` is treated as evolving, but every rule change should already be reviewed as a product
contract change. Once a stable `v1` is tagged, any change to a sentinel identity, exported
constant, or validation rule is subject to SemVer.

### Update policy

Validation changes are compatibility-sensitive because services, SDKs, CRDs, CLIs, and
persisted data may all depend on the same accepted value set.

- Do not silently widen or narrow the accepted values of an existing validator.
- Treat stricter validation as a breaking change unless all existing persisted and client-side
  values are known to be compatible.
- Treat looser validation as a contract change that must document downstream effects before
  release.
- Document the impact, migration path, and affected services for every rule change.
- Require review from the owners of the affected product concept, not only the maintainers of
  this Go module.
- Add or update test vectors for every validator change so services and future non-Go SDKs can
  stay aligned.
- Prefer adding a new explicitly named validator over broadening a concept name when two
  services use the same word for different product concepts.

## Usage

```go
import "github.com/formancehq/invariants"

if err := invariants.ValidateLedgerName(name); err != nil {
    // err matches one of invariants.ErrLedgerName* sentinels via errors.Is
}

if err := invariants.ValidateLedgerAccountAddress(address); err != nil {
    // err matches one of invariants.ErrLedgerAccountAddress* sentinels via errors.Is
}
```

## Contribution

Owners: `@formancehq/backend`. Open a PR; CI runs `go test` + `golangci-lint`.
