package registry

import (
	"context"
	"fmt"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func benchmarkOrganization(name string) resources.Organization {
	return resources.Organization{
		APIVersion: resources.OrgAPIVersion,
		Kind:       resources.OrgKind,
		Metadata: resources.Metadata{
			Name:        name,
			DisplayName: "Demo Organization",
			Labels: map[string]string{
				"env": "benchmark",
			},
			Annotations: map[string]string{
				"owner": "phase1-performance",
			},
		},
		Spec: resources.OrganizationSpec{
			Description:          "Benchmark organization",
			SovereignLocations:   []string{"in-gurugram"},
			DefaultPolicyProfile: "standard",
		},
		Status: resources.OrganizationStatus{
			Phase: resources.PhaseActive,
		},
	}
}

func BenchmarkOrganizationRegistryCreate(b *testing.B) {
	ctx := context.Background()
	reg := NewOrganizationRegistry()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		org := benchmarkOrganization(fmt.Sprintf("demo-org-%d", i))
		if err := reg.CreateOrganization(ctx, org); err != nil {
			b.Fatalf("CreateOrganization() error = %v", err)
		}
	}
}

func BenchmarkOrganizationRegistryGet(b *testing.B) {
	ctx := context.Background()
	reg := NewOrganizationRegistry()

	org := benchmarkOrganization("demo-org")
	if err := reg.CreateOrganization(ctx, org); err != nil {
		b.Fatalf("CreateOrganization() error = %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reg.GetOrganization(ctx, "demo-org")
		if err != nil {
			b.Fatalf("GetOrganization() error = %v", err)
		}
	}
}

func BenchmarkOrganizationRegistryList100(b *testing.B) {
	ctx := context.Background()
	reg := NewOrganizationRegistry()

	for i := 0; i < 100; i++ {
		org := benchmarkOrganization(fmt.Sprintf("demo-org-%03d", i))
		if err := reg.CreateOrganization(ctx, org); err != nil {
			b.Fatalf("CreateOrganization() error = %v", err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		items, err := reg.ListOrganizations(ctx)
		if err != nil {
			b.Fatalf("ListOrganizations() error = %v", err)
		}
		if len(items) != 100 {
			b.Fatalf("ListOrganizations() len = %d, want 100", len(items))
		}
	}
}

func BenchmarkOrganizationRegistryList1000(b *testing.B) {
	ctx := context.Background()
	reg := NewOrganizationRegistry()

	for i := 0; i < 1000; i++ {
		org := benchmarkOrganization(fmt.Sprintf("demo-org-%04d", i))
		if err := reg.CreateOrganization(ctx, org); err != nil {
			b.Fatalf("CreateOrganization() error = %v", err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		items, err := reg.ListOrganizations(ctx)
		if err != nil {
			b.Fatalf("ListOrganizations() error = %v", err)
		}
		if len(items) != 1000 {
			b.Fatalf("ListOrganizations() len = %d, want 1000", len(items))
		}
	}
}
