---
doc_type: threat_and_abuse_ledger
feature: FEATURE-0018
status: architecture-approved-security-reviewed
updated: 2026-08-30
---

# FEATURE-0018 Threat and Abuse Ledger

This architecture-stage ledger supplies review input; it is not the independent
security review and does not prove implementation resistance.

| ID | Threat / abuse case | Preventive architecture control | Required detection/proof | Residual owner |
|---|---|---|---|---|
| F18-THR-01 | Raw IdP group claim grants access | Only trusted provisioning into direct canonical Membership; external claims grant nothing | Raw-claim deny vs provisioned-group allow/deny pair | Future IdP/provisioner + FEATURE-0018 |
| F18-THR-02 | Nested/dynamic group hides beneficiary | Both forms rejected; RoleHolderRef group must have direct current members | Nested/dynamic schema and runtime denial | FEATURE-0018 |
| F18-THR-03 | Cross-tenant/scope confused deputy | Exact target binding, scope containment, resource narrowing and grantor ceiling | Seven-scope cross-tree and target mismatch matrix | FEATURE-0012/0018 |
| F18-THR-04 | Admin grants more than they hold | Publication-time `roleassignment.grant` reach plus one independent delegable witness per action covering scope/target/resource/time/delegation; no cross-assignment synthesis | Over-action/scope/duration/delegation, resource-A-to-B and resource-to-scope negatives | FEATURE-0018 |
| F18-THR-05 | Workflow writes an assignment directly | RoleAssignment controller is sole spec/status/lifecycle/effect publisher; every caller or workflow submits intent only | Writer/import/call-path ledger and wrong-writer proof | FEATURE-0018 |
| F18-THR-06 | Approval replay creates later authority | Approval current at effect publication; immutable provenance afterward; effect freshly authorized | Expired/replayed/changed-state approval cases | FEATURE-0018 |
| F18-THR-07 | Self-dealing via group/workload alias | Approval and access-review conflict sets cover direct Human holders, direct Workload/System responsible Humans, and, for AccessGroup-held authority, the owner, direct Human members and responsible Humans of all direct Workload/System members; unresolved review expansion fails closed for Retain/Replace and the complete set is rechecked at atomic decision/effect publication | Direct/group/group-workload-responsible-party, unresolved-expansion and replacement-union SoD cases | FEATURE-0018 |
| F18-THR-08 | JIT activation races deadline/replay | Atomic publication strictly before activationDeadline; exactly one TimeBound assignment | Before/at/after deadline and concurrent activation proof | FEATURE-0018 |
| F18-THR-09 | JIT or break-glass becomes weak/permanent/bypassable | Uniform phishing-resistant AAL2-or-higher activation floor; sole break-glass bypass rule, individual eligibility, duration ceilings, review/audit | JIT/break-glass assurance, duration, review and audit-failure abuse suite | FEATURE-0018 + external authenticator |
| F18-THR-10 | Revoked/suspended access remains usable | Current dependencies, exact time boundaries and authorization re-evaluation | Revoke/use, suspend/use and exact-expiry races | FEATURE-0018 |
| F18-THR-11 | Incomplete telemetry causes unsafe revocation or retention | Qualified immutable UsageEvidenceSummary; absence is not non-use; explicit overdue rule | Complete/stale/absent/incomplete workload cases | FEATURE-0018/evidence owner |
| F18-THR-12 | Exception silently broadens approved proposal | Intersect allowed sets, union mandatory compensating controls, minimum duration; reject rather than rewrite | Pairwise overlap/narrowing/terminal-retry cases | FEATURE-0018/0020 boundary |
| F18-THR-13 | Audit outage allows sensitive mutation | Required audit-obligation acceptance before publication; atomic failure publishes nothing | Failure injection for every authorization-changing boundary | FEATURE-0013/0018 |
| F18-THR-14 | Authorization result reused as bearer token | Result valid only at evaluation instant and bound to exact input/provenance | Cross-request/target replay denial | FEATURE-0018 |
| F18-THR-15 | Native cloud Allow bypasses Sovrunn or vice versa | Two-plane intersection; neither Allow substitutes and either Deny blocks | Four-way Sovrunn/native truth table and no-native-type sentinel | FEATURE-0018 + future execution owner |
| F18-THR-16 | Provider-native credential leaks into canonical evidence | No native credential/type/output; protected refs handled by later owner | Schema/import/log/redaction sentinels | Future adapter/execution owner |
| F18-THR-17 | Owner/responsible party disappears | Privilege-increasing actions block; privilege-reducing/recovery actions remain; existing access changes only through explicit dependencies | Inactive-owner increase/reduction/recovery matrix | FEATURE-0018 |
| F18-THR-18 | Idempotency collision crosses actor/target | Key binds actor, operation, exact target, scope and digest; replay rechecks visibility/authorization | Same/different digest, actor, target and changed-access cases | FEATURE-0018 |
| F18-THR-19 | Denial leaks target or beneficiary existence | Structural/local checks precede safe resolution; protected details redacted; inherited safe outcomes | Cross-scope malformed/unauthorized equivalence and redaction tests | FEATURE-0012/0018 |
| F18-THR-20 | Future governance or adapter behavior becomes accidental authority | Closed inventories/components/exclusions and architecture-drift sentinels | Forbidden-field/type/dependency scan | Architecture owner |
| F18-THR-21 | Temporary or broad Membership administration launders durable group-derived authority | Every AccessGroup assignment-effect expansion derives the exact newly enabled/extended envelope and requires applicable Membership authority, `roleassignment.grant`, and non-synthesized per-action scope/target/resource/validity witnesses; Standing effect requires Standing witnesses and trusted provisioner has no bypass | Empty/non-empty envelope, JIT-to-Standing, resource-A-to-B, stale refresh, self/accomplice add, provisioner overreach and membership/assignment race cases | FEATURE-0018 |
| F18-THR-22 | Membership or group-derived authority launders approval, review, JIT or break-glass eligibility | One closed non-empty `EligibilityRef` OR contract accepts exact Human PrincipalRefs only; AccessGroupRef, Membership and group-derived authority reject; eligibility is independently action-authorized and rechecked at decision/effect publication | Empty-group policy/rule, self/accomplice, reviewer kind/scope/target/action and publication-race cases | FEATURE-0018 |
| F18-THR-23 | Ordinary role grant or later rule publication launders approval, review, JIT or break-glass eligibility | RoleDefinition and RoleAssignment are prohibited EligibilityRef kinds; existing holders never gain eligibility from later policy/rule publication; exact named Human, applicable policy/rule and separate action are all required and rechecked | Self/accomplice role grant, broadly assigned role, later rule, approval/review/JIT/break-glass and concurrent role/rule publication cases | FEATURE-0018 |

No listed threat is accepted as an unbounded risk. Runtime residual risk cannot
be closed until implementation conformance, race testing, operational identity
and native-IAM integration evidence exist.
