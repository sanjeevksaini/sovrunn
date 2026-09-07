package govaccess

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
)

// ComposeFeature0013BundleView performs the sole FEATURE-0018 root composition
// of FEATURE-0013 profiles by loading evidence-owned strict local bytes through
// decision/bundle.Load and returning the read-only BundleView.
func ComposeFeature0013BundleView() (bundle.BundleView, *apiproblem.Problem) {
	return composeBundleView(evidence.LocalDecisionProfileBundleBytes())
}

func composeBundleView(profileBundleBytes []byte) (bundle.BundleView, *apiproblem.Problem) {
	loaded, prob := bundle.Load(profileBundleBytes, apivalid.ModeReadRepresentation)
	if prob != nil {
		return bundle.BundleView{}, prob
	}
	return loaded.View(), nil
}
