package exception

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

type DomainOutcome struct {
	sealed bool
	code   string
}

func (o DomainOutcome) Sealed() bool { return o.sealed }
func (o DomainOutcome) Code() string { return o.code }

// ApprovalRequirementDeriver owns Exception approval derivation.
type ApprovalRequirementDeriver struct{}

func (ApprovalRequirementDeriver) DeriverKind() string { return "Exception" }

func (ApprovalRequirementDeriver) DeriveForException(approvalreq.ExceptionInput) (model.ApprovalRequirement, bool) {
	return model.NewApprovalRequirementRequired(apimeta.TypedRef{
		APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "exception-default", UID: "approvalpolicy-exception-default",
	})
}

type PreparedExceptionIntent struct {
	sealed          bool
	claim           state.ParticipantClaim
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	record          model.ExceptionGrantRecord
	decisionFacts   model.ExceptionDecisionFacts
	expectedVersion string
	nextVersion     string
}

func (i PreparedExceptionIntent) Sealed() bool                             { return i.sealed }
func (i PreparedExceptionIntent) ParticipantClaim() state.ParticipantClaim { return i.claim }
func (i PreparedExceptionIntent) HasPublicationDerivedFields() bool        { return false }

func NewPreparedExceptionIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	record model.ExceptionGrantRecord,
	decisionFacts model.ExceptionDecisionFacts,
	expectedVersion string,
	nextVersion string,
) (PreparedExceptionIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return PreparedExceptionIntent{}, operation.NewMechanicalFailure("exception_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return PreparedExceptionIntent{}, operation.NewMechanicalFailure("exception_versions_incomplete")
	}
	if !decisionFacts.Sealed() {
		return PreparedExceptionIntent{}, operation.NewMechanicalFailure("exception_decisionfacts_unsealed")
	}
	raw, err := json.Marshal(struct {
		ResourceUID     string                       `json:"resourceUid"`
		ResourceName    string                       `json:"resourceName"`
		ScopeRef        apimeta.ScopeRef             `json:"scopeRef"`
		Record          model.ExceptionGrantRecord   `json:"record"`
		DecisionFacts   model.ExceptionDecisionFacts `json:"decisionFacts"`
		ExpectedVersion string                       `json:"expectedVersion"`
		NextVersion     string                       `json:"nextVersion"`
	}{
		ResourceUID:     resourceUID,
		ResourceName:    resourceName,
		ScopeRef:        scopeRef,
		Record:          record,
		DecisionFacts:   decisionFacts,
		ExpectedVersion: expectedVersion,
		NextVersion:     nextVersion,
	})
	if err != nil {
		return PreparedExceptionIntent{}, operation.NewMechanicalFailure("exception_intent_marshal")
	}
	sum := sha256.Sum256(raw)
	versionSet, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{
		ResourceUID: resourceUID, ResourceVersion: expectedVersion,
	}})
	if err != nil {
		return PreparedExceptionIntent{}, operation.NewMechanicalFailure("exception_versionset_invalid")
	}
	claim, err := state.NewParticipantClaim(state.ExceptionOwner, participantID, sum[:], versionSet, state.OriginatingResult)
	if err != nil {
		return PreparedExceptionIntent{}, operation.NewMechanicalFailure("exception_claim_invalid")
	}
	return PreparedExceptionIntent{
		sealed:          true,
		claim:           claim,
		resourceUID:     resourceUID,
		resourceName:    resourceName,
		scopeRef:        scopeRef,
		record:          record,
		decisionFacts:   decisionFacts,
		expectedVersion: expectedVersion,
		nextVersion:     nextVersion,
	}, operation.MechanicalFailure{}
}
