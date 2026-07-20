# Current Decision Summary

This file summarizes the currently binding architecture decisions. The detailed source of truth remains `docs/decisions/DECISION_INDEX.md` and individual DEC/RFC files.

## Accepted Decisions

- Sovrunn is a sovereign cloud-native PaaS platform.
- SDE is a future managed service inside Sovrunn.
- Reuse before build is mandatory (`DEC-0026`).
- Phase 2 is model/decision/audit/adapter only (`DEC-0027`).
- Policy evaluation must use a policy-engine abstraction (`DEC-0028`).
- Provider-neutral resource model is required before provider integrations.
- ResourcePool is the placement boundary (`DEC-0032`).
- ProviderCapability is the compatibility boundary (`DEC-0033`).
- PlacementDecision is required before provisioning (`DEC-0034`).
- Plugin taxonomy has three planes: provider/substrate, service management, service runtime (`DEC-0029`).
- MVP is governed PostgreSQL PaaS placement and provisioning on one substrate (`DEC-0030`).

## Decisions Requiring Future Revalidation

- First real policy adapter choice: OPA expected, Cedar evaluated later for authorization.
- First PostgreSQL runtime reuse choice: CloudNativePG, Crunchy, or Helm.
- Operation engine backend: simple v0 first; Temporal/Argo later.
- Secret backend: Kubernetes Secret for local MVP; Vault/External Secrets later.
