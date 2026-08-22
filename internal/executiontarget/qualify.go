package executiontarget

import (
	"fmt"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// Qualification evaluation reason codes (closed profile; REQ-F16-07).
const (
	ReasonProfileQualified      = "profile.qualified"
	ReasonProfileRejected       = "profile.rejected"
	ReasonProfileIndeterminate  = "profile.indeterminate"
	ReasonFixtureMissing        = "fixture.missing"
	ReasonFixtureLogicalTimeout = "fixture.logical_timeout"
)

// ObserverFault classifies malformed observer output that faults qualification
// without publishing a result (REQ-F16-07; design §5.1 step 12).
type ObserverFault struct {
	Detail string
}

func (f ObserverFault) Error() string {
	if f.Detail == "" {
		return "observer fault"
	}
	return "observer fault: " + f.Detail
}

// ProfileEvaluation is the closed synthetic-iaas profile outcome.
type ProfileEvaluation struct {
	Outcome     model.QualificationOutcome
	ReasonCodes []string
	Truths      model.FactSetTruths
	EvaluatedAt time.Time
}

// EvaluateObservation applies the closed profile to an observer proposal:
// any Unsupported → Rejected; else any Unknown → Indeterminate; else Qualified.
// Duplicate, missing, or malformed facts/provenance return ObserverFault.
func EvaluateObservation(obs Observation, evaluatedAt time.Time) (ProfileEvaluation, error) {
	if err := validateObservationProvenance(obs); err != nil {
		return ProfileEvaluation{}, err
	}
	truths, err := normalizeFacts(obs.Facts)
	if err != nil {
		return ProfileEvaluation{}, err
	}

	eval := ProfileEvaluation{
		Truths:      truths,
		EvaluatedAt: evaluatedAt.UTC(),
	}

	if hasTruth(truths, model.FactUnsupported) {
		eval.Outcome = model.OutcomeRejected
		eval.ReasonCodes = []string{ReasonProfileRejected}
		return eval, nil
	}
	if hasTruth(truths, model.FactUnknown) {
		eval.Outcome = model.OutcomeIndeterminate
		switch {
		case obs.MissingFixture:
			eval.ReasonCodes = []string{ReasonFixtureMissing, ReasonProfileIndeterminate}
		case obs.LogicalTimeout:
			eval.ReasonCodes = []string{ReasonFixtureLogicalTimeout, ReasonProfileIndeterminate}
		default:
			eval.ReasonCodes = []string{ReasonProfileIndeterminate}
		}
		return eval, nil
	}
	eval.Outcome = model.OutcomeQualified
	eval.ReasonCodes = []string{ReasonProfileQualified}
	return eval, nil
}

func validateObservationProvenance(obs Observation) error {
	if obs.ObserverID != model.ObserverID {
		return ObserverFault{Detail: "observerID must be " + model.ObserverID}
	}
	if obs.FactSchemaVersion != model.FactSchemaVersionV1 {
		return ObserverFault{Detail: "factSchemaVersion must be v1"}
	}
	if obs.ObserverRevision == "" && obs.FixtureRevision == "" {
		return ObserverFault{Detail: "observerRevision is required"}
	}
	return nil
}

func normalizeFacts(facts []ProposedFact) (model.FactSetTruths, error) {
	required := []string{
		FactNameComputeVM,
		FactNameStorageBlock,
		FactNameStorageObject,
		FactNameNetworkPrivate,
	}
	seen := make(map[string]model.FactTruth, len(required))
	for _, f := range facts {
		if f.Name == "" {
			return model.FactSetTruths{}, ObserverFault{Detail: "fact name is required"}
		}
		if !validFactTruth(f.Truth) {
			return model.FactSetTruths{}, ObserverFault{Detail: fmt.Sprintf("malformed truth for %s", f.Name)}
		}
		if _, dup := seen[f.Name]; dup {
			return model.FactSetTruths{}, ObserverFault{Detail: "duplicate fact " + f.Name}
		}
		seen[f.Name] = f.Truth
	}
	for _, name := range required {
		if _, ok := seen[name]; !ok {
			return model.FactSetTruths{}, ObserverFault{Detail: "missing fact " + name}
		}
	}
	if len(seen) != len(required) {
		return model.FactSetTruths{}, ObserverFault{Detail: "unexpected fact names"}
	}
	return model.FactSetTruths{
		Compute: model.FactComputeTruths{VM: seen[FactNameComputeVM]},
		Storage: model.FactStorageTruths{
			Block:  seen[FactNameStorageBlock],
			Object: seen[FactNameStorageObject],
		},
		Network: model.FactNetworkTruths{Private: seen[FactNameNetworkPrivate]},
	}, nil
}

func validFactTruth(t model.FactTruth) bool {
	switch t {
	case model.FactSupported, model.FactUnsupported, model.FactUnknown:
		return true
	default:
		return false
	}
}

func hasTruth(t model.FactSetTruths, want model.FactTruth) bool {
	return t.Compute.VM == want ||
		t.Storage.Block == want ||
		t.Storage.Object == want ||
		t.Network.Private == want
}

// BuildFactSet constructs an immutable NormalizedTargetFactSet from a successful
// profile evaluation and the captured fences for this attempt.
func BuildFactSet(
	uid string,
	target model.ExecutionTarget,
	obs Observation,
	eval ProfileEvaluation,
	stackGeneration int64,
	maintenanceEpoch int64,
	viability model.ViabilityFingerprint,
) model.NormalizedTargetFactSet {
	revision := obs.ObserverRevision
	if revision == "" {
		revision = obs.FixtureRevision
	}
	if revision == "" {
		revision = RegisteredObserverRevision
	}
	return model.NormalizedTargetFactSet{
		UID: uid,
		TargetRef: apimeta.TypedRef{
			APIVersion: model.APIVersionExecutionTarget,
			Kind:       model.KindExecutionTarget,
			Name:       target.Metadata.Name,
			UID:        target.Metadata.UID,
		},
		InfrastructureStackGeneration: stackGeneration,
		MaintenanceEpoch:              maintenanceEpoch,
		ViabilityFingerprint:          viability,
		ObserverID:                    model.ObserverID,
		ObserverRevision:              revision,
		FactVersion:                   model.FactVersionV1,
		ObservedAt:                    obs.ObservedAt.UTC().Format(time.RFC3339),
		ExpiresAt:                     obs.ExpiresAt.UTC().Format(time.RFC3339),
		Facts:                         eval.Truths,
	}
}

// BuildQualificationResult constructs an immutable TargetQualificationResult.
func BuildQualificationResult(
	uid string,
	target model.ExecutionTarget,
	factSetUID string,
	eval ProfileEvaluation,
	stackGeneration int64,
	maintenanceEpoch int64,
) model.TargetQualificationResult {
	reasons := append([]string(nil), eval.ReasonCodes...)
	return model.TargetQualificationResult{
		UID: uid,
		TargetRef: apimeta.TypedRef{
			APIVersion: model.APIVersionExecutionTarget,
			Kind:       model.KindExecutionTarget,
			Name:       target.Metadata.Name,
			UID:        target.Metadata.UID,
		},
		FactSetRef: apimeta.TypedRef{
			APIVersion: model.APIVersionFactSet,
			Kind:       model.KindNormalizedTargetFactSet,
			Name:       factSetUID,
			UID:        factSetUID,
		},
		InfrastructureStackGeneration: stackGeneration,
		MaintenanceEpoch:              maintenanceEpoch,
		ProfileVersion:                model.ProfileVersion,
		Outcome:                       eval.Outcome,
		ReasonCodes:                   reasons,
		EvaluatedAt:                   eval.EvaluatedAt.UTC().Format(time.RFC3339),
	}
}

func qualificationStateForOutcome(outcome model.QualificationOutcome) model.QualificationState {
	switch outcome {
	case model.OutcomeQualified:
		return model.QualificationQualified
	case model.OutcomeRejected:
		return model.QualificationRejected
	case model.OutcomeIndeterminate:
		return model.QualificationIndeterminate
	default:
		return model.QualificationUnqualified
	}
}
