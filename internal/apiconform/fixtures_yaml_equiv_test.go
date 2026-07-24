package apiconform

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// TestValidYAMLFixturesDecodeEquivalentToJSON verifies Task 14.2: each
// canonical JSON fixture has a strict JSON-compatible YAML twin, and
// DecodeYAML under ModeReadRepresentation produces an equivalent typed
// value to DecodeJSON (D-03a; F12-VALIDATION-001(2)).
func TestValidYAMLFixturesDecodeEquivalentToJSON(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)

	cases := []struct {
		base   string // fixture basename without extension
		newDst func() any
	}{
		{base: "project", newDst: func() any { return &Project{} }},
		{base: "resource-pool", newDst: func() any { return &ResourcePool{} }},
		{base: "discovered-database", newDst: func() any { return &DiscoveredDatabase{} }},
		{base: "plugin-definition", newDst: func() any { return &PluginDefinition{} }},
		{base: "adapter-configuration", newDst: func() any { return &AdapterConfiguration{} }},
		{base: "placement-evaluation-request", newDst: func() any { return &PlacementEvaluationRequest{} }},
		{base: "audit-event", newDst: func() any { return &AuditEvent{} }},
		{base: "operation", newDst: func() any { return &Operation{} }},
		{base: "operation-platform", newDst: func() any { return &Operation{} }},
		{base: "operation-organization", newDst: func() any { return &Operation{} }},
		{base: "operation-organizationunit", newDst: func() any { return &Operation{} }},
		{base: "operation-tenant", newDst: func() any { return &Operation{} }},
		{base: "operation-project", newDst: func() any { return &Operation{} }},
		{base: "operation-provider", newDst: func() any { return &Operation{} }},
	}

	if len(cases) != 14 {
		t.Fatalf("expected 14 Task 14.2 fixture pairs, got %d", len(cases))
	}

	dir := filepath.Join(root, ConformanceFixturesDir)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.base, func(t *testing.T) {
			t.Parallel()

			jsonPath := filepath.Join(dir, tc.base+".json")
			yamlPath := filepath.Join(dir, tc.base+".yaml")

			jsonRaw, err := os.ReadFile(jsonPath)
			if err != nil {
				t.Fatalf("read JSON fixture: %v", err)
			}
			yamlRaw, err := os.ReadFile(yamlPath)
			if err != nil {
				t.Fatalf("read YAML fixture: %v", err)
			}
			if len(strings.TrimSpace(string(yamlRaw))) == 0 {
				t.Fatalf("YAML fixture %s is empty", yamlPath)
			}

			fromJSON := tc.newDst()
			if prob := apivalid.DecodeJSON(jsonRaw, lim, pol, fromJSON); prob != nil {
				t.Fatalf("DecodeJSON: code=%s detail=%s violations=%v",
					prob.Code, prob.Detail, prob.Violations)
			}
			fromYAML := tc.newDst()
			if prob := apivalid.DecodeYAML(yamlRaw, lim, pol, fromYAML); prob != nil {
				t.Fatalf("DecodeYAML: code=%s detail=%s violations=%v",
					prob.Code, prob.Detail, prob.Violations)
			}
			if !reflect.DeepEqual(fromYAML, fromJSON) {
				t.Fatalf("YAML decode diverges from JSON decode:\nYAML=%#v\nJSON=%#v",
					fromYAML, fromJSON)
			}
		})
	}
}
