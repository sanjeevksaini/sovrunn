# Sovrunn Data Engine

Sovrunn Data Engine, abbreviated as SDE, is an AI-native sovereign datastore platform architecture.

This repository contains the architecture, specifications, implementation guidance, and RFC framework for SDE.

## Documentation

The main documentation is under:

```text
docs/
```

The MkDocs configuration is:

```text
mkdocs.yml
```

To build the documentation:

```bash
mkdocs build --strict
```

To serve the documentation locally:

```bash
mkdocs serve
```

Then open:

```text
http://127.0.0.1:8000
```

## Documentation Structure

```text
docs/
  foundation/
  architecture/
  specifications/
  implementation/
  rfc/
```

## Current Architecture Focus

The current architecture defines:

```text
Sovrunn Data Engine
  ├── SDE Control Plane
  ├── SDE Data Plane
  ├── SDE Runtime
  ├── Specifications
  ├── Implementation Layer
  └── RFC Framework
```

The SDE Control Plane includes a Management Plane Framework. Datastore Management Plane is modeled as the first pluggable management plane inside the SDE Control Plane.

## Key Terms

```text
SDE
  Sovrunn Data Engine

SDE Control Plane
  Management authority and governance plane for SDE

SDE Data Plane
  Runtime request execution plane

Management Plane Framework
  Framework for governed, pluggable management planes

Datastore Management Plane
  First pluggable management plane for downstream datastore lifecycle and operations

DMP Controller Runtime
  Executable runtime that hosts and reconciles DMP resources and workflows

dstoreOps
  Managed datastore operations capability powered by DMP
```

## Repository Notes

This repository is documentation-first.

Architecture and specifications are the source of truth. RFCs record decisions and changes. Implementation documents translate accepted architecture into code-ready structure.

## Local Validation

Recommended validation commands:

```bash
find . -name ".DS_Store" -delete
mkdocs build --strict
```

Optional Python virtual environment:

```bash
python3 -m venv .venv
source .venv/bin/activate
python -m pip install --upgrade pip
python -m pip install mkdocs mkdocs-material pymdown-extensions
mkdocs build --strict
```
