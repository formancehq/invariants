// Package invariants centralises the Formance product invariants that are shared
// across services and SDKs: identifier formats, validation rules, and bounds.
//
// The package exposes a primitive API (string in, error out) and exports
// sentinel errors that can be matched via errors.Is. Constants such as
// LedgerNameMaxLength are exported so clients can validate ahead of an RPC
// call.
package invariants
