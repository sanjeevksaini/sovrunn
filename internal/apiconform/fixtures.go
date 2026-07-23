package apiconform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// ConformanceFixturesDir is the repository-relative directory for FEATURE-0012
// canonical contract fixtures (tasks 14.1–14.4; F12-FIXTURE-001/002).
const ConformanceFixturesDir = "tests/conformance/fixtures"

// ConformanceNegativeFixturesDir holds Task 14.3 invalid fixtures used by
// Matrix D proofs that require negative evidence (D-17 mismatch, etc.).
const ConformanceNegativeFixturesDir = "tests/conformance/fixtures/negative"

// MatrixDScenarioCount is the fixed number of Architecture Matrix D scenarios
// (F12-FIXTURE-001). The conformance suite fails if any becomes unrepresentable.
const MatrixDScenarioCount = 17

// FixtureLoader reads JSON/YAML conformance fixtures from
// tests/conformance/fixtures/ (task 14.4).
type FixtureLoader struct {
	dir string // absolute path to ConformanceFixturesDir
}

// NewFixtureLoader constructs a loader rooted at moduleRoot/tests/conformance/fixtures.
// The fixtures directory must exist; individual fixtures are checked per scenario.
func NewFixtureLoader(moduleRoot string) (*FixtureLoader, error) {
	if strings.TrimSpace(moduleRoot) == "" {
		return nil, fmt.Errorf("fixture loader: module root is empty")
	}
	dir := filepath.Join(moduleRoot, ConformanceFixturesDir)
	st, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("fixture loader: fixtures dir %s: %w", dir, err)
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("fixture loader: %s is not a directory", dir)
	}
	return &FixtureLoader{dir: dir}, nil
}

// Dir returns the absolute fixtures directory path.
func (l *FixtureLoader) Dir() string {
	if l == nil {
		return ""
	}
	return l.dir
}

// Path joins a fixture-relative name (e.g. "project.json" or
// "negative/unknown-field.json") under the fixtures root.
func (l *FixtureLoader) Path(name string) string {
	return filepath.Join(l.dir, filepath.FromSlash(name))
}

// Read returns raw fixture bytes. name is relative to the fixtures directory.
func (l *FixtureLoader) Read(name string) ([]byte, error) {
	if l == nil {
		return nil, fmt.Errorf("fixture loader: nil loader")
	}
	raw, err := os.ReadFile(l.Path(name))
	if err != nil {
		return nil, fmt.Errorf("read fixture %s: %w", name, err)
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, fmt.Errorf("fixture %s is empty", name)
	}
	return raw, nil
}

// Exists reports whether the named fixture file is present and non-empty.
func (l *FixtureLoader) Exists(name string) bool {
	raw, err := l.Read(name)
	return err == nil && len(bytes.TrimSpace(raw)) > 0
}

// DecodeJSON decodes a JSON fixture under ModeReadRepresentation into dst.
func (l *FixtureLoader) DecodeJSON(name string, dst any) error {
	raw, err := l.Read(name)
	if err != nil {
		return err
	}
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)
	if prob := apivalid.DecodeJSON(raw, lim, pol, dst); prob != nil {
		return fmt.Errorf("DecodeJSON %s: code=%s detail=%s", name, prob.Code, prob.Detail)
	}
	return nil
}

// DecodeYAML decodes a strict JSON-compatible YAML fixture under
// ModeReadRepresentation into dst.
func (l *FixtureLoader) DecodeYAML(name string, dst any) error {
	raw, err := l.Read(name)
	if err != nil {
		return err
	}
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)
	if prob := apivalid.DecodeYAML(raw, lim, pol, dst); prob != nil {
		return fmt.Errorf("DecodeYAML %s: code=%s detail=%s", name, prob.Code, prob.Detail)
	}
	return nil
}

// RequireFixtures fails when any listed fixture is missing or empty so a
// Matrix D scenario cannot silently become unrepresentable.
func (l *FixtureLoader) RequireFixtures(names ...string) error {
	var missing []string
	for _, name := range names {
		if !l.Exists(name) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("unrepresentable: missing fixtures: %s", strings.Join(missing, ", "))
	}
	return nil
}

// MatrixDScenario maps one Architecture Matrix D row to fixture paths and a
// required-proof assertion (task 14.4; F12-FIXTURE-001).
type MatrixDScenario struct {
	Name          string
	Profile       apimeta.Profile
	Boundary      apimeta.Boundary
	RequiredProof string
	Fixtures      []string
	Prove         func(*FixtureLoader) error
}

