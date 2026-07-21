---
doc_type: repo_hygiene_review
title: Operational Cleanup Review
status: completed
phase: 2
ai_load_priority: low
---

# Operational Cleanup Review

## Purpose

This review records cleanup applied before using the regenerated Sovrunn repo as a Phase 2 Architecture Operating System baseline.

## Cleanup applied

- Removed generated/site/archive artifacts from source tree:
  - `site/`
  - `docs/generated-prompts/`
  - `docs/Archive.zip`
- Added `.gitignore` for generated prompts, site output, generated context packs, logs, automation outputs, archives, and macOS metadata.
- Replaced Mac-specific script shebangs with portable `#!/usr/bin/env bash`.
- Added `docs/engineering/go-version-standard.md` and normalized Go verification references away from stale Go 1.21 wording.
- Added `docs/features/FEATURE_INDEX.md` to map feature IDs to Kiro spec slugs.
- Replaced duplicated all-phase roadmap copy with a pointer file and kept `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md` canonical.
- Added individual decision records for major accepted decisions DEC-0026, DEC-0027, DEC-0028, DEC-0029, DEC-0030, and DEC-0036.
- Updated `mkdocs.yml` navigation to include Architecture Operating System, Phase 2, roadmap, traceability, handoff prompts, Structurizr, and RFC-0021 through RFC-0029.
- Hardened `scripts/phase2-scope-check.sh` to allow blocked future-scope phrases in Non-goals, Out of Scope, Deferred, or explicit negated statements.
- Hardened `scripts/feature-gate.sh` to use `docs/features/FEATURE_INDEX.md` instead of guessing Kiro spec paths from feature IDs.
- Strengthened Cursor/Kiro/reviewer prompt observability expectations.
- Corrected Structurizr future relationship tagging and preserved the `exports/` folder with `.gitkeep`.

## Validation performed

```text
bash -n scripts/*.sh
scripts/architecture-handoff-check.sh docs/reviews/architecture-decision-handoffs/ADH-EXAMPLE-FEATURE-0017-policy-engine-adapter.md
scripts/structurizr-check.sh
scripts/context-pack.sh /tmp/sovrunn_context_test.md
scripts/phase2-scope-check.sh FEATURE-0011
scripts/feature-gate.sh FEATURE-0001
```

Structurizr CLI and MkDocs were not installed in the regeneration environment, so full Structurizr syntax validation and MkDocs build validation must be run on the developer machine when those tools are available.

## Remaining recommended manual checks

```bash
make structurizr-lite
mkdocs build --strict
```

