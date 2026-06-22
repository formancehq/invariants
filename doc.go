// Package domain centralises the Formance-wide business invariants: identifier formats,
// validation rules, and bounds shared across services and SDKs.
//
// The package exposes a primitive API (string in, error out) and exports sentinel errors
// that can be matched via errors.Is. Constants such as LedgerNameMaxLength are exported so
// clients can validate ahead of an RPC call.
package domain
