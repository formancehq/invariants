// Package domain centralises the Formance transaction-domain invariants: identifier
// formats, validation rules, and bounds shared across services and SDKs around the
// platform's transaction model.
//
// The package exposes a primitive API (string in, error out) and exports sentinel errors
// that can be matched via errors.Is. Constants such as LedgerNameMaxLength are exported so
// clients can validate ahead of an RPC call.
package domain
