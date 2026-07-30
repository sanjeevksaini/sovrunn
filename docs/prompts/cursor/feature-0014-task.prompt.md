# FEATURE-0014 guarded Cursor task

Implement only `FEATURE-0014` Task `{{TASK_ID}}`.

## Closed execution boundary

- The approved task block below is the complete implementation authority for this run.
- Do not read other specification, architecture, handoff, phase, RFC, roadmap, context-summary, or feature files.
- Generic always-on Cursor rules describe broad Phase 1 expectations. Where they request registries, handlers, controllers, persistence, or unrelated tests, they do not apply to FEATURE-0014; this approved task boundary is controlling.
- Do not modify `requirements.md`, `design.md`, `tasks.md`, architecture, prompts, validators, automation scripts, or automation configuration.
- Write only the paths in **Writable paths**. Stop with `ARCHITECTURE_DECISION_REQUIRED` if the task cannot be completed within them.
- Implement no future task and no helpful adjacent refactor.
- Reuse the listed repository exemplars; do not search broadly for “similar” implementations.
- Preserve exactly five kinds: `Provider`, `ProviderLocation`, `ProviderDatacenter`, `DatacenterFailureDomain`, and `InfrastructureStack`.
- Reuse FEATURE-0012 grammar. Do not redefine metadata, scope, references, conditions, Problem Details, limits, or schema conventions.
- Do not implement a repository, registry, persistence/storage mechanism, state-provider or resolver interface, live handler/CRUD path, controller, reconciler, adapter/provider integration, connectivity model, capability/capacity/placement model, decision/audit/operation model, or provisioning runtime.
- Stateful behaviors remain `CONTRACT_ONLY / NO_TASK`. Tests may use explicitly supplied in-memory values and fixtures only; they must not create production state machinery.
- `FEATURE-0013` is `NOT_APPLICABLE`. Do not load or reuse its semantic artifacts.
- Do not add dependencies, weaken tests, add unfinished-work markers, or commit generated/build artifacts.
- Treat the repository as potentially containing unrelated user work; do not overwrite or revert it.

## Writable paths

{{WRITABLE_PATHS}}

## Required workflow

1. Read only the embedded context below and the writable files if they already exist.
2. Confirm the task can be completed without a new semantic or architectural choice.
3. Implement the smallest complete change with idiomatic Go and deterministic behavior.
4. Add the exact positive, negative, boundary, concurrency, and isolation tests required by the task—only where applicable.
5. Run the task verification commands from the task block, followed by the supplied validation command.
6. Run `git diff --name-only HEAD` and stop if any changed path is outside **Writable paths** or the pre-existing automation-state change.
7. Report files changed, tests run/results, boundary confirmation, security considerations, and any blocker. Do not commit.

## Approved task block

{{TASK_BLOCK}}

## Minimal hashed context

The following files are the complete semantic and engineering context allowlist for this task. Their hashes are recorded in `{{CONTEXT_MANIFEST}}`.

{{CONTEXT_CONTENT}}

## Final validation command

```bash
make feature-0014-cursor-boundary-check TASK={{TASK_ID}}
```
