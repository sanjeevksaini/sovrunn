---
doc_type: engineering_standard
title: Go Version Standard
status: active
phase: 2
ai_load_priority: high
---

# Go Version Standard

## Purpose

This file is the single source of truth for Go version expectations in Sovrunn prompts, scripts, and verification commands.

## Current standard

```text
Minimum supported Go version: 1.22
Default Docker verification image: golang:1.22
```

## Rules

- Do not hardcode Go versions independently in prompts or scripts.
- Scripts should use `GO_DOCKER_IMAGE` when running Docker-based Go verification.
- The default value is `golang:1.22` unless the repository `go.mod` or an accepted decision updates this standard.
- Kiro and Cursor prompts must refer to this file instead of embedding stale version requirements.

## Future update process

Changing the Go version standard requires:

1. updating this file,
2. updating `Makefile` defaults if required,
3. validating local and CI verification commands, and
4. recording the change in the relevant feature or decision review when the change is material.
