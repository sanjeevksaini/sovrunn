# Sovrunn Engineering Steering

## Source Documents

```text
docs/engineering/go-coding-guardrails.md
docs/engineering/ai-context-loading-standard.md
docs/engineering/ai-controlled-development.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
docs/engineering/context-engineering-standard.md
docs/features/FEATURE_SEQUENCE.md
docs/phase2/PHASE2_FEATURE_SEQUENCE.md
docs/phase2/PHASE2R_REBASELINE.md
docs/architecture/canonical/sovrunn-finalized-data-model.md
```
Follow docs/engineering/ai-context-loading-standard.md for context selection.

## Language

Use Go for Phase 2R platform core.

## Implementation Style

Prefer:

```text
simple code
explicit structs
small packages
deterministic validation
testable registries
clear API handlers
minimal dependencies
structured errors
```

Avoid:

```text
premature abstraction
large frameworks
implicit magic
future features hidden in current implementation
ResourcePool or ProviderCapability as active concepts
generic Provider as combined owner/operator
six-scope vocabulary in new code
```

## Package Direction

```text
cmd/sovrunn-api/
internal/api/
internal/audit/
internal/config/
internal/health/
internal/operation/
internal/registry/
internal/resources/
internal/server/
internal/validation/
tests/integration/
```

## Testing Rule

Every feature needs resource validation tests, registry tests, API handler tests, and error-path tests.

Run:

```bash
make fmt
make test
make vet
```

## Drift Checks

Phase 2R drift checks must verify:

- No active authority uses ResourcePool or ProviderCapability as mandatory
- No active authority uses generic Provider as combined concept
- Seven canonical scope kinds are used
- CloudEnrollment and CloudProviderParticipation are distinct
- ServiceClass has no active canonical authority
- Customer APIs have no provider-native objects
- No Go code change in architecture-only updates
