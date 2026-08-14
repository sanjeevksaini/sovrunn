package api

import (
	"net/http"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

// FaultDomainItemHandler serves FaultDomain GET and PATCH.
type FaultDomainItemHandler struct {
	Store       *cloudmodel.Store
	Grants      GrantResolver
	Publication *cloudmodel.PublicationCoordinator
}

// NewFaultDomainItemHandler constructs the item handler.
func NewFaultDomainItemHandler(
	store *cloudmodel.Store,
	grants GrantResolver,
	pub *cloudmodel.PublicationCoordinator,
) *FaultDomainItemHandler {
	return &FaultDomainItemHandler{Store: store, Grants: grants, Publication: pub}
}

// ServeHTTP dispatches GET and PATCH. PUT/DELETE/HEAD → 405.
func (h *FaultDomainItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	switch r.Method {
	case http.MethodGet:
		h.get(w, r, uid)
	case http.MethodPatch:
		h.patch(w, r, uid)
	default:
		writeMethodNotAllowed(w, r, "GET, PATCH")
	}
}

func (h *FaultDomainItemHandler) get(w http.ResponseWriter, r *http.Request, uid string) {
	principal, prob := authenticatePrincipal(r, h.Grants)
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if cloudmodel.HasForgedGrantHeader(r.Header) {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.NewForgedGrantDenial(
			actorRef(principal), requestID(r), newAuditUID("forged"),
		))
		return
	}
	fd, ok := h.Store.GetFaultDomain(uid)
	if uid == "" || !ok || !topologyReadable(h.Grants, fd.Metadata.ScopeRef, uid) {
		writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
		return
	}
	writeJSONSuccess(w, r, http.StatusOK, fd)
}

func (h *FaultDomainItemHandler) patch(w http.ResponseWriter, r *http.Request, uid string) {
	principal, prob := authenticatePrincipal(r, h.Grants)
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if cloudmodel.HasForgedGrantHeader(r.Header) {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.NewForgedGrantDenial(
			actorRef(principal), requestID(r), newAuditUID("forged"),
		))
		return
	}
	ifMatch := r.Header.Get("If-Match")
	if prob := validate.ValidateRequiredIfMatch(ifMatch); prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if !mergePatchMediaTypeOK(r) {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeUnsupportedMediaType).WithDetail("PATCH requires application/merge-patch+json"))
		return
	}

	phaseOne, prob := DecodePhaseOneHTTP(w, r, apivalid.DefaultLimits())
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if phaseOne.HasBootstrapGrant {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.NewForgedGrantDenial(
			actorRef(principal), requestID(r), newAuditUID("forged-body"),
		))
		return
	}

	current, ok := h.Store.GetFaultDomain(uid)
	if !ok || !topologyReadable(h.Grants, current.Metadata.ScopeRef, uid) {
		writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
		return
	}

	if len(h.Grants.CoarseLookup(cloudmodel.ActionTopologyWrite)) == 0 ||
		!topologyWritable(h.Grants, current.Metadata.ScopeRef, uid) {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("cloud-provider-admin topology.write grant is required"))
		return
	}

	if prob := DecodePhaseTwoStrict(validate.KindFaultDomain, validate.RequestClassPatch, phaseOne.Body, nil); prob != nil {
		writeClassifiedCreateDenial(w, r, h.Publication, principal, prob)
		return
	}
	mergedAny, prob := ApplyMergePatch(validate.KindFaultDomain, current, phaseOne.Body)
	if prob != nil {
		writeClassifiedCreateDenial(w, r, h.Publication, principal, prob)
		return
	}
	merged, ok := mergedAny.(model.FaultDomain)
	if !ok {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("merge produced unexpected type"))
		return
	}
	if prob := runFeature0015StageSet(r.Context(), validate.KindFaultDomain, validate.RequestClassPatch, merged); prob != nil {
		writeProblem(w, r, prob)
		return
	}

	h.Store.BeginPublication()
	defer h.Store.EndPublication()
	live, ok := h.Store.LookupFaultDomain(uid)
	if !ok {
		writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
		return
	}
	if prob := validate.CheckIfMatchCurrent(ifMatch, live.Metadata.ResourceVersion); prob != nil {
		writeProblem(w, r, prob)
		return
	}
	merged.Metadata = live.Metadata
	merged.Metadata.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	merged.Status = live.Status
	staged, stageProb := h.Store.StageUpdateFaultDomain(merged)
	if stageProb != nil {
		writeProblem(w, r, stageProb)
		return
	}
	event := cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
		UID: newAuditUID("patch"), RequestID: requestID(r),
		Action: cloudmodel.AuditActionResourcePatch, Outcome: apiconform.AuditOutcomeSucceeded,
		Reason: "faultdomain.patched", Actor: actorRef(principal),
		Subject: apimeta.TypedRef{
			APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain,
			Name: live.Metadata.Name, UID: uid,
		},
		SubjectResourceVersion: live.Metadata.ResourceVersion,
	})
	if pubProb := h.Publication.PublishSuccess(r.Context(), cloudmodel.AuditedSuccess{
		Staged: staged, Event: event,
	}); pubProb != nil {
		writeProblem(w, r, pubProb)
		return
	}
	updated, _ := h.Store.LookupFaultDomain(uid)
	writeJSONSuccess(w, r, http.StatusOK, updated)
}
