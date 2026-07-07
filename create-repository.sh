#!/usr/bin/env bash

set -e

ROOT="sovrunn"

mkdir -p "$ROOT"

cd "$ROOT"

# -------------------------------------------------------------------
# Root
# -------------------------------------------------------------------

touch README.md
touch LICENSE
touch .gitignore

# -------------------------------------------------------------------
# Documentation
# -------------------------------------------------------------------

mkdir -p docs/{foundation,specifications,architecture,implementation,rfc}

# -------------------------------------------------------------------
# Foundation
# -------------------------------------------------------------------

touch docs/foundation/vision.md
touch docs/foundation/constitution.md
touch docs/foundation/ontology.md
touch docs/foundation/ads.md
touch docs/foundation/style.md
touch docs/foundation/knowledge.md
touch docs/foundation/ownership.md
touch docs/foundation/dependencies.md

# -------------------------------------------------------------------
# SIR Specification
# -------------------------------------------------------------------

mkdir -p docs/specifications/sir

touch docs/specifications/sir/sir.md
touch docs/specifications/sir/adopted-standards.md
touch docs/specifications/sir/concepts.md
touch docs/specifications/sir/resources.md
touch docs/specifications/sir/relationships.md
touch docs/specifications/sir/operations.md
touch docs/specifications/sir/expressions.md
touch docs/specifications/sir/constraints.md
touch docs/specifications/sir/capabilities.md
touch docs/specifications/sir/conformance.md

# -------------------------------------------------------------------
# Protocol Specification
# -------------------------------------------------------------------

mkdir -p docs/specifications/protocol

touch docs/specifications/protocol/protocol.md
touch docs/specifications/protocol/native.md
touch docs/specifications/protocol/postgresql.md
touch docs/specifications/protocol/mysql.md
touch docs/specifications/protocol/mongodb.md
touch docs/specifications/protocol/redis.md
touch docs/specifications/protocol/grpc.md
touch docs/specifications/protocol/rest.md

# -------------------------------------------------------------------
# Engine Specification
# -------------------------------------------------------------------

mkdir -p docs/specifications/engine

touch docs/specifications/engine/engine.md
touch docs/specifications/engine/postgresql.md
touch docs/specifications/engine/mysql.md
touch docs/specifications/engine/cassandra.md
touch docs/specifications/engine/mongodb.md
touch docs/specifications/engine/redis.md
touch docs/specifications/engine/opensearch.md
touch docs/specifications/engine/neo4j.md
touch docs/specifications/engine/milvus.md
touch docs/specifications/engine/iceberg.md
touch docs/specifications/engine/delta-lake.md
touch docs/specifications/engine/parquet.md
touch docs/specifications/engine/s3.md

# -------------------------------------------------------------------
# Capability Specification
# -------------------------------------------------------------------

mkdir -p docs/specifications/capability

touch docs/specifications/capability/capability.md
touch docs/specifications/capability/transactions.md
touch docs/specifications/capability/indexing.md
touch docs/specifications/capability/search.md
touch docs/specifications/capability/graph.md
touch docs/specifications/capability/vector.md
touch docs/specifications/capability/object.md
touch docs/specifications/capability/streaming.md
touch docs/specifications/capability/cache.md
touch docs/specifications/capability/federation.md
touch docs/specifications/capability/security.md

# -------------------------------------------------------------------
# Platform Specifications
# -------------------------------------------------------------------

mkdir -p docs/specifications/serialization
mkdir -p docs/specifications/versioning
mkdir -p docs/specifications/reuse

touch docs/specifications/serialization/serialization.md
touch docs/specifications/versioning/versioning.md
touch docs/specifications/reuse/reuse.md

# -------------------------------------------------------------------
# Architecture
# -------------------------------------------------------------------

mkdir -p docs/architecture

touch docs/architecture/architecture.md
touch docs/architecture/runtime.md
touch docs/architecture/protocol-runtime.md
touch docs/architecture/semantic-translation.md
touch docs/architecture/planning.md
touch docs/architecture/execution-plan.md
touch docs/architecture/data-kernel.md
touch docs/architecture/engine-runtime.md
touch docs/architecture/federation.md
touch docs/architecture/metadata-runtime.md
touch docs/architecture/plugin-framework.md
touch docs/architecture/lifecycle.md

# -------------------------------------------------------------------
# Implementation
# -------------------------------------------------------------------

mkdir -p docs/implementation

touch docs/implementation/skeleton.md
touch docs/implementation/repository.md
touch docs/implementation/modules.md
touch docs/implementation/build.md
touch docs/implementation/testing.md
touch docs/implementation/coding.md
touch docs/implementation/naming.md

# -------------------------------------------------------------------
# RFC
# -------------------------------------------------------------------

touch docs/rfc/README.md

# -------------------------------------------------------------------
# Source Code (future)
# -------------------------------------------------------------------

mkdir -p cmd
mkdir -p internal
mkdir -p pkg
mkdir -p plugins
mkdir -p sdk
mkdir -p api
mkdir -p examples
mkdir -p test
mkdir -p scripts
mkdir -p configs
mkdir -p deployments

echo
echo "✅ Sovrunn repository structure created."
