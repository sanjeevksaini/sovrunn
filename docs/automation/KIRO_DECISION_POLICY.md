# Sovrunn Kiro Decision Policy

This policy controls how Sovrunn Feature Factory answers Kiro decisions during spec generation.

## Default workflow

For every Sovrunn feature:

1. Build a Feature.
2. Generate `requirements.md` only.
3. Wait for reviewer approval token `APPROVED_FOR_DESIGN`.
4. Generate `design.md` only.
5. Wait for reviewer approval token `APPROVED_FOR_TASKS`.
6. Generate `tasks.md` only.
7. Wait for reviewer approval token `APPROVED_FOR_CURSOR`.
8. Stop before implementation.
9. Cursor implements tasks one at a time.

## Auto-accept decisions

- Build a Feature.
- Generate `requirements.md` only when the current stage is requirements.
- Generate `design.md` only when `APPROVED_FOR_DESIGN` exists.
- Generate `tasks.md` only when `APPROVED_FOR_TASKS` exists.
- Proceed to Cursor only when `APPROVED_FOR_CURSOR` exists.

## Auto-reject decisions

- Implement code in Kiro.
- Modify source files during requirements/design/tasks generation.
- Create source files during spec generation.
- Add ServiceInstance, ServiceBinding, provisioning, persistence, UI, auth/RBAC, billing, marketplace, ServiceOps execution, or AI automation unless explicitly listed in the current feature scope.
- Add external dependencies unless explicitly approved in the design.
- Use Go 1.22 wildcard routing.
- Move to the next stage without an exact approval token.

## Pause decisions

Pause for architecture review when Kiro asks about:

- Global vs tenant/project scope.
- Plugin/capability ownership and boundaries.
- Runtime execution mode: in-process, sidecar, remote.
- Any cross-feature dependency not already in requirements.
- Any change that affects Sovrunn architecture, security, tenancy, or product scope.

## Standard response format

```text
Decision: <Kiro question>

Answer: <Accept | Reject | Pause>

Reason:
<short reason>

Instruction to Kiro:
<exact instruction>
```

## Golden rule

Kiro may decide wording, ordering, and formatting. Kiro may not decide architecture scope, implementation scope, or next-stage approval.
