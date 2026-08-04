# Feature Traceability Matrix

This matrix links features to tenets, decisions, RFCs, tests, and gate status.

## Completed Features (Phase 1 and Phase 2 History)

| Feature | Phase | Status | Tenets | Decisions | RFCs | Tests/Gates | Notes |
|---|---|---|---|---|---|---|---|
| FEATURE-0011 | Phase 2 | Implemented | TENET-017 | DEC-0026, DEC-0036 | RFC-0021 | Feature gate passed | Reuse assessment standard merged; ADH-2026-011 |
| FEATURE-0012 | Phase 2 | Implemented and Merged | TENET-015 | DEC-0026, DEC-0027, DEC-0036 | RFC-0022 | Merged through PR #14 | ADH-2026-012, ADH-2026-013 |
| FEATURE-0013 | Phase 2 | Implemented and Merged | TENET-012, TENET-013 | DEC-0026, DEC-0036 | RFC-0023 | Merged through PR #15 | ADH-2026-017 controlling |
| FEATURE-0014 | Phase 2 | Implemented and Merged | TENET-008 | DEC-0024; F14-AD-001..F14-AD-021 | RFC-0024 | Merged through PR #16 | ADH-2026-018/019 controlling; alpha model |

## Phase 2R Rebaselined Features

| Feature | Phase | Status | Decisions | Controlling Source | Notes |
|---|---|---|---|---|---|
| FEATURE-0015 | Phase 2R | Architecture ready; Kiro requirements not started | DEC-0037, DEC-0038, DEC-0041, DEC-0042, DEC-0054, DEC-0058 | ADH-2026-042, ADH-020/024/025/037/041; `FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md`; VS0-SCHEMA-008..015/060/061; VS0-WRITER-001/003/004/020; VS0-STATE-001/010/011; VS0-CF-MIG01/MIG02/X03 | Closed canonical cloud identity and append-only alpha migration boundary; manifest prepared, architecture human gate still required |
| FEATURE-0016 | Phase 2R | Planned | DEC-0036, DEC-0042 | ADH-2026-042 | Adapter boundary and ExecutionTarget qualification |
| FEATURE-0017 | Phase 2R | Planned | DEC-0028, DEC-0043 | ADH-2026-042 | Policy evaluation with DecisionRecord linkage |
| FEATURE-0018 | Phase 2R | Planned | DEC-0050 | ADH-2026-042 | Governance, IAM, approval, exception |
| FEATURE-0019 | Phase 2R | Planned | DEC-0041, DEC-0043, DEC-0055 | ADH-2026-042 | Sovereignty facts and evidence |
| FEATURE-0020 | Phase 2R | Planned | DEC-0050 | ADH-2026-042 | EffectiveGovernanceContext resolution |
| FEATURE-0021 | Phase 2R | Planned | DEC-0038, DEC-0039, DEC-0047 | ADH-2026-042 | CloudEnrollment, entitlement, quota |
| FEATURE-0022 | Phase 2R | Planned | DEC-0040, DEC-0044, DEC-0049 | ADH-2026-042 | ServiceTypeDefinition, ServiceOffering, ServiceRequirementSet |
| FEATURE-0023 | Phase 2R | Planned | DEC-0043, DEC-0045, DEC-0046 | ADH-2026-042 | Sovereignty and placement DecisionRecord profiles |
| FEATURE-0024 | Phase 2R | Planned | DEC-0029 | ADH-2026-042 | Plugin taxonomy and synthetic execution |
| FEATURE-0025 | Phase 2R | Planned | DEC-0043 | ADH-2026-042 | AI-readable bounded explanation |
| FEATURE-0026 | Phase 2R | Planned | DEC-0058 | ADH-2026-042, VS-000, VS-000-CONTRACT-1 | Slice 0 integration and conformance; exact test contracts in `VS-000-contract-registry.yaml` |

## Rule

Each implemented feature must reference at least one accepted decision or architecture document unless it is purely mechanical maintenance.
