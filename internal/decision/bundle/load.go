package bundle

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// Load decodes a declarative decision-profile bundle from local bytes using
// apivalid.StrictDecode (design §8; RID-09 / closure 13). Media type is
// sniffed from leading content (JSON object/array → application/json;
// otherwise application/yaml). After decode it builds the derivative
// (id,version) BundleView, rejects duplicates, and applies offline trust
// fail-closed checks. No network fetch, digest computation, or signature
// verification is performed.
func Load(data []byte, mode apivalid.DecodeMode) (Bundle, *apiproblem.Problem) {
	contentType := sniffContentType(data)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	req.Header.Set("Content-Type", contentType)

	var raw Bundle
	if prob := apivalid.StrictDecode(httptest.NewRecorder(), req, apivalid.DefaultLimits(), mode, &raw); prob != nil {
		return Bundle{}, prob
	}

	if prob := validateBundleTrust(raw.Spec); prob != nil {
		return Bundle{}, prob
	}

	view, prob := buildView(raw.Spec)
	if prob != nil {
		return Bundle{}, prob
	}
	raw.view = view
	return raw, nil
}

// sniffContentType selects StrictDecode media type from leading bytes.
// JSON objects/arrays use application/json; all other local fixtures use
// application/yaml so JSON/YAML decode equivalence is exercised through the
// same StrictDecode path (design §8; RID-09).
func sniffContentType(data []byte) string {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return "application/json"
	}
	return "application/yaml"
}
