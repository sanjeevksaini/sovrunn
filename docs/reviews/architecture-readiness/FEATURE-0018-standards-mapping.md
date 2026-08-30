---
doc_type: security_standards_mapping
feature: FEATURE-0018
status: architecture-reviewed-and-approved
updated: 2026-08-30
---

# FEATURE-0018 Security Standards and Shared-Responsibility Mapping

## 1. Claim boundary

This mapping validates architecture coverage and identifies downstream proof.
It is not a certification, legal conclusion, formal ISO conformity assessment,
or claim that operational controls already exist. F18-RD-01..24 remain the
architecture authority.

Verified framework baselines for this review are NIST SP 800-207 (final), NIST
SP 800-63-4 (final, July 2025), NIST SP 800-53 Rev. 5, NIST SP 800-162 (2019
update), CIS Controls v8.1, CSA CCM 4.1 and OWASP ASVS 5.0.0. Official reference
links are retained so a reviewer can recheck currency.

## 2. NIST mapping

| Framework/control theme | FEATURE-0018 architectural response | Proof / owner |
|---|---|---|
| [NIST SP 800-207 Zero Trust](https://csrc.nist.gov/pubs/sp/800/207/final): no implicit trust; per-resource decision | Exact actor/action/target/scope; non-bearer result; every operation re-evaluates current dependencies; default deny | Authorization truth table, missing/stale/deny cases; FEATURE-0018 |
| [NIST SP 800-207A](https://csrc.nist.gov/pubs/sp/800/207/a/final): identity-centric multi-cloud enforcement | Human/Workload/System PrincipalRef; provider-neutral actions; Sovrunn/native-IAM intersection | Two-plane matrix and no-native-type sentinel; FEATURE-0018 plus future execution owner |
| SP 800-53 AC-2 account lifecycle | Membership provenance/freshness/Guest expiry, responsible parties, review/revocation, and grant-envelope validation when group assignment effect is restored or extended; Membership cannot create approval/review/privileged eligibility | Lifecycle/expiry/review, empty/non-empty envelope and empty-group eligibility-denial proof; FEATURE-0018; identity account lifecycle external |
| AC-3 access enforcement | Exact assignment/action/target/scope and fail-closed composition | Full allow/deny matrix; FEATURE-0018 |
| AC-5 separation of duties | Approval and access-review conflict sets span direct holders/requesters and, for AccessGroup-held authority, the owner, direct Human members and responsible Humans of direct Workload/System members; approver/reviewer eligibility is exact-Human PrincipalRef-only, action-independent and rechecked; roles, assignments, groups and Membership never qualify; unresolved review expansion fails closed for Retain/Replace at atomic decision/effect publication | Direct, indirect, group-workload-responsible-Human, unsupported eligibility-kind, replacement-union and unresolved-expansion SoD tests; FEATURE-0018 |
| AC-6 least privilege | Scoped roles, per-action scope/target/resource-reach/validity grantor dominance for direct grants and AccessGroup Membership assignment-effect expansion, monotonic privilege, exact-Human PrincipalRef-only eligibility, and JIT/break-glass TimeBound access | Overgrant/resource-reach/duration/delegation, role/rule eligibility laundering, group-derived eligibility, JIT-to-Standing laundering, provisioner-overreach and membership/assignment/rule-race negatives; FEATURE-0018 |
| IA-2/IA-5 authentication/assurance | External authenticated identity plus provider-neutral AssuranceEvidence; fresh phishing-resistant privileged assurance | Assurance fixtures and future authenticator/IdP proof; external identity owner + FEATURE-0018 validation |
| AU-2/AU-3/AU-9/AU-12 audit selection/content/protection/generation | Exact bounded taxonomy; FEATURE-0013 envelope/projection/protection; obligation accepted before publication | Taxonomy, redaction, tamper/append and audit-failure atomicity; FEATURE-0013/0018 |
| CA-7 continuous monitoring | Standing certification, immutable usage evidence and break-glass retrospective review | Review modes, evidence freshness/coverage and due-time proof; FEATURE-0018 |
| CM change control | Immutable published definitions, exact pins, supersession, suspension/restoration and retirement | Version/lifecycle/no-silent-migration proof; FEATURE-0018 |
| IR emergency response | Bounded break-glass with immediate audit and mandatory retrospective review | Emergency allow/deny/review/race proof; FEATURE-0018 |
| [NIST SP 800-63-4](https://csrc.nist.gov/pubs/sp/800/63/4/final) identity/authentication/federation | Identity proofing/authentication/federation remain external; canonical carrier records issuer/subject and assurance level/age/phishing resistance only | Paper OIDC/WebAuthn mapping and future adapter assessment |
| [NIST SP 800-162](https://doi.org/10.6028/NIST.SP.800-162) policy attributes | Actor, scope, target, time, assurance and context are explicit typed inputs, but v1 adds no generic ABAC language | Schema-absence and typed-input proof; FEATURE-0018 |

## 3. CIS Controls v8.1 mapping

The [official CIS Controls navigator](https://www.cisecurity.org/controls/cis-controls-navigator/v8)
identifies Controls 5 (Account Management), 6 (Access Control Management) and 8
(Audit Log Management) as the relevant operational families.

| CIS safeguard/theme | FEATURE-0018 coverage | Responsibility / residual gap |
|---|---|---|
| 5.1/5.5 account and service-account inventory | Principal refs, belonging, Workload/System responsible party and review evidence | FEATURE-0018 tracks authorization relationships; external identity owner inventories/authenticates accounts |
| 5.3 disable dormant accounts | Suspended/revoked/expired dependencies deny; usage evidence is never silently interpreted as non-use | Dormancy policy/account disablement external; F18 consumes current state and supports review remediation |
| 5.4 dedicated administrator privilege | Privileged actions are classified, individual, strongly authenticated and TimeBound where required | Persona/account separation is external; roles and enforcement are FEATURE-0018 |
| 6.1/6.2 grant and revoke processes | Exact grant intents, sole publishers, AccessGroup membership-enabled grant envelope, exact-Human PrincipalRef-only workflow eligibility, explicit revoke and audit | FEATURE-0018 |
| 6.5 MFA for administrative access | Fresh phishing-resistant AAL2-or-higher assurance for every JIT and break-glass activation | Authenticator enforcement external; evidence validation FEATURE-0018 |
| 6.6 authorization-system inventory | Closed resources/actions/targets/scopes and dependency ledgers | FEATURE-0018 plus platform operations |
| 6.8 RBAC and reviews | Immutable roles, scoped assignments, resource-reach-aware grantor ceiling, the same dominance when group membership activates assignments, non-empty exact-Human PrincipalRef-only reviewer eligibility with independent action checks, rejection of role/assignment qualification, and beneficiary-independent Standing certification including accountable Humans behind group-member workload/system identities | FEATURE-0018 |
| 8.1/8.2/8.5 audit process/collection/detail | Exact event taxonomy and protected actor/action/subject/scope/outcome/correlation evidence | FEATURE-0013/0018 |
| 8.9/8.10/8.11 centralize/retain/review | Retention/projection class refs and immutable linkage exist | Operational log platform, durations and review cadence remain deployment/governance responsibility |
| 8.12 service-provider logs | Native provider authorization remains independent and must be correlated by later execution owner | Future integration/execution feature; not FEATURE-0018 |

## 4. CSA CCM 4.1 shared responsibility

The [CSA Cloud Controls Matrix 4.1](https://cloudsecurityalliance.org/research/cloud-controls-matrix)
is used at domain/responsibility level; no unlicensed clause-level conformity is
claimed.

| CCM domain | CloudPlatform/Sovrunn responsibility | CloudProvider/native responsibility | Customer Organization responsibility |
|---|---|---|---|
| IAM | Canonical actions/roles, assignments, membership-enabled grant-envelope enforcement, exact-Human PrincipalRef-only eligibility, approvals, privileged flows, reviews and safe projections | Native controller identity, credentials, provider policy and native allow/deny | Membership intent, group ownership, scoped grants/revocations and reviews within delegated authority; roles, groups and trusted provisioning never create workflow eligibility |
| LOG | Canonical DecisionRecord/AuditEvent obligation and correlation | Native IAM/provider action logs and retention under provider controls | Authorized review and investigation within visibility policy |
| GRC | FEATURE-0018 rule envelopes/evidence; FEATURE-0020 later resolves effective governance | Provider control evidence and restrictions | Narrower customer governance, approval and exception proposals |
| IPY | Provider-neutral contracts; no native type in core | Adapter-native translation/compatibility under later owner | No direct native dependency required for ordinary journeys |
| IVS | No infrastructure execution authority | Least-privileged native enforcement and infrastructure operation | No implied native administrative authority from Sovrunn roles |

Native IAM Allow cannot replace Sovrunn Allow, and Sovrunn evidence cannot
override native Deny. Direct out-of-band IaaS administration remains provider-
owned and outside FEATURE-0018.

## 5. OWASP ASVS 5.0.0 applicability profile

The [OWASP ASVS project](https://owasp.org/www-project-application-security-verification-standard/)
identifies 5.0.0 as the latest stable release and recommends version-qualified
identifiers. FEATURE-0018 pins:

- Level 2 for every eventual FEATURE-0018 application/API surface;
- Level 3 for privileged-access, approval, exception, role/security
  administration and protected evidence surfaces;
- architecture-stage applicability to ASVS chapters 2 (validation/business
  logic), 4 (API), 6 (authentication), 7 (session), 8 (authorization), 9/10
  (token/OAuth/OIDC where a future adapter applies), 13 (configuration), 14
  (data protection), 15 (secure architecture) and 16 (logging/error handling).

Exact `v5.0.0-x.y.z` requirement selection belongs to requirements/conformance
after architecture approval. That selection may add verification evidence but
cannot change F18-RD semantics. No ASVS level claim is made before the exact
requirement matrix and executable proof pass independent review.

## 6. ISO/IEC boundary

ISO/IEC 27001/27002 identity, access control, privileged access, segregation,
logging, monitoring and review concepts are covered at architecture-theme
level. ISO/IEC 24760 identity lifecycle, ISO/IEC 29146 access-management and
ISO/IEC/IEEE 42010 architecture-description concepts inform separation and
traceability. Formal clause mapping, Statement of Applicability, audit sampling,
licensed text assessment and certification are explicitly outside this record.

## 7. Future protocol compatibility, paper only

| Standard | Compatible canonical seam | Prohibited architecture leakage |
|---|---|---|
| OIDC | issuer/subject → PrincipalRef; authentication/assurance → AssuranceEvidence | Raw token, claim-to-role authority, IdP-native role/group in canonical schema |
| SCIM | trusted provisioning → PrincipalRef/Membership/AccessGroup with provenance/freshness | Direct SCIM/external-group authorization, nesting/dynamic behavior |
| WebAuthn | phishing-resistance and assurance freshness → AssuranceEvidence | WebAuthn credential/factor payload in core |
| SPIFFE | workload identity → Workload PrincipalRef through future trust adapter | SPIFFE runtime, bundle or workload API in FEATURE-0018 |

## 8. Review outcome

Coverage passed independent architecture security review, with these
honest residuals: executable controls do not yet exist; real identity/native
IAM responsibilities require later adapter/execution evidence; operational log
retention/monitoring is deployment-owned; ISO conformity is unassessed. Status:
**ARCHITECTURE EVIDENCE REVIEWED AND APPROVED — EXECUTABLE PROOF PENDING**.
