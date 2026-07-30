# Feature Traceability Matrix

This matrix links features to tenets, decisions, RFCs, tests, and gate status.

| Feature | Phase | Status | Tenets | Decisions | RFCs | Tests/Gates | Notes |
|---|---|---|---|---|---|---|---|
| FEATURE-0011 | Phase 2 | Implemented | TENET-017 | DEC-0026, DEC-0036 | RFC-0021 | Feature gate passed | Reuse assessment standard merged; ADH-2026-011 |
| FEATURE-0012 | Phase 2 | Implemented and Merged | TENET-015 | DEC-0026, DEC-0027, DEC-0036 | RFC-0022 | Checkpoint 18 passed; final human approval 2026-07-24; merged through PR #14 as commit `a1b74fb` into `phase2-reuse-first-paas-fabric-foundation` | ADH-2026-012, ADH-2026-013; canonical API/resource architecture approved |
| FEATURE-0013 | Phase 2 | Implemented and Merged | TENET-012, TENET-013 | DEC-0026, DEC-0036 | RFC-0023 | Final feature gate passed; PR #15 merged 2026-07-29 | ADH-2026-017 is the approved single replacement handoff for the complete architecture and is bound to its exact digest; ADH-2026-014/015/016 are retained only as historical provenance and excluded from downstream instructions. FEATURE-0012's six ScopeKind values are preserved; ServiceInstance is a typed subject under its governing scope. FEATURE-0014 may consume the inherited grammar but must not redefine DecisionRecord, AuditEvent, ScopeKind, or metadata.scopeRef. |
| FEATURE-0014 | Phase 2 | Implemented and Merged | TENET-008 | DEC-0024; F14-AD-001..F14-AD-021 | RFC-0024 | Final feature gate passed; PR #16 merged 2026-07-30 as commit `1ed47ac` | ADH-2026-018 remains controlling and ADH-2026-019 clarifies optional geographic descriptors. Owns exactly Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, and InfrastructureStack; reuses FEATURE-0012 grammar; FEATURE-0013 adoption is NOT_APPLICABLE; FEATURE-0015/0016/0053 semantics remain excluded. |
| FEATURE-0015 | Phase 2 | Planned | TENET-009 | DEC-0024, DEC-0025 | RFC-0024 | Pending | ResourcePool and ProviderCapability |
| FEATURE-0016 | Phase 2 | Planned | TENET-017 | DEC-0018 | RFC-0021 | Pending | Adapter boundary foundation |
| FEATURE-0017 | Phase 2 | Planned | TENET-007, TENET-017 | DEC-0020 | RFC-0025 | Pending | Policy evaluation abstraction |
| FEATURE-0023 | Phase 2 | Planned | TENET-009, TENET-012 | DEC-0026 | RFC-0026 | Pending | PlacementDecision v0 |
| FEATURE-0024 | Phase 2 | Planned | TENET-017 | DEC-0021 | RFC-0027 | Pending | Plugin taxonomy |
| FEATURE-0026 | Phase 2 | Planned | TENET-012, TENET-013 | DEC-0019 | Multiple | Pending | Phase 2 integration demo |

## Rule

Each implemented feature must reference at least one accepted decision or architecture document unless it is purely mechanical maintenance.
