package roledefinition

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

type PreparedRoleDefinitionIntent struct {
	sealed          bool
	claim           state.ParticipantClaim
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            model.RoleDefinitionSpec
	expectedVersion string
	nextVersion     string
}

func (i PreparedRoleDefinitionIntent) Sealed() bool                             { return i.sealed }
func (i PreparedRoleDefinitionIntent) ParticipantClaim() state.ParticipantClaim { return i.claim }
func (i PreparedRoleDefinitionIntent) HasPublicationDerivedFields() bool        { return false }

func NewPreparedRoleDefinitionIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.RoleDefinitionSpec,
	expectedVersion string,
	nextVersion string,
) (PreparedRoleDefinitionIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return PreparedRoleDefinitionIntent{}, operation.NewMechanicalFailure("roledefinition_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return PreparedRoleDefinitionIntent{}, operation.NewMechanicalFailure("roledefinition_versions_incomplete")
	}
	if spec.Version == "" || len(spec.Actions) == 0 {
		return PreparedRoleDefinitionIntent{}, operation.NewMechanicalFailure("roledefinition_spec_invalid")
	}
	raw, err := json.Marshal(struct {
		ResourceUID     string                   `json:"resourceUid"`
		ResourceName    string                   `json:"resourceName"`
		ScopeRef        apimeta.ScopeRef         `json:"scopeRef"`
		Spec            model.RoleDefinitionSpec `json:"spec"`
		ExpectedVersion string                   `json:"expectedVersion"`
		NextVersion     string                   `json:"nextVersion"`
	}{
		ResourceUID:     resourceUID,
		ResourceName:    resourceName,
		ScopeRef:        scopeRef,
		Spec:            spec,
		ExpectedVersion: expectedVersion,
		NextVersion:     nextVersion,
	})
	if err != nil {
		return PreparedRoleDefinitionIntent{}, operation.NewMechanicalFailure("roledefinition_intent_marshal")
	}
	sum := sha256.Sum256(raw)
	versionSet, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{
		ResourceUID: resourceUID, ResourceVersion: expectedVersion,
	}})
	if err != nil {
		return PreparedRoleDefinitionIntent{}, operation.NewMechanicalFailure("roledefinition_versionset_invalid")
	}
	claim, err := state.NewParticipantClaim(
		state.RoleDefinitionOwner, participantID, sum[:], versionSet, state.OriginatingResult,
	)
	if err != nil {
		return PreparedRoleDefinitionIntent{}, operation.NewMechanicalFailure("roledefinition_claim_invalid")
	}
	return PreparedRoleDefinitionIntent{
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

func effectiveClassification(spec model.RoleDefinitionSpec) model.RoleClassification {
	if spec.BaseClassification.Valid() {
		return spec.BaseClassification
	}
	if spec.SupersedesRef != nil {
		return model.RoleClassificationPrivileged
	}
	return model.RoleClassificationOrdinary
}