// MatrixDScenarios returns the fixed seventeen Matrix D scenario rows in
// requirements order. Each entry binds fixtures plus a required-proof
// assertion; Prove fails when the scenario is unrepresentable.
func MatrixDScenarios() []MatrixDScenario {
	return []MatrixDScenario{
		{
			Name:          "Customer creates Project",
			Profile:       apimeta.ProfileManagedResource,
			Boundary:      apimeta.BoundaryCustomerFacing,
			RequiredProof: "Stable identity, scope, strict validation, spec/status.",
			Fixtures:      []string{"project.json", "project.yaml"},
			Prove:         proveCustomerCreatesProject,
		},
		{
			Name:          "Operator registers ResourcePool",
			Profile:       apimeta.ProfileManagedResource,
			Boundary:      apimeta.BoundaryOperatorFacing,
			RequiredProof: "Provider-neutral capabilities; no vendor fields in core.",
			Fixtures:      []string{"resource-pool.json", "resource-pool.yaml"},
			Prove:         proveOperatorRegistersResourcePool,
		},
		{
			Name:          "Adapter discovers external database",
			Profile:       apimeta.ProfileObservedExternalResource,
			Boundary:      apimeta.BoundaryAdapterFacing,
			RequiredProof: "Provenance, freshness, stale/deleted semantics.",
			Fixtures:      []string{"discovered-database.json", "discovered-database.yaml"},
			Prove:         proveAdapterDiscoversExternalDatabase,
		},
		{
			Name:          "Publisher releases plugin contract",
			Profile:       apimeta.ProfileVersionedDefinition,
			Boundary:      apimeta.BoundaryPluginFacing,
			RequiredProof: "Immutable published version and compatibility metadata.",
			Fixtures:      []string{"plugin-definition.json", "plugin-definition.yaml"},
			Prove:         provePublisherReleasesPluginContract,
		},
		{
			Name:          "Operator installs plugin",
			Profile:       apimeta.ProfileManagedResource,
			Boundary:      apimeta.BoundaryOperatorFacing,
			RequiredProof: "Definition, installation, and execution remain separate.",
			Fixtures: []string{
				"plugin-definition.json",
				"operation.json",
				"operation-platform.json",
			},
			Prove: proveOperatorInstallsPlugin,
		},
		{
			Name:          "Operator configures adapter",
			Profile:       apimeta.ProfileManagedResource,
			Boundary:      apimeta.BoundaryAdapterFacing,
			RequiredProof: "Secret references and native-config isolation.",
			Fixtures:      []string{"adapter-configuration.json", "adapter-configuration.yaml"},
			Prove:         proveOperatorConfiguresAdapter,
		},
		{
			Name:          "Placement engine evaluates request",
			Profile:       apimeta.ProfileTransientRequestResult,
			Boundary:      apimeta.BoundaryInternalEngineFacing,
			RequiredProof: "Typed request/result without forced persistence.",
			Fixtures:      []string{"placement-evaluation-request.json", "placement-evaluation-request.yaml"},
			Prove:         provePlacementEngineEvaluatesRequest,
		},
		{
			Name:          "Decision becomes auditable",
			Profile:       apimeta.ProfileImmutableRecord,
			Boundary:      apimeta.BoundaryGovernanceOnly,
			RequiredProof: "Immutable subject/actor/input refs for FEATURE-0013.",
			Fixtures:      []string{"audit-event.json", "audit-event.yaml"},
			Prove:         proveDecisionBecomesAuditable,
		},
		{
			Name:          "Future provisioning executes",
			Profile:       apimeta.ProfileLongRunningOperation,
			Boundary:      apimeta.BoundaryPluginFacing,
			RequiredProof: "Idempotency, progress, retry, cancel, terminal result; Operation.scopeRef equals resolved target governance scope (D-17).",
			Fixtures: []string{
				"operation.json",
				"operation-platform.json",
				"operation-organization.json",
				"operation-organizationunit.json",
				"operation-tenant.json",
				"operation-project.json",
				"operation-provider.json",
				"negative/operation-scope-target-kind-mismatch.json",
				"negative/operation-scope-target-uid-mismatch.json",
			},
			Prove: proveFutureProvisioningExecutes,
		},
		{
			Name:          "Portal lists large collections",
			Profile:       apimeta.ProfileListEnvelope,
			Boundary:      apimeta.BoundaryCustomerFacing,
			RequiredProof: "Bounded opaque pagination and deterministic ordering.",
			Fixtures:      []string{"project.json"},
			Prove:         provePortalListsLargeCollections,
		},
		{
			Name:          "Provider disconnects",
			Profile:       apimeta.ProfileObservedExternalResource,
			Boundary:      apimeta.BoundaryAdapterFacing,
			RequiredProof: "Current, stale, unknown, absent remain distinct.",
			Fixtures:      []string{"discovered-database.json"},
			Prove:         proveProviderDisconnects,
		},
		{
			Name:          "External object recreated under same name",
			Profile:       apimeta.ProfileObservedExternalResource,
			Boundary:      apimeta.BoundaryAdapterFacing,
			RequiredProof: "UID prevents stale-reference rebinding.",
			Fixtures:      []string{"discovered-database.json"},
			Prove:         proveExternalObjectRecreatedUnderSameName,
		},
		{
			Name:          "Cross-tenant reference attempted",
			Profile:       apimeta.ProfileManagedResource,
			Boundary:      apimeta.BoundaryCustomerFacing,
			RequiredProof: "Denial without target-existence disclosure.",
			Fixtures:      []string{"project.json"},
			Prove:         proveCrossTenantReferenceAttempted,
		},
		{
			Name:          "New cloud provider added",
			Profile:       apimeta.ProfileManagedResource,
			Boundary:      apimeta.BoundaryAdapterFacing,
			RequiredProof: "No customer/core schema change.",
			Fixtures: []string{
				"adapter-configuration.json",
				"resource-pool.json",
				"project.json",
			},
			Prove: proveNewCloudProviderAdded,
		},
		{
			Name:          "New data-service plugin added",
			Profile:       apimeta.ProfileVersionedDefinition,
			Boundary:      apimeta.BoundaryPluginFacing,
			RequiredProof: "No core grammar redesign.",
			Fixtures: []string{
				"plugin-definition.json",
				"operation.json",
				"operation-provider.json",
			},
			Prove: proveNewDataServicePluginAdded,
		},
		{
			Name:          "AI explains denial",
			Profile:       apimeta.ProfileImmutableRecord,
			Boundary:      apimeta.BoundaryGovernanceOnly,
			RequiredProof: "Stable codes and safe context, no message scraping.",
			Fixtures:      []string{"audit-event.json"},
			Prove:         proveAIExplainsDenial,
		},
		{
			Name:          "Phase 1 resource migrated",
			Profile:       apimeta.ProfileManagedResource,
			Boundary:      apimeta.BoundaryCustomerFacing,
			RequiredProof: "Explicit version/migration; no silent reinterpretation.",
			Fixtures:      []string{"project.json"},
			Prove:         provePhase1ResourceMigrated,
		},
	}
}

