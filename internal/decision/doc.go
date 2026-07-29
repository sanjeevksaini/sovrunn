// Package decision defines the FEATURE-0013 Decision Record and AuditEvent
// domain value types and closed vocabularies.
//
// Domain ownership (RID-01, design closure 7): canonical FEATURE-0013 domain
// types live in this package tree. internal/apiconform remains limited to
// TypeBinding registration, conformance adapters, and executable checks; it
// does not become the canonical domain owner and must not receive divergent
// copies of FEATURE-0012 base fields. The FEATURE-0013 AuditEvent
// payload/linkage extension value types are owned here.
//
// This package is stdlib-only at the domain root and may import internal/apimeta
// only. It does not import internal/server, internal/api, or runtime/provider
// packages.
package decision
