package governanceprofile

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

type PreparedGovernanceProfileIntent struct {
	sealed          bool
	claim           state.ParticipantClaim
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            model.GovernanceProfileSpec
	expectedVersion string
	nextVersion     string
}

func (i PreparedGovernanceProfileIntent) Sealed() bool                             { return i.sealed }
func (i PreparedGovernanceProfileIntent) ParticipantClaim() state.ParticipantClaim { return i.claim }
func (i PreparedGovernanceProfileIntent) HasPublicationDerivedFields() bool        { return false }

func NewPreparedGovernanceProfileIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.GovernanceProfileSpec,
	expectedVersion string,
	nextVersion string,
) (PreparedGovernanceProfileIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return PreparedGovernanceProfileIntent{}, operation.NewMechanicalFailure("governanceprofile_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return PreparedGovernanceProfileIntent{}, operation.NewMechanicalFailure("governanceprofile_versions_incomplete")
	}
	if spec.Version == "" || !spec.PublicationState.Valid() {
		return PreparedGovernanceProfileIntent{}, operation.NewMechanicalFailure("governanceprofile_spec_invalid")
	}
	raw, err := json.Marshal(struct {
		ResourceUID     string                      `json:"resourceUid"`
		ResourceName    string                      `json:"resourceName"`
		ScopeRef        apimeta.ScopeRef            `json:"scopeRef"`
		Spec            model.GovernanceProfileSpec `json:"spec"`
		ExpectedVersion string                      `json:"expectedVersion"`
		NextVersion     string                      `json:"nextVersion"`
	}{
		ResourceUID:     resourceUID,
		ResourceName:    resourceName,
		ScopeRef:        scopeRef,
		Spec:            spec,
		ExpectedVersion: expectedVersion,
		NextVersion:     nextVersion,
	})
	if err != nil {
		return PreparedGovernanceProfileIntent{}, operation.NewMechanicalFailure("governanceprofile_intent_marshal")
	}
	sum := sha256.Sum256(raw)
	versionSet, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{
		ResourceUID: resourceUID, ResourceVersion: expectedVersion,
	}})
	if err != nil {
		return PreparedGovernanceProfileIntent{}, operation.NewMechanicalFailure("governanceprofile_versionset_invalid")
	}
	claim, err := state.NewParticipantClaim(state.GovernanceProfileOwner, participantID, sum[:], versionSet, state.OriginatingResult)
	if err != nil {
		return PreparedGovernanceProfileIntent{}, operation.NewMechanicalFailure("governanceprofile_claim_invalid")
	}
	return PreparedGovernanceProfileIntent{
		sealed:          true,
		claim:           claim,
		resourceUID:     resourceUID,
		resourceName:    resourceName,
		scopeRef:        scopeRef,
		spec:            spec,
		expectedVersion: expectedVersion,
		nextVersion:     nextVersion,
	}, operation.MechanicalFailure{}
}
