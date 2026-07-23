package apiconform

import (
	"reflect"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

// TypeBindings is the concrete schema-to-Go TypeBinding registry for
// FEATURE-0012 (D-01b, F12-NAMING-005, F12-VERIFY-001(13)).
//
// Bindings live in apiconform so apischema stays free of concrete contract
// types and never imports this package (no cycle). Each entry maps one
// canonical schema or _common sub-schema to its derivative Go type.
// VerifyGoTypeAgainstSchema is the authoritative consistency check;
// fixture round-tripping is supporting evidence only.
var TypeBindings = []apischema.TypeBinding{
	// _common sub-schemas
	{SchemaPath: "api/schemas/_common/type-meta.json", GoType: reflect.TypeOf(apimeta.TypeMeta{})},
	{SchemaPath: "api/schemas/_common/object-meta.json", GoType: reflect.TypeOf(apimeta.ObjectMeta{})},
	{SchemaPath: "api/schemas/_common/typed-ref.json", GoType: reflect.TypeOf(apimeta.TypedRef{})},
	{SchemaPath: "api/schemas/_common/scope-ref.json", GoType: reflect.TypeOf(apimeta.ScopeRef{})},
	{SchemaPath: "api/schemas/_common/owner-ref.json", GoType: reflect.TypeOf(apimeta.OwnerRef{})},
	{SchemaPath: "api/schemas/_common/condition.json", GoType: reflect.TypeOf(apicond.Condition{})},
	{SchemaPath: "api/schemas/_common/problem.json", GoType: reflect.TypeOf(apiproblem.Problem{})},
	{SchemaPath: "api/schemas/_common/violation.json", GoType: reflect.TypeOf(apiproblem.Violation{})},
	{SchemaPath: "api/schemas/_common/page.json", GoType: reflect.TypeOf(apimeta.Page{})},

	// Eight canonical contract schemas (Matrix D)
	{SchemaPath: "api/schemas/project.json", GoType: reflect.TypeOf(Project{})},
	{SchemaPath: "api/schemas/resource-pool.json", GoType: reflect.TypeOf(ResourcePool{})},
	{SchemaPath: "api/schemas/discovered-database.json", GoType: reflect.TypeOf(DiscoveredDatabase{})},
	{SchemaPath: "api/schemas/plugin-definition.json", GoType: reflect.TypeOf(PluginDefinition{})},
	{SchemaPath: "api/schemas/adapter-configuration.json", GoType: reflect.TypeOf(AdapterConfiguration{})},
	{SchemaPath: "api/schemas/placement-evaluation-request.json", GoType: reflect.TypeOf(PlacementEvaluationRequest{})},
	{SchemaPath: "api/schemas/operation.json", GoType: reflect.TypeOf(Operation{})},
	{SchemaPath: "api/schemas/audit-event.json", GoType: reflect.TypeOf(AuditEvent{})},
}
