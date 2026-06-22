# formancehq/domain

Formance-wide business domain contracts. This package centralises the invariants — formats,
validations, bounds — that the Formance platform shares across services (ledger, payments,
wallets, reconciliation, etc.) and SDKs.

The package keeps a minimal API surface (primitive Go types in, `error` out) so it can be
imported from any Go tool — server, CLI, operator, SDK — and so the same rules can later be
exported as a declarative spec for non-Go clients.

## Scope

Current contents:

- Account address format (charset, max length).
- Asset format (regex, precision range, canonical parse/format round-trip).
- Identifier names (ledger, numscript, signing-key id): printable-ASCII charset bounded for
  HTTP/2 metadata trailers used by paginated list cursors.
- Metadata key/value validation (null-byte rejection).

Future additions:

- Page-size and request-limit constants.
- Cross-service business enums.
- Declarative specs (JSON Schema / proto descriptors) for non-Go SDKs.

## Stability

`v0.x` is treated as evolving but already covers the invariants used by `ledger-v3-poc`.
Once a stable `v1` is tagged, any change to a sentinel identity or to a validation rule
becomes a breaking change subject to SemVer.

## Usage

```go
import "github.com/formancehq/domain"

if err := domain.ValidateLedgerName(name); err != nil {
    // err matches one of domain.ErrLedgerName* sentinels via errors.Is
}
```

## Contribution

Owners: `@formancehq/backend`. Open a PR; CI runs `go test` + `golangci-lint`.
