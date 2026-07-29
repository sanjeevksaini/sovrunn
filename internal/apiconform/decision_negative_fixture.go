package apiconform

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/validate"
)

// KindDecisionNegativeFixture is the on-disk T-030 negative fixture envelope
// kind. It is a conformance-only wrapper (not a domain resource).
const KindDecisionNegativeFixture = "DecisionNegativeFixture"

// parsedNegativeDecisionFixture is the decoded T-030 envelope.
type parsedNegativeDecisionFixture struct {
	WantViolationCode  apiproblem.ViolationCode
	WantViolationField string
	DecisionRecordJSON []byte
	BundleJSON         []byte
	RecordOptions      validate.DecisionRecordOptions
}

type negativeDecisionFixtureEnvelope struct {
	Kind               string          `json:"kind"`
	WantViolationCode  string          `json:"wantViolationCode"`
	WantViolationField string          `json:"wantViolationField"`
	Options            json.RawMessage `json:"options,omitempty"`
	Bundle             json.RawMessage `json:"bundle,omitempty"`
	DecisionRecord     json.RawMessage `json:"decisionRecord"`
}

type negativeFixtureOptionsDTO struct {
	ParallelScopeFields  []string         `json:"parallelScopeFields,omitempty"`
	TrustRequired        *bool            `json:"trustRequired,omitempty"`
	ExpectedTrustState   string           `json:"expectedTrustState,omitempty"`
	Classifications      []string         `json:"classifications,omitempty"`
	DeclaredSensitivity  string           `json:"declaredSensitivity,omitempty"`
	ContentCategories    []string         `json:"contentCategories,omitempty"`
	ProhibitedCategories []string         `json:"prohibitedCategories,omitempty"`
	CanonicalView        []string         `json:"canonicalView,omitempty"`
	AllowedAudiences     []string         `json:"allowedAudiences,omitempty"`
	SupportedObligations map[string]bool  `json:"supportedObligations,omitempty"`
	Established          []establishedDTO `json:"established,omitempty"`
}

type establishedDTO struct {
	Ref          apimeta.TypedRef               `json:"ref"`
	Relationship *decision.DecisionRelationship `json:"relationship,omitempty"`
}

func parseNegativeDecisionFixture(raw []byte) (parsedNegativeDecisionFixture, error) {
	var env negativeDecisionFixtureEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return parsedNegativeDecisionFixture{}, fmt.Errorf("decode DecisionNegativeFixture: %w", err)
	}
	if kind := strings.TrimSpace(env.Kind); kind != "" && kind != KindDecisionNegativeFixture {
		return parsedNegativeDecisionFixture{}, fmt.Errorf("unexpected negative fixture kind %q", kind)
	}
	code := apiproblem.ViolationCode(strings.TrimSpace(env.WantViolationCode))
	if code == "" || !validate.Valid(code) {
		return parsedNegativeDecisionFixture{}, fmt.Errorf("wantViolationCode %q is not a closed DECISION_* code", env.WantViolationCode)
	}
	if len(bytesTrimSpace(env.DecisionRecord)) == 0 {
		return parsedNegativeDecisionFixture{}, fmt.Errorf("decisionRecord is required")
	}
	out := parsedNegativeDecisionFixture{
		WantViolationCode:  code,
		WantViolationField: strings.TrimSpace(env.WantViolationField),
		DecisionRecordJSON: append([]byte(nil), env.DecisionRecord...),
		BundleJSON:         append([]byte(nil), bytesTrimSpace(env.Bundle)...),
	}
	if len(bytesTrimSpace(env.Options)) > 0 {
		opts, err := parseNegativeFixtureOptions(env.Options)
		if err != nil {
			return parsedNegativeDecisionFixture{}, err
		}
		out.RecordOptions = opts
	}
	return out, nil
}

func parseNegativeFixtureOptions(raw []byte) (validate.DecisionRecordOptions, error) {
	var dto negativeFixtureOptionsDTO
	if err := json.Unmarshal(raw, &dto); err != nil {
		return validate.DecisionRecordOptions{}, fmt.Errorf("decode negative fixture options: %w", err)
	}
	out := validate.DecisionRecordOptions{
		ParallelScopeFields:  append([]string(nil), dto.ParallelScopeFields...),
		TrustRequired:        dto.TrustRequired,
		ExpectedTrustState:   dto.ExpectedTrustState,
		DeclaredSensitivity:  decision.Sensitivity(dto.DeclaredSensitivity),
		ContentCategories:    append([]string(nil), dto.ContentCategories...),
		ProhibitedCategories: append([]string(nil), dto.ProhibitedCategories...),
		CanonicalView:        append([]string(nil), dto.CanonicalView...),
		AllowedAudiences:     append([]string(nil), dto.AllowedAudiences...),
	}
	if len(dto.SupportedObligations) > 0 {
		out.SupportedObligations = make(map[string]bool, len(dto.SupportedObligations))
		for k, v := range dto.SupportedObligations {
			out.SupportedObligations[k] = v
		}
	}
	if len(dto.Classifications) > 0 {
		out.Classifications = make([]apimeta.DataClassification, 0, len(dto.Classifications))
		for _, c := range dto.Classifications {
			out.Classifications = append(out.Classifications, apimeta.DataClassification(c))
		}
	}
	if len(dto.Established) > 0 {
		out.Established = make([]validate.RelationshipNode, 0, len(dto.Established))
		for _, n := range dto.Established {
			out.Established = append(out.Established, validate.RelationshipNode{
				Ref:          n.Ref,
				Relationship: n.Relationship,
			})
		}
	}
	return out, nil
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}
