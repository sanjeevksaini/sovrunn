# ACR-2026-002: FEATURE-0018 Governance, IAM, Approval and Exception Foundation

## Title

Adopt the canonical FEATURE-0018 control-plane authorization, approval, access-review and exception foundation.

## Change Type

Extension and new decision.

## Current Approved Position

The Phase 2R baseline names FEATURE-0018 and contains preliminary identity and
governance concepts, but its canonical inventory still binds assignments
directly to principals, omits access-only groups and several retained workflow
contracts, and does not close the authorization algebra, writers, lifecycle,
evidence, target binding, validation or native-IAM trust boundary.

## Proposed Change

Apply the jointly governed ADH-2026-070 package and DEC-0060 as the FEATURE-0018
architecture authority. The change:

1. adopts the closed resource and supporting-value inventories in F18-RD-02;
2. adds direct-principal and direct-member `AccessGroup` role holders without
   trusting external group claims or adding nested/dynamic groups;
3. defines immutable role versions, scoped role assignments, Standing and
   TimeBound validity, privileged/JIT/break-glass controls, access reviews,
   approvals and immutable exception evidence;
4. closes grant-union plus guardrail-intersection authorization, exact action
   target bindings, deterministic validation, safe denial, sole writers and
   audit-before-publication;
5. registers FEATURE-0018 adoption of FEATURE-0013 and additive metadata for
   the existing FEATURE-0016 ExecutionTarget actions without changing either
   dependency's semantic ownership;
6. limits `GovernanceProfile` v1 to FEATURE-0018-owned governance composition
   and preserves FEATURE-0020 ownership of assignment and effective resolution;
7. establishes the Sovrunn/native-IAM intersection boundary without creating,
   translating or operating provider-native IAM; and
8. preserves progressive disclosure so routine user journeys remain simple.

The exact behavioral authority remains F18-RD-01 through F18-RD-24 in the
content-bound ADH-2026-070 core and Appendices A and B. Repository summaries and
evidence ledgers must reference those decisions and must not become competing
semantic authorities.

## Reason

Enterprise Sovrunn operations require one provider-neutral, zero-trust control-
plane authorization foundation before later governance and service features can
be specified. Closing the model now prevents direct IdP-claim authorization,
multiple assignment writers, ambiguous scopes, approval-as-runtime-authority,
provider-native IAM leakage and future-feature ownership leakage.

## Alternatives Considered

- Direct-principal-only RBAC: rejected because it creates avoidable enterprise
  administration overhead and lacks auditable access-group ownership.
- Trust OIDC group claims directly: rejected because claim freshness,
  provenance and lifecycle would become implicit authorization authority.
- Provider-native IAM as Sovrunn's canonical model: rejected because it breaks
  provider neutrality and cannot replace Sovrunn control-plane authorization.
- Generic ABAC expressions or a new policy/workflow engine: deferred; they add
  complexity and duplicate FEATURE-0017 or future adapter ownership.
- One combined governance/workflow resource: rejected because it conflates
  immutable definitions, requests, decisions, grants and review evidence.

## Reuse Assessment

The complete FEATURE-0011-format assessment is recorded in
`docs/features/FEATURE-0018-governance-iam-approval-exception-foundation.md`.
The proposed disposition is Build only for provider-neutral Sovrunn domain
contracts and deterministic fake controllers; Reuse for FEATURE-0012,
FEATURE-0013, FEATURE-0016 and FEATURE-0017 foundations; and Deferred Wrap for
future IdP, policy-engine and workflow integrations. No external candidate owns
Sovrunn's complete canonical authorization and evidence boundary.

## Impacted Docs

- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/decisions/DECISION_INDEX.md`
- canonical model, canonical contract catalog and glossary
- FEATURE-0013 adoption and FEATURE-0016 action-registration authorities
- Phase 2R sequence, feature index, roadmap and traceability matrices
- FEATURE-0018 architecture, digest, industry matrix, feature contract, reuse
  evidence, control manifest and architecture-readiness evidence

The active architecture baseline/version files are not changed until DEC-0060,
the exact applied package and the final architecture evidence receive explicit
human approval.

## Impacted Features

- FEATURE-0011: format reused unchanged.
- FEATURE-0012: common resource/scope/validation authority reused unchanged.
- FEATURE-0013: three adopting profiles and a bounded AuditEvent taxonomy are
  registered; the envelope remains FEATURE-0013-owned.
- FEATURE-0016: three existing action meanings gain exact target-binding and
  intrinsic-class registrations; meaning and lifecycle remain FEATURE-0016-owned.
- FEATURE-0017: seam and fake reused within the approved v1 limit.
- FEATURE-0018: sole current implementation scope.
- FEATURE-0019 onward: exclusions only; no future behavior is approved.

## RFC Disposition

No existing RFC is replaced. RFC-0012's enterprise governance intent,
RFC-0021 reuse governance, RFC-0022 resource/API foundations and RFC-0023
decision/audit envelope remain valid and are reused without reinterpretation.
The closed FEATURE-0018 semantics are sufficiently governed by ACR-2026-002,
DEC-0060, the canonical architecture and ADH-2026-070; a second focused RFC
would duplicate those authorities. A later real IdP, policy/workflow or native-
IAM integration must make its own RFC/ADH disposition.

## Backward Compatibility Impact

This is a pre-implementation canonical contract correction. No live customer
state or supported FEATURE-0018 API exists. Existing lower-authority
experimental descriptions are reconciled or retired; no compatibility alias,
dual writer or migration runtime is introduced.

## Phase Impact

Allowed in Phase 2R as deterministic contracts, in-memory fixtures and local
conformance only. No real IdP, policy engine, workflow engine, provider IAM,
persistence, provider execution or deployment topology is authorized.

## Security Impact

Positive: default deny, exact target binding, scoped least privilege, separation
of duties, explicit responsibility, time-bounded privilege, immutable evidence,
audit-before-publication, revocation, review, fail-closed dependency handling,
and independent native-IAM enforcement. Independent security-review renewal-05
passed the exact 32-file payload with no blocking findings.

## Observability Impact

FEATURE-0018 registers its bounded AuditEvent taxonomy and three DecisionProfile
adoptions while preserving FEATURE-0013 ownership. Protected audit-obligation
acceptance precedes publication; operational telemetry never substitutes for
authoritative audit or decision evidence.

## Decision Required

Yes.

## Recommendation

Accepted after review of the exact applied package, approval evidence and
independent security-review result. DEC-0060 and the active architecture
baseline metadata are reconciled in the same controlled approval change.

## Approval

- Proposed by: Codex, applying the architecture-owner-directed ADH-2026-070 handoff
- Reviewed by: Independent security-review renewal-05; no blocking findings
- Approved by: Sanjeev Kumar, Sovrunn Architecture Owner
- Date: 2026-08-30
- Approved input: final-approval manifest SHA-256 `9f69ea2fb89ba04aa253a37a0cf5f63124753c2b1e9cf08cf0a5c028369d8efc`
- Security evidence: renewal-05 SHA-256 `2f95272aebb23fbfdca507de6434461f5e66378ddc206e9ff5b382848a905ccf`
- Scope: Architecture and reuse approval only; requirements, design, tasks and implementation remain separately gated.
