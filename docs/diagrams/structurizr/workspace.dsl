workspace "Sovrunn" "Sovrunn sovereign PaaS architecture model" {

    !identifiers hierarchical

    model {
        providerOperator = person "Provider Operator" "Operates local cloud, MSP, or enterprise PaaS services."
        customerAdmin = person "Customer Admin" "Requests and manages governed PaaS services."
        platformArchitect = person "Sovrunn Architect" "Maintains approved architecture baseline, DEC/RFC records, and C4 model."

        sovrunn = softwareSystem "Sovrunn Platform" "Cloud-native sovereign PaaS platform for governed service catalog, provider-neutral placement, lifecycle orchestration, audit, evidence, and AI-assisted operations." {
            api = container "Sovrunn API Server" "Exposes customer/provider APIs for organizations, tenants, projects, service catalog, placement, operations, and audit." "Go"
            controlPlane = container "Sovrunn Control Plane" "Coordinates governance, policy context, placement decisions, operations, plugin execution, audit, and AI-readable explanations." "Go"
            policyAdapter = container "Policy Evaluation Adapter" "Abstracts policy evaluation through OPA/Cedar-compatible adapters." "Go interface"
            placementEngine = container "Placement Decision Engine" "Matches service runtime requirements against provider capabilities and effective policy context." "Go"
            operationController = container "Operation Controller" "Tracks long-running operations and plugin execution steps." "Go"
            registry = container "Registry" "Stores Sovrunn resources, service catalog, profiles, decisions, operations, and audit records." "PostgreSQL/YugabyteDB later"
            pluginBoundary = container "Plugin Execution Boundary" "Executes provider, service management, and service runtime plugins through governed contracts." "Go plugin boundary"
            aiDecisionContext = container "AI Decision Context" "Produces AI-readable decision explanations and operation context." "Structured JSON"
        }

        providerPlugin = softwareSystem "Provider/Substrate Plugin" "Executes infrastructure operations against Kubernetes, VM, bare metal, or provider APIs."
        pgManagementPlugin = softwareSystem "PostgreSQL Management Plugin" "Plans PostgreSQL service lifecycle and runtime requirements."
        pgRuntimePlugin = softwareSystem "PostgreSQL Runtime Plugin" "Creates and manages PostgreSQL runtime resources through reused operators or Helm."

        kubernetes = softwareSystem "Kubernetes / k3s Substrate" "Initial local substrate for Phase 3 MVP."
        postgresOperator = softwareSystem "PostgreSQL Operator / Helm" "Reused PostgreSQL runtime foundation such as CloudNativePG, Crunchy, or Helm."
        policyEngines = softwareSystem "OPA / Cedar Policy Engines" "Reusable policy engines integrated through PolicyEngineAdapter in later phases."
        observability = softwareSystem "OpenTelemetry / Prometheus / Grafana" "Reusable observability stack for logs, metrics, traces, and operational views."
        secretProvider = softwareSystem "Vault / External Secrets" "Reusable secrets and credential management integrated through SecretProviderAdapter later."
        gitRepo = softwareSystem "Sovrunn Git Repository" "Source of truth for Architecture Operating System docs, DEC/RFC records, Kiro specs, and Structurizr DSL."
        kiro = softwareSystem "Kiro" "Architecture/spec studio that applies approved Architecture Decision Handoffs to repo docs and feature specs."
        chatgpt = softwareSystem "ChatGPT Project" "Architecture tradeoff discussion and Architecture Decision Handoff generation."
        cursor = softwareSystem "Cursor" "Go implementation studio that implements approved Kiro tasks."

        customerAdmin -> sovrunn.api "Requests governed PaaS services"
        providerOperator -> sovrunn.api "Configures providers, capabilities, plugins, and service catalog"
        platformArchitect -> chatgpt "Discusses architecture tradeoffs"
        chatgpt -> gitRepo "Reads Architecture Operating System context and produces handoff"
        platformArchitect -> kiro "Approves handoff for repo update"
        kiro -> gitRepo "Updates architecture docs, DEC/RFCs, feature specs, and Structurizr DSL"
        cursor -> gitRepo "Implements approved tasks and tests"

        sovrunn.api -> sovrunn.controlPlane "Submits service requests and management actions"
        sovrunn.controlPlane -> sovrunn.policyAdapter "Evaluates effective governance, security, data, and cost policies"
        sovrunn.policyAdapter -> policyEngines "Delegates policy evaluation through adapter" "" "Future"
        sovrunn.controlPlane -> sovrunn.placementEngine "Requests placement decisions"
        sovrunn.placementEngine -> sovrunn.registry "Reads ResourcePool and ProviderCapability data"
        sovrunn.controlPlane -> sovrunn.operationController "Creates and tracks operations"
        sovrunn.operationController -> sovrunn.pluginBoundary "Invokes plugin execution"
        sovrunn.controlPlane -> sovrunn.registry "Reads/writes resources, decisions, operations, and audit events"
        sovrunn.controlPlane -> sovrunn.aiDecisionContext "Creates AI-readable decision explanations"
        sovrunn.controlPlane -> observability "Emits metrics, traces, logs, and health signals"
        sovrunn.pluginBoundary -> providerPlugin "Calls provider/substrate operations"
        sovrunn.pluginBoundary -> pgManagementPlugin "Calls service management planning"
        sovrunn.pluginBoundary -> pgRuntimePlugin "Calls runtime lifecycle operations"
        sovrunn.pluginBoundary -> secretProvider "Resolves credential references" "" "Future"
        providerPlugin -> kubernetes "Provisions/validates substrate resources"
        pgRuntimePlugin -> postgresOperator "Creates PostgreSQL runtime using reused operator or Helm"
    }

    views {
        systemContext sovrunn "SystemContext" {
            include *
            autolayout lr
            description "Sovrunn system context showing users, AI/spec/coding workflow, plugin ecosystem, and reused foundations."
        }

        container sovrunn "Containers" {
            include *
            autolayout lr
            description "Sovrunn core containers and their relationships to plugins and reused OSS foundations."
        }

        dynamic sovrunn "ArchitectureHandoffWorkflow" "ChatGPT to Kiro to Cursor governed architecture handoff workflow" {
            platformArchitect -> chatgpt "Discuss architecture tradeoff"
            chatgpt -> gitRepo "Load baseline, DEC index, phase scope, roadmap"
            chatgpt -> platformArchitect "Produce Architecture Decision Handoff"
            platformArchitect -> kiro "Approve handoff for application"
            kiro -> gitRepo "Validate and update architecture docs, specs, traceability, Structurizr DSL"
            platformArchitect -> cursor "Start implementation only from approved Kiro tasks"
            cursor -> gitRepo "Update Go code and tests"
        }

        dynamic sovrunn "GovernedPostgreSQLProvisioning" "Governed PostgreSQL PaaS provisioning flow" {
            customerAdmin -> sovrunn.api "Request PostgreSQL ServiceInstance"
            sovrunn.api -> sovrunn.controlPlane "Submit request"
            sovrunn.controlPlane -> sovrunn.policyAdapter "Evaluate policy"
            sovrunn.controlPlane -> sovrunn.placementEngine "Create PlacementDecision"
            sovrunn.controlPlane -> sovrunn.operationController "Create Operation"
            sovrunn.operationController -> sovrunn.pluginBoundary "Execute plugin chain"
            sovrunn.pluginBoundary -> pgManagementPlugin "Plan PostgreSQL lifecycle"
            sovrunn.pluginBoundary -> providerPlugin "Prepare substrate"
            sovrunn.pluginBoundary -> pgRuntimePlugin "Create runtime"
            pgRuntimePlugin -> postgresOperator "Apply operator/Helm resource"
            sovrunn.controlPlane -> sovrunn.registry "Record status, decision, audit"
            sovrunn.controlPlane -> sovrunn.aiDecisionContext "Generate explanation"
        }

        styles {
            element "Person" {
                shape person
                background #08427b
                color #ffffff
            }
            element "Software System" {
                background #1168bd
                color #ffffff
            }
            element "Container" {
                background #438dd5
                color #ffffff
            }
            relationship "Future" {
                dashed true
            }
        }
    }
}
