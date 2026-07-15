# Sovrunn Product Steering

## Product Identity

Sovrunn is a sovereign cloud-native PaaS platform.

It provides an opinionated product layer above proven open-source and open-standard infrastructure.

## SDE Positioning

SDE is a major differentiated capability inside Sovrunn. SDE is not the entire platform.

## Target Customers

```text
government-scale platforms
local cloud providers
local colocation providers
regulated enterprises
on-prem/cloud enterprises
system integrators
```

## Product Goals

Sovrunn should provide:

```text
organization-first governance
multi-tenant service consumption
service catalog
service plans
Service Management Plane registry
ServiceOps plugin framework
capability registry
operation framework
policy inheritance
audit aggregation
backup and archival governance
cloud management across multiple sovereign datacenter locations
AI-assisted operations
SDE
```

## Phase 1 Goal

Phase 1 builds the platform grammar:

```text
Organization
OrganizationUnit
Tenant
Project
Operation
ServiceClass
ServicePlan
Plugin
Capability
ServiceInstance
ServiceBinding
```

## Phase 1 Non-Goals

Do not build yet:

```text
production UI
billing engine
marketplace
multi-cluster federation implementation
persistent database storage
Kubernetes CRDs
GitOps controller
ServiceOps plugin execution
real datastore provisioning
AI agent execution
SDE transformation
```

## Product Rule

Expose simple, governed service consumption. Do not expose raw Kubernetes complexity to tenants.