// AssertAllMatrixDScenarios runs every Matrix D scenario proof. It fails if
// the table is incomplete or any scenario is unrepresentable.
func AssertAllMatrixDScenarios(l *FixtureLoader) error {
	if l == nil {
		return fmt.Errorf("AssertAllMatrixDScenarios: nil fixture loader")
	}
	scenarios := MatrixDScenarios()
	if len(scenarios) != MatrixDScenarioCount {
		return fmt.Errorf("Matrix D table length = %d, want %d", len(scenarios), MatrixDScenarioCount)
	}
	seen := make(map[string]struct{}, len(scenarios))
	for _, sc := range scenarios {
		if sc.Name == "" {
			return fmt.Errorf("Matrix D scenario has empty name")
		}
		if _, dup := seen[sc.Name]; dup {
			return fmt.Errorf("duplicate Matrix D scenario %q", sc.Name)
		}
		seen[sc.Name] = struct{}{}
		if sc.RequiredProof == "" {
			return fmt.Errorf("scenario %q: missing RequiredProof", sc.Name)
		}
		if sc.Prove == nil {
			return fmt.Errorf("scenario %q: missing Prove assertion", sc.Name)
		}
		if err := l.RequireFixtures(sc.Fixtures...); err != nil {
			return fmt.Errorf("scenario %q: %w", sc.Name, err)
		}
		if err := sc.Prove(l); err != nil {
			return fmt.Errorf("scenario %q: %w", sc.Name, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Required-proof implementations (grammar representation only; no domain runtime)
// ---------------------------------------------------------------------------

func proveCustomerCreatesProject(l *FixtureLoader) error {
	var p Project
	if err := l.DecodeJSON("project.json", &p); err != nil {
		return err
	}
	if p.APIVersion == "" || p.Kind != "Project" {
		return fmt.Errorf("stable type identity missing: apiVersion=%q kind=%q", p.APIVersion, p.Kind)
	}
	if p.Metadata.Name == "" || p.Metadata.UID == "" {
		return fmt.Errorf("stable identity missing: name=%q uid=%q", p.Metadata.Name, p.Metadata.UID)
	}
	if p.Metadata.ScopeRef == nil || apimeta.ScopeKind(p.Metadata.ScopeRef.Kind) != apimeta.ScopeTenant {
		return fmt.Errorf("Tenant scopeRef required, got %+v", p.Metadata.ScopeRef)
	}
	if p.Spec.Description == "" {
		return fmt.Errorf("spec must carry desired-state fields")
	}
	if p.Status.Phase == "" {
		return fmt.Errorf("status.phase must be present under read representation")
	}
	var fromYAML Project
	if err := l.DecodeYAML("project.yaml", &fromYAML); err != nil {
		return err
	}
	if !reflect.DeepEqual(p, fromYAML) {
		return fmt.Errorf("JSON/YAML project fixture diverge")
	}
	return nil
}

func proveOperatorRegistersResourcePool(l *FixtureLoader) error {
	var pool ResourcePool
	if err := l.DecodeJSON("resource-pool.json", &pool); err != nil {
		return err
	}
	if pool.Spec.CapabilityClass == "" {
		return fmt.Errorf("capabilityClass required for provider-neutral pool")
	}
	if err := rejectVendorNativeTokens(pool.Spec.CapabilityClass); err != nil {
		return err
	}
	raw, err := l.Read("resource-pool.json")
	if err != nil {
		return err
	}
	if err := rejectVendorNativeTokens(string(raw)); err != nil {
		return fmt.Errorf("resource-pool fixture: %w", err)
	}
	if pool.Metadata.ScopeRef == nil || apimeta.ScopeKind(pool.Metadata.ScopeRef.Kind) != apimeta.ScopeProvider {
		return fmt.Errorf("Provider scopeRef required, got %+v", pool.Metadata.ScopeRef)
	}
	return nil
}

func proveAdapterDiscoversExternalDatabase(l *FixtureLoader) error {
	var db DiscoveredDatabase
	if err := l.DecodeJSON("discovered-database.json", &db); err != nil {
		return err
	}
	if db.Provenance.SourceRef.Name == "" || db.Provenance.ObservedAt == "" {
		return fmt.Errorf("provenance sourceRef/observedAt required")
	}
	if db.Freshness.FreshnessState == "" {
		return fmt.Errorf("freshnessState required")
	}
	switch db.Status.ObservationState {
	case ObservationStateCurrent, ObservationStateStale, ObservationStateUnknown, ObservationStateAbsent:
		// ok
	default:
		return fmt.Errorf("observationState %q not in Current/Stale/Unknown/Absent", db.Status.ObservationState)
	}
	return nil
}

func provePublisherReleasesPluginContract(l *FixtureLoader) error {
	var def PluginDefinition
	if err := l.DecodeJSON("plugin-definition.json", &def); err != nil {
		return err
	}
	if def.Spec.Version == "" || def.Spec.CompatibilityRange == "" {
		return fmt.Errorf("version and compatibilityRange required")
	}
	if def.Spec.PublicationState != PublicationStatePublished {
		return fmt.Errorf("publicationState=%q, want Published", def.Spec.PublicationState)
	}
	return nil
}

func proveOperatorInstallsPlugin(l *FixtureLoader) error {
	var def PluginDefinition
	if err := l.DecodeJSON("plugin-definition.json", &def); err != nil {
		return err
	}
	raw, err := l.Read("plugin-definition.json")
	if err != nil {
		return err
	}
	lower := strings.ToLower(string(raw))
	for _, banned := range []string{`"install"`, `"installed"`, `"execute"`, `"execution"`, `"runtime"`} {
		if strings.Contains(lower, banned) {
			return fmt.Errorf("plugin definition must not embed installation/execution fields (%s)", banned)
		}
	}
	var op Operation
	if err := l.DecodeJSON("operation.json", &op); err != nil {
		return err
	}
	if def.Kind == op.Kind {
		return fmt.Errorf("definition and operation must remain distinct kinds")
	}
	if op.Spec.Action == "" {
		return fmt.Errorf("operation action required to prove execution is a separate contract")
	}
	return nil
}

func proveOperatorConfiguresAdapter(l *FixtureLoader) error {
	var cfg AdapterConfiguration
	if err := l.DecodeJSON("adapter-configuration.json", &cfg); err != nil {
		return err
	}
	if cfg.Spec.CredentialsSecretRef.Name == "" || cfg.Spec.CredentialsSecretRef.Kind == "" {
		return fmt.Errorf("credentialsSecretRef required (secret reference, not raw secret)")
	}
	if cfg.Spec.NativeConfigRef == nil || cfg.Spec.NativeConfigRef.Name == "" {
		return fmt.Errorf("nativeConfigRef required for native-config isolation")
	}
	raw, err := l.Read("adapter-configuration.json")
	if err != nil {
		return err
	}
	lower := strings.ToLower(string(raw))
	for _, banned := range []string{"password", "secretvalue", "apikey", "private_key", "connectionstring"} {
		if strings.Contains(lower, banned) {
			return fmt.Errorf("adapter configuration must not embed raw secret-like values (%s)", banned)
		}
	}
	return nil
}

func provePlacementEngineEvaluatesRequest(l *FixtureLoader) error {
	var req PlacementEvaluationRequest
	if err := l.DecodeJSON("placement-evaluation-request.json", &req); err != nil {
		return err
	}
	if req.Request.SubjectRef.Name == "" || req.Request.RequiredCapabilityClass == "" {
		return fmt.Errorf("typed request fields required")
	}
	if req.Result == nil || req.Result.Outcome == "" || req.Result.ReasonCode == "" {
		return fmt.Errorf("typed result required without implying persistence")
	}
	// TransientRequestResult: no forced status.phase lifecycle field on the type.
	rv := reflect.TypeOf(req)
	if _, ok := rv.FieldByName("Status"); ok {
		return fmt.Errorf("PlacementEvaluationRequest must not force a Status persistence shape")
	}
	return nil
}

func proveDecisionBecomesAuditable(l *FixtureLoader) error {
	var ev AuditEvent
	if err := l.DecodeJSON("audit-event.json", &ev); err != nil {
		return err
	}
	if ev.Record.ActorRef.Name == "" || ev.Record.SubjectRef.Name == "" {
		return fmt.Errorf("actorRef and subjectRef required")
	}
	if ev.Record.RequestID == "" || ev.Record.ReasonCode == "" {
		return fmt.Errorf("requestId and reasonCode required for auditable decisions")
	}
	if ev.Record.Action == "" || ev.Record.Outcome == "" {
		return fmt.Errorf("action and outcome required")
	}
	return nil
}

func proveFutureProvisioningExecutes(l *FixtureLoader) error {
	variants := []struct {
		file string
		kind apimeta.ScopeKind // expected canonical op scope; Platform uses nil scopeRef
	}{
		{"operation-platform.json", apimeta.ScopePlatform},
		{"operation-organization.json", apimeta.ScopeOrganization},
		{"operation-organizationunit.json", apimeta.ScopeOrganizationUnit},
		{"operation-tenant.json", apimeta.ScopeTenant},
		{"operation-project.json", apimeta.ScopeProject},
		{"operation-provider.json", apimeta.ScopeProvider},
	}
	if len(variants) != 6 {
		return fmt.Errorf("expected six Operation scope variants, got %d", len(variants))
	}

	for _, v := range variants {
		var op Operation
		if err := l.DecodeJSON(v.file, &op); err != nil {
			return err
		}
		if err := assertOperationProvisioningShape(&op); err != nil {
			return fmt.Errorf("%s: %w", v.file, err)
		}
		opScope := apimeta.CanonicalScopeIdentity(op.Metadata.ScopeRef)
		if opScope.Kind != v.kind {
			return fmt.Errorf("%s: CanonicalScopeIdentity.kind=%q, want %q", v.file, opScope.Kind, v.kind)
		}
		if v.kind == apimeta.ScopePlatform {
			if op.Metadata.ScopeRef != nil {
				return fmt.Errorf("%s: Platform must use canonical nil scopeRef", v.file)
			}
			if opScope.UID != apimeta.PlatformScopeUID {
				return fmt.Errorf("%s: Platform UID=%q, want %q", v.file, opScope.UID, apimeta.PlatformScopeUID)
			}
		} else if op.Metadata.ScopeRef == nil || op.Metadata.ScopeRef.UID == "" {
			return fmt.Errorf("%s: non-platform scopeRef.uid required", v.file)
		}

		targetScope := resolveOperationTargetGovernanceScope(&op)
		if mismatch := apivalid.CheckOperationTargetScopeMatch(opScope, targetScope); mismatch != nil {
			return fmt.Errorf("%s: D-17 match failed: code=%s field=%s message=%s",
				v.file, mismatch.Code, mismatch.Field, mismatch.Message)
		}
	}

	// Canonical representative also participates (Platform nil scope).
	var canonical Operation
	if err := l.DecodeJSON("operation.json", &canonical); err != nil {
		return err
	}
	if err := assertOperationProvisioningShape(&canonical); err != nil {
		return fmt.Errorf("operation.json: %w", err)
	}
	if mismatch := apivalid.CheckOperationTargetScopeMatch(
		apimeta.CanonicalScopeIdentity(canonical.Metadata.ScopeRef),
		resolveOperationTargetGovernanceScope(&canonical),
	); mismatch != nil {
		return fmt.Errorf("operation.json: D-17 match failed: %#v", mismatch)
	}

	// Negative: kind mismatch and UID mismatch → OPERATION_TARGET_SCOPE_MISMATCH at /metadata/scopeRef.
	if err := assertOperationScopeMismatchFixture(l, "negative/operation-scope-target-kind-mismatch.json"); err != nil {
		return err
	}
	if err := assertOperationScopeMismatchFixture(l, "negative/operation-scope-target-uid-mismatch.json"); err != nil {
		return err
	}

	// Unavailable / unauthorized target → SafeDenial 404 without mismatch disclosure.
	denied := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	if denied == nil || absent == nil {
		return fmt.Errorf("SafeDenial/absent Problem must be non-nil")
	}
	deniedJSON, err := json.Marshal(denied)
	if err != nil {
		return err
	}
	absentJSON, err := json.Marshal(absent)
	if err != nil {
		return err
	}
	if !bytes.Equal(deniedJSON, absentJSON) {
		return fmt.Errorf("SafeDenial must be byte-identical to absent 404:\ndenied=%s\nabsent=%s", deniedJSON, absentJSON)
	}
	if denied.Status != 404 || denied.Code != apiproblem.CodeResourceNotFound {
		return fmt.Errorf("SafeDenial status/code = %d %q, want 404 RESOURCE_NOT_FOUND", denied.Status, denied.Code)
	}
	if len(denied.Violations) != 0 {
		return fmt.Errorf("SafeDenial must not disclose mismatch violations, got %#v", denied.Violations)
	}
	return nil
}

func provePortalListsLargeCollections(l *FixtureLoader) error {
	var item Project
	if err := l.DecodeJSON("project.json", &item); err != nil {
		return err
	}
	lim := apivalid.DefaultLimits()
	if lim.MaxPageSize <= 0 || lim.DefaultPageSize <= 0 || lim.DefaultPageSize > lim.MaxPageSize {
		return fmt.Errorf("list limits not bounded: default=%d max=%d", lim.DefaultPageSize, lim.MaxPageSize)
	}
	if lim.MaxPageSize > 200 {
		return fmt.Errorf("MaxPageSize=%d exceeds reviewed default bound 200", lim.MaxPageSize)
	}
	list := apimeta.ListEnvelope[Project]{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: item.APIVersion,
			Kind:       "ProjectList",
		},
		Items: []Project{item},
		Page: apimeta.Page{
			NextPageToken: "opaque-page-token-conformance-001",
		},
	}
	if len(list.Items) > lim.MaxPageSize {
		return fmt.Errorf("page size %d exceeds MaxPageSize %d", len(list.Items), lim.MaxPageSize)
	}
	raw, err := json.Marshal(list)
	if err != nil {
		return err
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return err
	}
	if _, ok := top["apiVersion"]; !ok {
		return fmt.Errorf("ListEnvelope must promote apiVersion to top level")
	}
	if _, ok := top["kind"]; !ok {
		return fmt.Errorf("ListEnvelope must promote kind to top level")
	}
	if _, nested := top["TypeMeta"]; nested {
		return fmt.Errorf("ListEnvelope must not nest TypeMeta")
	}
	if strings.Contains(list.Page.NextPageToken, "offset=") ||
		strings.Contains(strings.ToLower(list.Page.NextPageToken), "select ") {
		return fmt.Errorf("nextPageToken must remain opaque (no offsets/SQL)")
	}
	return nil
}

func proveProviderDisconnects(l *FixtureLoader) error {
	var db DiscoveredDatabase
	if err := l.DecodeJSON("discovered-database.json", &db); err != nil {
		return err
	}
	states := []ObservationState{
		ObservationStateCurrent,
		ObservationStateStale,
		ObservationStateUnknown,
		ObservationStateAbsent,
	}
	seen := map[ObservationState]struct{}{}
	for _, s := range states {
		seen[s] = struct{}{}
	}
	if len(seen) != 4 {
		return fmt.Errorf("observation states must remain four distinct values")
	}
	if _, ok := seen[db.Status.ObservationState]; !ok {
		return fmt.Errorf("fixture observationState %q not in closed vocabulary", db.Status.ObservationState)
	}
	// Distinctness: assigning each state yields a different value.
	var samples [4]ObservationState
	copy(samples[:], states)
	for i := 0; i < len(samples); i++ {
		for j := i + 1; j < len(samples); j++ {
			if samples[i] == samples[j] {
				return fmt.Errorf("observation states collide: %q", samples[i])
			}
		}
	}
	return nil
}

func proveExternalObjectRecreatedUnderSameName(l *FixtureLoader) error {
	var db DiscoveredDatabase
	if err := l.DecodeJSON("discovered-database.json", &db); err != nil {
		return err
	}
	if db.Metadata.Name == "" || db.Metadata.UID == "" {
		return fmt.Errorf("name and uid required")
	}
	if db.Metadata.Name == db.Metadata.UID {
		return fmt.Errorf("uid must be distinct from name to prevent stale rebinding")
	}
	if db.Status.ExternalName == "" {
		return fmt.Errorf("externalName required to model same-name recreation")
	}
	// Same external name with a different platform UID is a different object.
	reborn := db
	reborn.Metadata.UID = "ffffffffffffffffffffffffffffffff"
	if reborn.Metadata.UID == db.Metadata.UID {
		return fmt.Errorf("uid collision would allow stale-reference rebinding")
	}
	if reborn.Status.ExternalName != db.Status.ExternalName {
		return fmt.Errorf("external name should be comparable across recreation")
	}
	return nil
}

func proveCrossTenantReferenceAttempted(l *FixtureLoader) error {
	var p Project
	if err := l.DecodeJSON("project.json", &p); err != nil {
		return err
	}
	if p.Metadata.ScopeRef == nil {
		return fmt.Errorf("project scopeRef required as reference carrier")
	}
	denied := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	dj, err := json.Marshal(denied)
	if err != nil {
		return err
	}
	aj, err := json.Marshal(absent)
	if err != nil {
		return err
	}
	if !bytes.Equal(dj, aj) {
		return fmt.Errorf("cross-tenant denial must match absent 404 bytes")
	}
	if len(denied.Violations) != 0 || denied.Detail != absent.Detail {
		return fmt.Errorf("denial must not disclose target existence detail")
	}
	return nil
}

func proveNewCloudProviderAdded(l *FixtureLoader) error {
	var cfg AdapterConfiguration
	if err := l.DecodeJSON("adapter-configuration.json", &cfg); err != nil {
		return err
	}
	var pool ResourcePool
	if err := l.DecodeJSON("resource-pool.json", &pool); err != nil {
		return err
	}
	var project Project
	if err := l.DecodeJSON("project.json", &project); err != nil {
		return err
	}
	// Customer/core Project remains unchanged in shape (no adapter/native fields).
	raw, err := l.Read("project.json")
	if err != nil {
		return err
	}
	lower := strings.ToLower(string(raw))
	for _, banned := range []string{"nativeconfig", "adapterclass", "providerid", "aws", "azure", "gcp", "kubernetes"} {
		if strings.Contains(lower, banned) {
			return fmt.Errorf("customer Project fixture must not carry provider-native tokens (%s)", banned)
		}
	}
	if cfg.Spec.AdapterClass == "" || pool.Spec.CapabilityClass == "" {
		return fmt.Errorf("adapter/pool provider extension points required")
	}
	if project.Kind != "Project" {
		return fmt.Errorf("customer core kind unexpectedly changed")
	}
	return nil
}

func proveNewDataServicePluginAdded(l *FixtureLoader) error {
	var def PluginDefinition
	if err := l.DecodeJSON("plugin-definition.json", &def); err != nil {
		return err
	}
	var op Operation
	if err := l.DecodeJSON("operation-provider.json", &op); err != nil {
		return err
	}
	if def.APIVersion == "" || op.APIVersion == "" {
		return fmt.Errorf("explicit apiVersion required on definition and operation")
	}
	if def.Kind == "Project" || op.Kind == "Project" {
		return fmt.Errorf("plugin addition must not redefine core Project grammar")
	}
	if def.Spec.Version == "" || op.Spec.Action == "" {
		return fmt.Errorf("plugin contract + operation action prove extension without core redesign")
	}
	return nil
}

func proveAIExplainsDenial(l *FixtureLoader) error {
	var ev AuditEvent
	if err := l.DecodeJSON("audit-event.json", &ev); err != nil {
		return err
	}
	if ev.Record.ReasonCode == "" {
		return fmt.Errorf("stable reasonCode required (no message scraping)")
	}
	denied := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
	if denied.Code == "" || denied.Title == "" {
		return fmt.Errorf("Problem must expose stable code/title for explainability")
	}
	if denied.Code != apiproblem.CodeResourceNotFound {
		return fmt.Errorf("SafeDenial code=%q, want RESOURCE_NOT_FOUND", denied.Code)
	}
	// Safe context: no secret-like fields on the denial Problem.
	raw, err := json.Marshal(denied)
	if err != nil {
		return err
	}
	lower := strings.ToLower(string(raw))
	for _, banned := range []string{"password", "token", "private_key", "connectionstring"} {
		if strings.Contains(lower, banned) {
			return fmt.Errorf("denial Problem must not embed secrets (%s)", banned)
		}
	}
	_ = ev.Record.Outcome // outcome is structured context alongside reasonCode
	return nil
}

func provePhase1ResourceMigrated(l *FixtureLoader) error {
	var p Project
	if err := l.DecodeJSON("project.json", &p); err != nil {
		return err
	}
	group, version, ok := apimeta.ParseAPIVersion(p.APIVersion)
	if !ok || group == "" || version == "" {
		return fmt.Errorf("migration requires explicit group/version, got %q", p.APIVersion)
	}
	if !apimeta.IsKnownVersion(version) {
		return fmt.Errorf("version %q is not an approved explicit maturity form", version)
	}
	if p.Kind == "" || p.Metadata.UID == "" {
		return fmt.Errorf("migrated resource must retain stable kind/uid identity")
	}
	return nil
}

func assertOperationProvisioningShape(op *Operation) error {
	if op == nil {
		return fmt.Errorf("nil Operation")
	}
	if op.Spec.TargetRef.Name == "" || op.Spec.TargetRef.Kind == "" {
		return fmt.Errorf("targetRef required")
	}
	if op.Spec.Action == "" {
		return fmt.Errorf("action required")
	}
	if op.Spec.IdempotencyKey == "" {
		return fmt.Errorf("idempotencyKey required")
	}
	if op.Spec.RequestID == "" {
		return fmt.Errorf("requestId required for correlation")
	}
	if op.Status.Phase == "" {
		return fmt.Errorf("status.phase required")
	}
	if op.Status.Retryable == "" {
		return fmt.Errorf("status.retryable required")
	}
	if op.Status.CancelRequested == "" {
		return fmt.Errorf("status.cancelRequested required")
	}
	// Progress and terminal fields prove LongRunningOperation shape without execution.
	if op.Status.ProgressPercent < 0 || op.Status.ProgressPercent > 100 {
		return fmt.Errorf("progressPercent=%d out of range", op.Status.ProgressPercent)
	}
	switch op.Status.Phase {
	case OperationPhaseSucceeded, OperationPhaseFailed, OperationPhaseCancelled:
		if op.Status.TerminalCode == "" {
			return fmt.Errorf("terminal phase %q requires terminalCode", op.Status.Phase)
		}
	case OperationPhasePending, OperationPhaseRunning:
		// non-terminal; terminal fields may be empty
	default:
		return fmt.Errorf("unknown operation phase %q", op.Status.Phase)
	}
	return nil
}

// resolveOperationTargetGovernanceScope derives the authoritative target
// governance ScopeIdentity from the Operation fixture for D-17 checks.
// When targetRef.kind is itself a Matrix B scope kind, that kind+uid is used.
// Otherwise (PluginDefinition, ResourcePool, …) the Operation's canonical
// scopeRef (nil → Platform) is the fixture-declared governance scope.
func resolveOperationTargetGovernanceScope(op *Operation) apimeta.ScopeIdentity {
	kind := apimeta.ScopeKind(op.Spec.TargetRef.Kind)
	switch kind {
	case apimeta.ScopeOrganization,
		apimeta.ScopeOrganizationUnit,
		apimeta.ScopeTenant,
		apimeta.ScopeProject,
		apimeta.ScopeProvider:
		return apimeta.ScopeIdentity{Kind: kind, UID: op.Spec.TargetRef.UID}
	case apimeta.ScopePlatform:
		return apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}
	default:
		return apimeta.CanonicalScopeIdentity(op.Metadata.ScopeRef)
	}
}

func assertOperationScopeMismatchFixture(l *FixtureLoader, name string) error {
	var op Operation
	if err := l.DecodeJSON(name, &op); err != nil {
		return err
	}
	opScope := apimeta.CanonicalScopeIdentity(op.Metadata.ScopeRef)
	targetScope := resolveOperationTargetGovernanceScope(&op)
	v := apivalid.CheckOperationTargetScopeMatch(opScope, targetScope)
	if v == nil {
		return fmt.Errorf("%s: expected OPERATION_TARGET_SCOPE_MISMATCH, got nil (op=%+v target=%+v)",
			name, opScope, targetScope)
	}
	if v.Code != apiproblem.ViolationOperationTargetScopeMismatch {
		return fmt.Errorf("%s: code=%q, want OPERATION_TARGET_SCOPE_MISMATCH", name, v.Code)
	}
	if v.Field != "/metadata/scopeRef" {
		return fmt.Errorf("%s: field=%q, want /metadata/scopeRef", name, v.Field)
	}
	return nil
}

func rejectVendorNativeTokens(s string) error {
	lower := strings.ToLower(s)
	for _, banned := range []string{
		"aws_", "amazon.", "azure.", "gcp.", "google-cloud",
		"kubernetes.io/", "k8s.io/", "eks.", "aks.", "gke.",
	} {
		if strings.Contains(lower, banned) {
			return fmt.Errorf("provider-native token %q not allowed in core fixture", banned)
		}
	}
	return nil
}
