# Naming Standards

Document:
  ID: implementation-naming
  Title: Naming Standards
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: CONTRACT
  Version: 1.1
  Status: Draft

Purpose:
  - Define implementation naming standards for Sovrunn Data Engine
  - Align code names with glossary terms
  - Clarify Datastore Management Plane as a pluggable management plane
  - Prevent ambiguous package, type, API, plugin, and resource names
  - Help AI coding agents generate consistent names

Core Naming Rule:
  Names in code should reflect canonical glossary terms.

  Do not invent synonyms for architecture terms.

Product Names:
  Canonical:
    - Sovrunn
    - Sovrunn Data Engine
    - SDE

  Go Package Prefix:
    - sde

  Binary Names:
    - sde-control-plane
    - sde-data-plane
    - sde-management-plane-controller
    - sde-dmp-controller
    - sde-plugin-runner
    - sde-cli

Plane Names:
  Use:
    - controlplane
    - dataplane
    - runtime
    - managementplane
    - dmp

  Avoid:
    - cp
    - dp
    - databaseplatform
    - engineplane
    - fixed dmp subsystem wording

Management Plane Names:
  Use:
    - Management Plane Framework
    - Pluggable Management Plane
    - Datastore Management Plane
    - DMP Controller Runtime
    - sde-dmp-controller

  Meaning:
    - DMP means the pluggable management plane.
    - DMP Controller Runtime means the executable runtime hosting and reconciling DMP.
    - sde-dmp-controller is the binary name for the DMP Controller Runtime.

  Avoid:
    - Management Plane Plugin as a generic synonym
    - Calling sde-dmp-controller the whole DMP
    - Treating DMP as a fixed hard-coded subsystem

Datastore Names:
  Use:
    - downstreamdatastore
    - downstreamengine
    - datastoreinstance
    - datastorerequest
    - datastoreprofile
    - datastorepolicy

  Avoid:
    - database
    - dbplatform
    - enginedataplane

Plugin Names:
  Use:
    - protocolplugin
    - engineplugin
    - datastoreoperatorplugin
    - infrastructureprovider
    - foundationprovider

  Avoid:
    - managementplugin
    - provider without qualifier
    - databaseplugin
    - engineprovider for infrastructure

Package Naming:
  Rules:
    - Use lowercase Go package names.
    - Avoid underscores.
    - Avoid hyphens in Go package names.
    - Use short but precise names.
    - Avoid generic names like common, util, helper unless truly generic.

Examples:
  Good:
    - protocol
    - engine
    - managementplane
    - dmp
    - datastoreoperator
    - infrastructureprovider
    - foundationprovider
    - capability
    - workflow
    - namespace
    - result
    - errormodel

  Bad:
    - plugins2
    - db
    - helpers
    - common
    - provider
    - manager
    - engineplane

Directory Naming:
  Rules:
    - Directory names may use hyphen for external plugin names.
    - Go package directories should avoid hyphen if imported as packages.
    - Use canonical domain names.

Examples:
  - plugins/protocol/postgresql
  - plugins/engine/delta-lake
  - plugins/management-plane/datastore-management-plane
  - plugins/datastore-operator/postgresql
  - plugins/infrastructure-provider/kubernetes
  - plugins/foundation-provider/secrets

Type Naming:
  Good:
    - ManagementPlane
    - ManagementPlaneManifest
    - ManagementPlaneRegistry
    - ManagementPlaneController
    - DatastoreManagementPlane
    - DMPControllerRuntime
    - ProtocolPlugin
    - EnginePlugin
    - DatastoreOperatorPlugin
    - InfrastructureProvider
    - FoundationProvider
    - DatastoreRequestController
    - TenantNamespaceManager
    - ExecutionPlan
    - ExecutionContext
    - ResultModel
    - ErrorModel

  Bad:
    - DMPManager when meaning DMP itself
    - DBManager
    - PluginManager without domain
    - Provider without qualifier
    - EngineDataPlane
    - DatabaseEngine

Resource Naming:
  Kubernetes-style resources should use canonical SDE names.

  Examples:
    - ManagementPlane
    - ManagementPlaneManifest
    - TenantNamespace
    - DatastoreRequest
    - DatastoreInstance
    - DatastoreProfile
    - DatastorePolicy
    - BackupPolicy
    - ScalingPolicy
    - MaintenancePolicy
    - MonitoringPolicy
    - CredentialPolicy
    - ConnectionProfile
    - DatastoreOperation
    - DatastoreWorkflow

API Naming:
  REST paths:
    - /v1/management-planes
    - /v1/management-planes/datastore-management-plane
    - /v1/tenants/{tenantId}/namespaces
    - /v1/tenants/{tenantId}/datastore-requests
    - /v1/tenants/{tenantId}/datastore-instances
    - /v1/plugins/protocol
    - /v1/plugins/engine
    - /v1/dmp/workflows

Event Naming:
  Use domain.event.version pattern.

  Examples:
    - sde.managementplane.registered.v1
    - sde.dmp.started.v1
    - sde.datastore.request.created.v1
    - sde.datastore.instance.ready.v1
    - sde.plugin.admitted.v1
    - sde.workflow.completed.v1
    - sde.runtime.error.observed.v1

AI Naming:
  Current Scope:
    - AI Control Plane is reserved.
    - Tenant AI Agent is reserved.

  Use:
    - aicontrolplane
    - tenantaiagent
    - aiactionpolicy
    - aiworkflowrequest
    - airecommendation

Deprecated Names:
  Do Not Use:
    - Database Platform
    - Downstream Database
    - Engine Data Plane
    - Management Plane Plugin
    - Provider without qualifier
    - Non-Destructible Auto-Tuning

Replacement:
  - Database Platform → Sovrunn Data Engine
  - Downstream Database → Downstream Datastore
  - Engine Data Plane → Datastore Data Plane
  - Management Plane Plugin → Pluggable Management Plane, Datastore Management Plane, Datastore Operator Plugin, Infrastructure Provider, or Foundation Provider
  - Provider → Foundation Provider or Infrastructure Provider
  - Non-Destructible Auto-Tuning → Non-Destructive Auto-Tuning or AI Safe Autotuning

Invariants:
  - Naming must reflect architecture boundaries.
  - Naming must not blur Control Plane and Data Plane responsibilities.
  - Naming must not blur DMP and DMP Controller Runtime.
  - Naming must not blur Engine Plugin and Datastore Operator Plugin.
  - Naming must not blur Infrastructure Provider and Foundation Provider.
  - New architecture terms require glossary update.
