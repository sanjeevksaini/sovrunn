---
doc_type: feature
id: FEATURE-0007
title: Plugin and Capability Registry
status: draft
phase: 1
depends_on: [FEATURE-0006]
ai_load_priority: feature
ai_summary: Implements plugin registry and declared lifecycle capabilities.
---

# FEATURE-0007 Plugin and Capability Registry

## 1. Objective

Implement `Plugin` and `Capability` registry primitives.

This does not execute plugins yet. It only records what a plugin claims to support.

## 2. Plugin Resource

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "Plugin",
  "metadata": {
    "name": "postgres.dstoreops.basic"
  },
  "spec": {
    "pluginType": "dStoreOps",
    "version": "0.1.0",
    "serviceClassRefs": ["datastore.postgresql"],
    "deploymentMode": "compiled-in",
    "description": "Basic PostgreSQL datastore operations plugin"
  },
  "status": {
    "phase": "Active"
  }
}
```

## 3. Capability Resource

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "Capability",
  "metadata": {
    "name": "postgres-basic-provision"
  },
  "spec": {
    "pluginRef": "postgres.dstoreops.basic",
    "serviceClassRef": "datastore.postgresql",
    "operation": "Provision",
    "supported": true
  },
  "status": {
    "phase": "Active"
  }
}
```

## 4. Allowed Plugin Types

```text
dStoreOps
cacheOps
streamOps
objectOps
gatewayOps
faasOps
lbOps
k8sOps
bigDataOps
sdeOps
```

## 5. Allowed Capability Operations

```text
Validate
Plan
Provision
Configure
Bind
Observe
Scale
Upgrade
Backup
Restore
RotateCredentials
Unbind
Delete
```

## 6. API

| Method | Path |
|---|---|
| POST | `/v1/plugins` |
| GET | `/v1/plugins/{name}` |
| GET | `/v1/plugins` |
| PUT | `/v1/plugins/{name}` |
| DELETE | `/v1/plugins/{name}` |
| POST | `/v1/capabilities` |
| GET | `/v1/capabilities/{name}` |
| GET | `/v1/capabilities?pluginRef=postgres.dstoreops.basic` |
| GET | `/v1/capabilities?serviceClassRef=datastore.postgresql` |
| DELETE | `/v1/capabilities/{name}` |

## 7. Validation

- Plugin serviceClassRefs must exist.
- Capability pluginRef must exist.
- Capability serviceClassRef must exist.
- Capability operation must be allowed.
- Delete Plugin fails if Capability entries still reference it.

## 8. Acceptance Criteria

- Register PostgreSQL dStoreOps plugin.
- Register Provision/Bind/Observe/Delete capabilities.
- Reject capability for missing plugin.
- Reject capability for missing ServiceClass.
- Can query capabilities by plugin and service class.
