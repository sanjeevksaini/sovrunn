// Package cloudmodel implements the FEATURE-0015 canonical cloud model
// foundation store and coordination primitives (design §3.1, DD-01, DD-05).
//
// Resource types live in cloudmodel/model; pure validators live in
// cloudmodel/validate; the assigned ISO dataset lives in cloudmodel/isocodes.
// This package owns the in-memory registries, name/pair uniqueness,
// resourceVersion, and publication lock/staging primitives used by later
// audit and handler tasks.
package cloudmodel
