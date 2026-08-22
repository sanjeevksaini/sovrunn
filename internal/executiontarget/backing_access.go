package executiontarget

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	cmmodel "github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// BackingGrantEvaluator is the inherited FEATURE-0012 grant surface consumed by
// BackingAccessProvider. Handlers adapt their GrantResolver to this interface.
type BackingGrantEvaluator interface {
	AuthorizeExact(action string, scope apimeta.ScopeIdentity, resourceUID, relatedPlatformUID string) bool
}

// BackingDisposition is the closed F0016-internal result of principal-aware
// backing access (ADH-2026-066; DD-13.2). Handlers never receive a lease, raw
// backing object, or raw existence boolean.
type BackingDisposition int

const (
	// BackingSafeDenied maps to the existing audited safe 404.
	BackingSafeDenied BackingDisposition = iota
	// BackingAuthorizationDenied maps to the existing audited 403.
	BackingAuthorizationDenied
	// BackingAllowed carries an immutable paired snapshot for route use.
	BackingAllowed
)

// BackingView is the immutable allowed snapshot returned by BackingAccessProvider.
// Participation generation is intentionally absent: it is never a fence.
type BackingView struct {
	ParticipationUID              string
	ParticipationEffectiveActive  bool
	ParticipationCloudProviderUID string
	StackUID                      string
	StackPhase                    string
	StackScopeUID                 string
	StackGeneration               int64
}

// StackActive reports whether the InfrastructureStack phase is Active.
func (v BackingView) StackActive() bool {
	return v.StackPhase == string(cmmodel.InfrastructureStackPhaseActive)
}

// Fingerprint returns the closed viability fingerprint inputs.
func (v BackingView) Fingerprint() model.ViabilityFingerprint {
	return model.ComputeViabilityFingerprint(
		v.ParticipationUID,
		v.ParticipationEffectiveActive,
		v.StackUID,
		v.StackPhase,
	)
}

// AsViability converts the allowed view into the lifecycle BackingViability shape.
func (v BackingView) AsViability() BackingViability {
	return BackingViability{
		ParticipationUID:             v.ParticipationUID,
		ParticipationEffectiveActive: v.ParticipationEffectiveActive,
		ParticipationScopeUID:        v.ParticipationCloudProviderUID,
		StackUID:                     v.StackUID,
		StackPhase:                   v.StackPhase,
		StackScopeUID:                v.StackScopeUID,
		StackGeneration:              v.StackGeneration,
	}
}

// BackingAccessResult is the closed union returned by Access.
type BackingAccessResult struct {
	Disposition BackingDisposition
	View        BackingView // meaningful only when Disposition == BackingAllowed
}

// BackingAccessProvider is the F0016-owned, principal-aware sole resolver for
// FEATURE-0015 backing reads used by create/qualify/GET/LIST (ADH-2026-066).
type BackingAccessProvider struct {
	store *cloudmodel.Store
}

// NewBackingAccessProvider constructs a provider over the FEATURE-0015 store.
func NewBackingAccessProvider(store *cloudmodel.Store) *BackingAccessProvider {
	return &BackingAccessProvider{store: store}
}

// Access evaluates caller-relative safe access under one paired backing-read
// lease and releases that lease before returning. Ordinary create/GET/LIST
// paths use this method.
func (p *BackingAccessProvider) Access(
	_ context.Context,
	_ string, // principal retained for the closed provider contract
	grantAction string,
	participationUID, stackUID string,
	grants BackingGrantEvaluator,
) BackingAccessResult {
	if p == nil || p.store == nil || grants == nil {
		return BackingAccessResult{Disposition: BackingSafeDenied}
	}
	var out BackingAccessResult
	p.store.WithPairedBackingRead(participationUID, stackUID, func(snap cloudmodel.PairedBackingSnapshot) {
		out = evaluatePairedBacking(snap, grantAction, grants)
	})
	return out
}

// WithFinalQualificationLease retains a fresh paired lease for the duration of
// fn when access is Allowed. fn must acquire the F0016 lifecycle mutex itself;
// the required order is lease → lifecycle mutex → audit → publication.
// When access is not Allowed, fn is not invoked and the disposition is returned.
func (p *BackingAccessProvider) WithFinalQualificationLease(
	_ context.Context,
	_ string,
	grantAction string,
	participationUID, stackUID string,
	grants BackingGrantEvaluator,
	fn func(BackingView),
) BackingDisposition {
	if p == nil || p.store == nil || grants == nil {
		return BackingSafeDenied
	}
	var disposition BackingDisposition = BackingSafeDenied
	p.store.WithPairedBackingRead(participationUID, stackUID, func(snap cloudmodel.PairedBackingSnapshot) {
		res := evaluatePairedBacking(snap, grantAction, grants)
		disposition = res.Disposition
		if disposition != BackingAllowed || fn == nil {
			return
		}
		fn(res.View)
	})
	return disposition
}

func evaluatePairedBacking(
	snap cloudmodel.PairedBackingSnapshot,
	grantAction string,
	grants BackingGrantEvaluator,
) BackingAccessResult {
	if !snap.ParticipationOK || !snap.StackOK {
		return BackingAccessResult{Disposition: BackingSafeDenied}
	}
	cloudProviderUID := snap.Participation.Spec.CloudProviderRef.UID
	if cloudProviderUID == "" {
		return BackingAccessResult{Disposition: BackingSafeDenied}
	}
	scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: cloudProviderUID}
	if !grants.AuthorizeExact(grantAction, scope, "", "") {
		return BackingAccessResult{Disposition: BackingAuthorizationDenied}
	}
	view := BackingView{
		ParticipationUID:              snap.Participation.Metadata.UID,
		ParticipationEffectiveActive:  participationEffectiveActive(snap.Participation),
		ParticipationCloudProviderUID: cloudProviderUID,
		StackUID:                      snap.Stack.Metadata.UID,
		StackPhase:                    string(snap.Stack.Status.Phase),
		StackScopeUID:                 apimeta.CanonicalScopeIdentity(snap.Stack.Metadata.ScopeRef).UID,
		StackGeneration:               snap.Stack.Metadata.Generation,
	}
	return BackingAccessResult{Disposition: BackingAllowed, View: view}
}

func participationEffectiveActive(part cmmodel.CloudProviderParticipation) bool {
	return part.Status.Phase == cmmodel.ParticipationPhaseActive &&
		!part.Status.PlatformSuspended &&
		!part.Status.ProviderSuspended
}
