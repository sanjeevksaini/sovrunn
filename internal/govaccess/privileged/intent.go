package privileged

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"reflect"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/policyseam"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

type DomainOutcome struct {
	sealed bool
	code   string
}

func (o DomainOutcome) Sealed() bool { return o.sealed }
func (o DomainOutcome) Code() string { return o.code }

// ApprovalRequirementDeriver owns Privileged approval derivation.
type ApprovalRequirementDeriver struct{}

func (ApprovalRequirementDeriver) DeriverKind() string { return "Privileged" }

func (ApprovalRequirementDeriver) DeriveForPrivileged(approvalreq.PrivilegedInput) (model.ApprovalRequirement, bool) {
	return model.NewApprovalRequirementRequired(apimeta.TypedRef{
		APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "privileged-default", UID: "approvalpolicy-privileged-default",
	})
}

type PreparedPrivilegedIntent struct {
	sealed           bool
	claim            state.ParticipantClaim
	resourceUID      string
	resourceName     string
	scopeRef         apimeta.ScopeRef
	spec             model.PrivilegedAccessRequestSpec
	assurance        model.AssuranceEvidence
	expectedVersion  string
	nextVersion      string
	policyEvaluation policyseam.Evaluation
}

func (i PreparedPrivilegedIntent) Sealed() bool                             { return i.sealed }
func (i PreparedPrivilegedIntent) ParticipantClaim() state.ParticipantClaim { return i.claim }
func (i PreparedPrivilegedIntent) HasPublicationDerivedFields() bool        { return false }

func NewPreparedPrivilegedIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.PrivilegedAccessRequestSpec,
	assurance model.AssuranceEvidence,
	port policyseam.Port,
	expectedVersion string,
	nextVersion string,
) (PreparedPrivilegedIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return PreparedPrivilegedIntent{}, operation.NewMechanicalFailure("privileged_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return PreparedPrivilegedIntent{}, operation.NewMechanicalFailure("privileged_versions_incomplete")
	}
	if spec.Mode == "" || spec.ActivationDeadline.IsZero() {
		return PreparedPrivilegedIntent{}, operation.NewMechanicalFailure("privileged_spec_invalid")
	}
	evaluation, fail := evaluateWithZeroPrivilegedRequest(port)
	if fail.Reason() != "" {
		return PreparedPrivilegedIntent{}, fail
	}
	raw, err := json.Marshal(struct {
		ResourceUID     string                            `json:"resourceUid"`
		ResourceName    string                            `json:"resourceName"`
		ScopeRef        apimeta.ScopeRef                  `json:"scopeRef"`
		Spec            model.PrivilegedAccessRequestSpec `json:"spec"`
		Assurance       model.AssuranceEvidence           `json:"assurance"`
		Eval            policyseam.Evaluation             `json:"evaluation"`
		ExpectedVersion string                            `json:"expectedVersion"`
		NextVersion     string                            `json:"nextVersion"`
	}{
		ResourceUID:     resourceUID,
		ResourceName:    resourceName,
		ScopeRef:        scopeRef,
		Spec:            spec,
		Assurance:       assurance,
		Eval:            evaluation,
		ExpectedVersion: expectedVersion,
		NextVersion:     nextVersion,
	})
	if err != nil {
		return PreparedPrivilegedIntent{}, operation.NewMechanicalFailure("privileged_intent_marshal")
	}
	sum := sha256.Sum256(raw)
	versionSet, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{
		ResourceUID: resourceUID, ResourceVersion: expectedVersion,
	}})
	if err != nil {
		return PreparedPrivilegedIntent{}, operation.NewMechanicalFailure("privileged_versionset_invalid")
	}
	claim, err := state.NewParticipantClaim(state.PrivilegedOwner, participantID, sum[:], versionSet, state.OriginatingResult)
	if err != nil {
		return PreparedPrivilegedIntent{}, operation.NewMechanicalFailure("privileged_claim_invalid")
	}
	return PreparedPrivilegedIntent{
		sealed:           true,
		claim:            claim,
		resourceUID:      resourceUID,
		resourceName:     resourceName,
		scopeRef:         scopeRef,
		spec:             spec,
		assurance:        assurance,
		expectedVersion:  expectedVersion,
		nextVersion:      nextVersion,
		policyEvaluation: evaluation,
	}, operation.MechanicalFailure{}
}

func evaluateWithZeroPrivilegedRequest(port policyseam.Port) (policyseam.Evaluation, operation.MechanicalFailure) {
	method := reflect.ValueOf(port).MethodByName("EvaluatePrivileged")
	if !method.IsValid() || method.Type().NumIn() != 2 || method.Type().NumOut() != 2 {
		return policyseam.Evaluation{}, operation.NewMechanicalFailure("policyseam_method_unavailable")
	}
	callIn := []reflect.Value{
		reflect.ValueOf(context.Background()),
		reflect.Zero(method.Type().In(1)),
	}
	out := method.Call(callIn)
	eval, ok := out[0].Interface().(policyseam.Evaluation)
	if !ok {
		return policyseam.Evaluation{}, operation.NewMechanicalFailure("policyseam_eval_cast")
	}
	fail, ok := out[1].Interface().(operation.MechanicalFailure)
	if !ok {
		return policyseam.Evaluation{}, operation.NewMechanicalFailure("policyseam_failure_cast")
	}
	return eval, fail
}
