# FEATURE-0012 Task 17.1 — Baseline Protected-Review Evidence

Feature: FEATURE-0012 — API, Resource Naming, Status, and Validation Standard
Task: 17.1 Configure CODEOWNERS for baseline protected review
Design: D-11
Requirements: F12-EVOLVE-002, F12-VERIFY-001(10)

Evidence collected at (UTC): 2026-07-23T16:43:18Z
Repository: `https://github.com/sanjeevksaini/sovrunn`
Default branch inspected: `main`

## CODEOWNERS inspection

| Check | Result |
|---|---|
| Existing `.github/CODEOWNERS` before Task 17.1 | Absent |
| Existing root `CODEOWNERS` before Task 17.1 | Absent |
| Existing in-repo CODEOWNERS owner pattern | None (no prior CODEOWNERS file) |
| Existing GitHub repository owner identity | `sanjeevksaini` (from `git remote` / GitHub repo owner) |

Task 17.1 created `.github/CODEOWNERS` and assigned baseline paths to the
existing repository owner identity `@sanjeevksaini`. No team or user
identity was invented for this task.

### CODEOWNERS entry added

```text
api/schemas/baseline/** @sanjeevksaini
```

This path covers baseline schema snapshots, `BASELINE_MANIFEST.json`, and
`BASELINE_APPROVALS.json` under `api/schemas/baseline/`.

## Branch-protection evidence (external; not configured by repository code)

Branch protection is external GitHub settings evidence. This task
**collected** the following observed state via the GitHub API. Repository
files in this change **did not configure** branch protection and must not be
read as having done so.

Observed `main` branch protection (source:
`GET /repos/sanjeevksaini/sovrunn/branches/main/protection`):

| Setting | Observed value |
|---|---|
| Pull-request reviews required | yes |
| Required approving review count | 1 |
| Dismiss stale reviews | true |
| Require conversation resolution | true |
| Require code owner reviews | **false** |
| Enforce admins | false |
| Allow force pushes | false |
| Allow deletions | false |
| Required status checks (strict) | true (no required contexts/checks listed) |
| Required signatures | false |

### PENDING_HUMAN_REVIEW items (external settings)

The following remain **PENDING_HUMAN_REVIEW** because they are external
GitHub settings that this coding agent must not change:

1. Enable `require_code_owner_reviews` on `main` so CODEOWNERS for
   `api/schemas/baseline/**` is enforced on protected merges.
2. Confirm whether additional required status-check contexts should be
   attached to the existing strict status-check rule.
3. Confirm whether `enforce_admins` should be enabled for baseline
   governance hardening.

Status for those external settings: **PENDING_HUMAN_REVIEW**

## Explicit non-claims

- This evidence document does **not** claim that repository code configured
  branch protection.
- This evidence document does **not** write an approval token, invent a
  reviewer, invent a review date, or accept residual risk on behalf of a
  human (those remain Task 17.3 / human gates).
- `BASELINE_MANIFEST.json` remains an integrity mechanism; protected review
  (CODEOWNERS + branch protection) remains the human governance boundary.
