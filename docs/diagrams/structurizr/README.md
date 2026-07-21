# Structurizr Integration

This directory contains the Sovrunn C4 architecture-as-code model.

## Purpose

`workspace.dsl` is the durable diagram model for approved Sovrunn architecture.
It complements Markdown architecture docs, DEC/RFC records, and the Architecture Operating System.

Structurizr does not decide architecture. It visualizes approved architecture.

## Source-of-truth rule

When approved architecture changes any of the following, update `workspace.dsl`:

- platform/system boundaries,
- major containers,
- plugin boundaries,
- external OSS/reuse relationships,
- deployment/runtime relationships,
- major dynamic flows,
- ChatGPT -> Kiro -> Cursor architecture handoff workflow.

## Local usage

From the repo root:

```bash
make structurizr-lite
```

Then open:

```text
http://localhost:8080
```

The Makefile target mounts this directory into Structurizr Lite.

## Validation

Run:

```bash
make structurizr-check
```

The check verifies that `workspace.dsl` exists and optionally validates it when Structurizr CLI is installed.

## Cloud publishing

Optional later:

```bash
export STRUCTURIZR_WORKSPACE_ID="..."
export STRUCTURIZR_API_KEY="..."
export STRUCTURIZR_API_SECRET="..."
make structurizr-push
```

Do not commit Structurizr credentials.
