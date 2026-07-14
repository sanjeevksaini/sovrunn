---
doc_type: architecture
title: GitOps Desired-State Model
status: draft
phase: 1
ai_load_priority: phase
ai_summary: Defines how Sovrunn Phase 1 resources should be shaped so they can later be managed through GitOps.
---

# GitOps Desired-State Model

## 1. Purpose

Sovrunn must support GitOps as a first-class access channel in later phases.

Phase 1 does not implement Argo CD, Flux, or a Sovrunn GitOps controller. However, all resources must be designed so they can be safely represented as declarative YAML/JSON and stored in Git.

The principle is:

```text
Git stores desired state.
Sovrunn reconciles actual state.
Operations record what changed.
```

## 2. GitOps Channel Meaning

GitOps is an access channel alongside:

```text
Portal
CLI
API
SDKs
AI Assistant
```

Example flow:

```text
developer commits ServiceInstance YAML
  -> pull request review
  -> GitOps sync applies resource
  -> Sovrunn validates policy and references
  -> Sovrunn records Operation
  -> Sovrunn updates status
```

## 3. Resource Shape

All resources must be representable as YAML/JSON:

```yaml
apiVersion: platform.sovrunn.io/v1alpha1
kind: ServiceInstance
metadata:
  name: nhm-prod-postgres
  labels:
    environment: production
spec:
  organizationRef: nic
  organizationUnitRef: ministry-health
  tenantRef: national-health-mission
  projectRef: production
  serviceClassRef: datastore.postgresql
  servicePlanRef: postgres-small-ha
  parameters:
    databaseName: nhm
```

Users author:

```text
apiVersion
kind
metadata.name
metadata.labels
metadata.annotations
spec
```

Users do not author:

```text
status
operation records
resourceVersion
observed state
runtime-generated secrets
```

## 4. Stable Names

GitOps requires stable names.

Good:

```text
national-health-mission
production
nhm-prod-postgres
postgres-small-ha
```

Bad:

```text
service-abc123
temp-db-final-latest
prod_new_2
```

## 5. Explicit References

References must be explicit:

```yaml
organizationRef: nic
organizationUnitRef: ministry-health
tenantRef: national-health-mission
projectRef: production
serviceClassRef: datastore.postgresql
servicePlanRef: postgres-small-ha
```

Avoid hidden inference for core relationships in Phase 1.

## 6. Status Separation

GitOps systems should ignore or avoid committing status.

If status is provided in a create/update request, recommended Phase 1 behavior is:

```text
reject user-authored status for mutating API requests
```

## 7. Idempotent Apply

Future GitOps apply should be idempotent.

Applying the same manifest multiple times should produce the same desired state.

Phase 1 REST create/update may be separate, but resource design should prepare for:

```text
server-side apply
upsert
patch
drift detection
reconciliation
```

## 8. Directory Layout Recommendation

Future GitOps repos may use:

```text
sovrunn-config/
  organizations/
    nic.yaml
  organization-units/
    ministry-health.yaml
  tenants/
    national-health-mission.yaml
  projects/
    production.yaml
  catalog/
    service-classes/
    service-plans/
  plugins/
    postgres.dstoreops.basic.yaml
  service-instances/
    nhm-prod-postgres.yaml
  service-bindings/
    nhm-app-postgres-binding.yaml
```

## 9. Promotion Model

Environment promotion should happen through Git branches or folders.

Use Project for workload/environment grouping.

Do not overload Tenant for environment promotion.

## 10. Secrets

Do not store raw secrets in Git.

Use references:

```yaml
secretRef: existing-secret-ref
keyRef: existing-key-ref
externalSecretRef: vault/path/name
```

## 11. AI and GitOps

AI Assistant should be able to draft GitOps manifests, but not bypass review.

Recommended flow:

```text
AI drafts manifest
  -> human reviews pull request
  -> GitOps sync applies
  -> Sovrunn validates
  -> Operation records result
```

## 12. Phase 1 Non-Goals

Do not implement yet:

```text
Argo CD integration
Flux integration
Git webhook receiver
Sovrunn GitOps controller
server-side apply
three-way merge
drift detection
Git-based approval engine
```

## 13. Acceptance Criteria

Phase 1 is GitOps-ready if:

```text
all resources serialize cleanly as YAML/JSON
spec is declarative
status is system-owned
references are explicit
names are stable
demo manifests can be stored in Git
resources can be recreated from manifest examples
```

## 14. Final Principle

Sovrunn should be API-first and GitOps-friendly from the beginning, even before GitOps automation exists.
