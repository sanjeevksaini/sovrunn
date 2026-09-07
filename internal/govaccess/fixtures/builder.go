package fixtures

// ResourceKind is the closed FEATURE-0018 owned fixture kind set.
type ResourceKind string

const (
	KindAccessGroup             ResourceKind = "AccessGroup"
	KindMembership              ResourceKind = "Membership"
	KindRoleDefinition          ResourceKind = "RoleDefinition"
	KindRoleAssignment          ResourceKind = "RoleAssignment"
	KindPrivilegedAccessRequest ResourceKind = "PrivilegedAccessRequest"
	KindAccessReview            ResourceKind = "AccessReview"
	KindApprovalPolicy          ResourceKind = "ApprovalPolicy"
	KindApprovalRequest         ResourceKind = "ApprovalRequest"
	KindExceptionGrant          ResourceKind = "ExceptionGrant"
	KindGovernanceProfileV1     ResourceKind = "GovernanceProfile.v1"
)

// Resource is the synthetic deterministic fixture graph node.
type Resource struct {
	Kind    ResourceKind
	UID     string
	RefUIDs []string
}

// FixtureSet is an immutable synthetic fixture graph.
type FixtureSet struct {
	Resources          []Resource
	ProhibitedConcepts []ProhibitedConcept
}

// Builder is an immutable fixture builder.
type Builder struct {
	set FixtureSet
}

// NewBuilder constructs an empty fixture builder.
func NewBuilder() Builder {
	return Builder{set: FixtureSet{
		Resources:          []Resource{},
		ProhibitedConcepts: []ProhibitedConcept{},
	}}
}

// WithResource appends a deep-copied synthetic resource and returns a new builder.
func (b Builder) WithResource(r Resource) Builder {
	next := copyFixtureSet(b.set)
	next.Resources = append(next.Resources, copyResource(r))
	return Builder{set: next}
}

// WithProhibitedConcept appends a deep-copied prohibited concept and returns a new builder.
func (b Builder) WithProhibitedConcept(c ProhibitedConcept) Builder {
	next := copyFixtureSet(b.set)
	next.ProhibitedConcepts = append(next.ProhibitedConcepts, c)
	return Builder{set: next}
}

// Build returns a deep-copied immutable fixture set.
func (b Builder) Build() FixtureSet {
	return copyFixtureSet(b.set)
}

func copyFixtureSet(in FixtureSet) FixtureSet {
	out := FixtureSet{
		Resources:          make([]Resource, 0, len(in.Resources)),
		ProhibitedConcepts: append([]ProhibitedConcept(nil), in.ProhibitedConcepts...),
	}
	for _, r := range in.Resources {
		out.Resources = append(out.Resources, copyResource(r))
	}
	return out
}

func copyResource(in Resource) Resource {
	out := in
	out.RefUIDs = append([]string(nil), in.RefUIDs...)
	return out
}
