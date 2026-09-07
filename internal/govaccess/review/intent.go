package review

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
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

type PreparedReviewIntent struct {
	sealed          bool
	claim           state.ParticipantClaim
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            model.AccessReviewSpec
	expectedVersion string
	nextVersion     string
	beneficiarySet  []model.PrincipalRef
}

func (i PreparedReviewIntent) Sealed() bool                             { return i.sealed }
func (i PreparedReviewIntent) ParticipantClaim() state.ParticipantClaim { return i.claim }
func (i PreparedReviewIntent) HasPublicationDerivedFields() bool        { return false }

func NewPreparedReviewIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.AccessReviewSpec,
	expectedVersion string,
	nextVersion string,
	beneficiaries []model.PrincipalRef,
) (PreparedReviewIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return PreparedReviewIntent{}, operation.NewMechanicalFailure("review_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return PreparedReviewIntent{}, operation.NewMechanicalFailure("review_versions_incomplete")
	}
	if len(spec.ReviewerEligibility) == 0 || spec.AccessReviewRuleRef.UID == "" {
		return PreparedReviewIntent{}, operation.NewMechanicalFailure("review_spec_invalid")
	}
	raw, err := json.Marshal(struct {
		ResourceUID     string                 `json:"resourceUid"`
		ResourceName    string                 `json:"resourceName"`
		ScopeRef        apimeta.ScopeRef       `json:"scopeRef"`
		Spec            model.AccessReviewSpec `json:"spec"`
		ExpectedVersion string                 `json:"expectedVersion"`
		NextVersion     string                 `json:"nextVersion"`
		Beneficiaries   []model.PrincipalRef   `json:"beneficiaries"`
	}{
		ResourceUID:     resourceUID,
		ResourceName:    resourceName,
		ScopeRef:        scopeRef,
		Spec:            spec,
		ExpectedVersion: expectedVersion,
		NextVersion:     nextVersion,
		Beneficiaries:   append([]model.PrincipalRef(nil), beneficiaries...),
	})
	if err != nil {
		return PreparedReviewIntent{}, operation.NewMechanicalFailure("review_intent_marshal")
	}
	sum := sha256.Sum256(raw)
	versionSet, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{
		ResourceUID: resourceUID, ResourceVersion: expectedVersion,
	}})
	if err != nil {
		return PreparedReviewIntent{}, operation.NewMechanicalFailure("review_versionset_invalid")
	}
	claim, err := state.NewParticipantClaim(state.ReviewOwner, participantID, sum[:], versionSet, state.OriginatingResult)
	if err != nil {
		return PreparedReviewIntent{}, operation.NewMechanicalFailure("review_claim_invalid")
	}
	return PreparedReviewIntent{
		sealed:          true,
		claim:           claim,
		resourceUID:     resourceUID,
		resourceName:    resourceName,
		scopeRef:        scopeRef,
		spec:            spec,
		expectedVersion: expectedVersion,
		nextVersion:     nextVersion,
		beneficiarySet:  append([]model.PrincipalRef(nil), beneficiaries...),
	}, operation.MechanicalFailure{}
}

func hasBeneficiaryConflict(reviewer model.PrincipalRef, beneficiaries []model.PrincipalRef) bool {
	for _, b := range beneficiaries {
		if reviewer.Equal(b) {
			return true
		}
	}
	return false
}
