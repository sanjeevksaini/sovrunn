package validation

import (
	"context"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func BenchmarkValidateOrganization(b *testing.B) {
	ctx := context.Background()

	org := resources.Organization{
		APIVersion: resources.OrgAPIVersion,
		Kind:       resources.OrgKind,
		Metadata: resources.Metadata{
			Name:        "demo-org",
			DisplayName: "Demo Organization",
		},
		Spec: resources.OrganizationSpec{
			Description:          "Demo organization",
			SovereignLocations:   []string{"in-gurugram"},
			DefaultPolicyProfile: "standard",
		},
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		errs := ValidateOrganization(ctx, org)
		if len(errs) != 0 {
			b.Fatalf("unexpected validation errors: %v", errs)
		}
	}
}

func BenchmarkValidateNamePath(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		errs := ValidateNamePath(ctx, "demo-org")
		if len(errs) != 0 {
			b.Fatalf("unexpected validation errors: %v", errs)
		}
	}
}
