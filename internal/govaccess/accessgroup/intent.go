package accessgroup

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

type PreparedAccessGroupIntent struct {
	sealed          bool
	claim           state.ParticipantClaim
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            model.AccessGroupSpec
	expectedVersion string
	nextVersion     string
}

func (i PreparedAccessGroupIntent) Sealed() bool                             { return i.sealed }
func (i PreparedAccessGroupIntent) ParticipantClaim() state.ParticipantClaim { return i.claim }
func (i PreparedAccessGroupIntent) HasPublicationDerivedFields() bool        { return false }

func NewPreparedAccessGroupIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.AccessGroupSpec,
	expectedVersion string,
	nextVersion string,
) (PreparedAccessGroupIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return PreparedAccessGroupIntent{}, operation.NewMechanicalFailure("accessgroup_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return PreparedAccessGroupIntent{}, operation.NewMechanicalFailure("accessgroup_versions_incomplete")
	}
	if spec.OwnerRef.PrincipalType != model.PrincipalTypeHuman {
		return PreparedAccessGroupIntent{}, operation.NewMechanicalFailure("accessgroup_owner_must_be_human")
	}
	raw, err := json.Marshal(struct {
		ResourceUID     string                `json:"resourceUid"`
		ResourceName    string                `json:"resourceName"`
		ScopeRef        apimeta.ScopeRef      `json:"scopeRef"`
		Spec            model.AccessGroupSpec `json:"spec"`
		ExpectedVersion string                `json:"expectedVersion"`
		NextVersion     string                `json:"nextVersion"`
	}{
		ResourceUID:     resourceUID,
		ResourceName:    resourceName,
		ScopeRef:        scopeRef,
		Spec:            spec,
		ExpectedVersion: expectedVersion,
		NextVersion:     nextVersion,
	})
	if err != nil {
		return PreparedAccessGroupIntent{}, operation.NewMechanicalFailure("accessgroup_intent_marshal")
	}
	sum := sha256.Sum256(raw)
	versionSet, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{
		ResourceUID: resourceUID, ResourceVersion: expectedVersion,
	}})
	if err != nil {
		return PreparedAccessGroupIntent{}, operation.NewMechanicalFailure("accessgroup_versionset_invalid")
	}
	claim, err := state.NewParticipantClaim(
		state.AccessGroupOwner, participantID, sum[:], versionSet, state.OriginatingResult,
	)
	if err != nil {
		return PreparedAccessGroupIntent{}, operation.NewMechanicalFailure("accessgroup_claim_invalid")
	}
	return PreparedAccessGroupIntent{
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
