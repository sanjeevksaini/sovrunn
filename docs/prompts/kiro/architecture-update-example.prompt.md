# Example Kiro Architecture Update Request

Use `docs/prompts/kiro/architecture-update.prompt.md`.

Approved handoff:

```text
Decision title: FEATURE-0017 PolicyEngineAdapter scope
Classification: Clarification
Approval status: Approved
Decision: FEATURE-0017 defines PolicyEngineAdapter, PolicyEvaluationRequest,
PolicyEvaluationResult, OPA adapter placeholder, and Cedar adapter placeholder.
FEATURE-0017 does not implement full OPA/Cedar integration in Phase 2.
```

Task:

Validate this handoff against the Sovrunn Architecture Operating System and update the impacted files:

- `docs/architecture/policy-evaluation-abstraction.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `.kiro/specs/policy-evaluation-abstraction/requirements.md`
- `.kiro/specs/policy-evaluation-abstraction/design.md`
- `.kiro/specs/policy-evaluation-abstraction/tasks.md`

Do not update Go code.
Do not implement full OPA/Cedar integration.
Do not change `CURRENT_ARCHITECTURE_BASELINE.md` unless validation shows a baseline update is required.
