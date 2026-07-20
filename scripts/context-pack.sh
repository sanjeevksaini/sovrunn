#!/usr/bin/env bash
set -euo pipefail

OUT="${1:-docs/context/SOVRUNN_CONTEXT_PACK.generated.md}"

mkdir -p "$(dirname "$OUT")"

{
  echo "# Sovrunn Generated Context Pack"
  echo
  echo "> Generated from approved repo documents. Do not edit this generated file directly."
  echo

  echo "## Architecture Version"
  cat docs/context/ARCHITECTURE_VERSION.md
  echo

  echo "## Current Architecture Baseline"
  cat docs/context/CURRENT_ARCHITECTURE_BASELINE.md
  echo

  echo "## Context Pack"
  cat docs/context/SOVRUNN_CONTEXT_PACK.md
  echo

  echo "## Current Phase Context"
  cat docs/context/CURRENT_PHASE_CONTEXT.md
  echo

  echo "## Current Decision Summary"
  cat docs/context/CURRENT_DECISION_SUMMARY.md
  echo

  echo "## Constitution"
  sed -n '1,220p' docs/foundation/constitution.md
  echo

  echo "## Development Phases"
  sed -n '1,260p' docs/architecture/development-phases.md
  echo

  echo "## Phase 2 Scope"
  cat docs/phase2/PHASE2_SCOPE.md
  echo

  echo "## Phase 2 Feature Sequence"
  cat docs/phase2/PHASE2_FEATURE_SEQUENCE.md
  echo

  echo "## Feature Roadmap"
  sed -n '1,260p' docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md
  echo

  echo "## Decision Index"
  cat docs/decisions/DECISION_INDEX.md
  echo

  echo "## Glossary"
  sed -n '1,320p' docs/glossary.md
  echo
} > "$OUT"

echo "Generated $OUT"
