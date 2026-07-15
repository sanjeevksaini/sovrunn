package registry

import (
	"context"
	"errors"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Sentinel errors — handlers map these to APIError codes without
// inspecting error message strings.
var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
)

// OrganizationRegistryIface is the storage contract for Organization
// resources. Implementations may be in-memory (Phase 1) or durable
// (future phases). The interface keeps handlers decoupled from storage.
type OrganizationRegistryIface interface {
	CreateOrganization(ctx context.Context, org resources.Organization) error
	GetOrganization(ctx context.Context, name string) (resources.Organization, error)
	ListOrganizations(ctx context.Context) ([]resources.Organization, error)
	UpdateOrganization(ctx context.Context, name string, org resources.Organization) (resources.Organization, error)
	DeleteOrganization(ctx context.Context, name string) error
}
