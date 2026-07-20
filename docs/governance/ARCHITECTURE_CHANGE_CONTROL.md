# Architecture Change Control

This document defines how Sovrunn architecture changes over time.

## Purpose

Sovrunn architecture must evolve deliberately, not accidentally through chat suggestions, generated prompts, or implementation drift.

## Change Classification

Every proposed architecture change must be classified as one of:

- clarification: explains existing architecture without changing it,
- extension: adds compatible detail without replacing an approved decision,
- correction: fixes inconsistency or error,
- replacement: changes or supersedes an approved decision,
- new decision: introduces a new approved architecture direction.

## Approval Rules

Clarifications may update explanatory docs.

Extensions require review against the current baseline and may require DEC/RFC updates.

Corrections require affected docs to be listed and reviewed.

Replacements and new decisions require:

- architecture change request,
- impacted docs list,
- impacted features list,
- compatibility impact,
- phase impact,
- human approval,
- updated DEC/RFC records,
- updated current baseline if accepted.

## Architecture Change Request

Use `docs/templates/ARCHITECTURE_CHANGE_REQUEST.md` for significant changes.

## Non-Authoritative Sources

The following do not change architecture by themselves:

- ChatGPT answers,
- generated prompts,
- local notes,
- roadmap placeholders,
- implementation code that bypassed review.

## Authoritative Sources

The following may define approved architecture:

- current architecture baseline,
- accepted DEC files,
- approved RFC files,
- approved architecture docs,
- phase scope docs.
