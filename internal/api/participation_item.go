package api

import (
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
)

// ParticipationItemHandler serves CloudProviderParticipation GET only.
// Lifecycle mutations are registered exclusively on explicit action routes
// (Task 15); PATCH is not registered and returns 405 from this handler.
type ParticipationItemHandler struct {
	Store       *cloudmodel.Store
	Grants      GrantResolver
	Publication *cloudmodel.PublicationCoordinator
}

// NewParticipationItemHandler constructs the item handler.
func NewParticipationItemHandler(
	store *cloudmodel.Store,
	grants GrantResolver,
	pub *cloudmodel.PublicationCoordinator,
) *ParticipationItemHandler {
	return &ParticipationItemHandler{Store: store, Grants: grants, Publication: pub}
}

// ServeHTTP dispatches GET only. PATCH/PUT/DELETE/HEAD → 405.
func (h *ParticipationItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	switch r.Method {
	case http.MethodGet:
		h.get(w, r, uid)
	default:
		writeMethodNotAllowed(w, r, "GET")
	}
}

func (h *ParticipationItemHandler) get(w http.ResponseWriter, r *http.Request, uid string) {
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
	part, ok := h.Store.GetParticipation(uid)
	if uid == "" || !ok || !participationReadable(h.Grants, part.Metadata.ScopeRef, uid) {
		writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
		return
	}
	writeJSONSuccess(w, r, http.StatusOK, part)
}
