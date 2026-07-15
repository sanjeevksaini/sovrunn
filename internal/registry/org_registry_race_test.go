package registry

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestOrganizationRegistry_ConcurrentStress(t *testing.T) {
	reg := NewOrganizationRegistry()
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			name := fmt.Sprintf("org-%d", id%6)
			org := resources.Organization{
				APIVersion: resources.OrgAPIVersion,
				Kind:       resources.OrgKind,
				Metadata:   resources.Metadata{Name: name},
				Status:     resources.OrganizationStatus{Phase: resources.PhaseActive},
			}
			_ = reg.CreateOrganization(ctx, org)
			_, _ = reg.GetOrganization(ctx, name)
			_, _ = reg.ListOrganizations(ctx)
			update := org
			update.Metadata.DisplayName = fmt.Sprintf("display-%d", id)
			_, _ = reg.UpdateOrganization(ctx, name, update)
			_ = reg.DeleteOrganization(ctx, name)
		}(i)
	}
	wg.Wait()
}
