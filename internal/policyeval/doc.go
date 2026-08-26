// Package policyeval is the FEATURE-0017 engine-neutral policy evaluation seam.
//
// It provides an in-process evaluation boundary for conformance and adapter
// abstraction. This package is not a policy engine, not a policy system, and
// does not select or execute production policy engines.
//
// The injected UTC time source is a dependency-injection pattern for
// deterministic chronology (evaluatedAt), not a resource, store, or platform
// service. FEATURE-0017 does not import FEATURE-0016 clock types.
package policyeval
