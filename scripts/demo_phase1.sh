#!/opt/homebrew/bin/bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"

echo "Checking health..."
curl -fsS "$BASE_URL/healthz"
echo
curl -fsS "$BASE_URL/readyz"
echo

echo "Creating Organization..."
curl -fsS -X POST "$BASE_URL/v1/organizations" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"Organization",
    "metadata":{"name":"nic","displayName":"National Informatics Centre"},
    "spec":{"description":"Central government cloud organization","sovereignLocations":["in-delhi-1","in-mumbai-1"]}
  }'
echo

echo "Creating OrganizationUnit..."
curl -fsS -X POST "$BASE_URL/v1/organization-units" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"OrganizationUnit",
    "metadata":{"name":"ministry-health","displayName":"Ministry of Health"},
    "spec":{"organizationRef":"nic","description":"Health ministry OU"}
  }'
echo

echo "Creating Tenant..."
curl -fsS -X POST "$BASE_URL/v1/tenants" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"Tenant",
    "metadata":{"name":"national-health-mission","displayName":"National Health Mission"},
    "spec":{"organizationRef":"nic","organizationUnitRef":"ministry-health","isolationProfile":"namespace"}
  }'
echo

echo "Creating Project..."
curl -fsS -X POST "$BASE_URL/v1/projects" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"Project",
    "metadata":{"name":"production","displayName":"Production"},
    "spec":{"organizationRef":"nic","organizationUnitRef":"ministry-health","tenantRef":"national-health-mission","environmentType":"production"}
  }'
echo

echo "Registering ServiceClass..."
curl -fsS -X POST "$BASE_URL/v1/service-classes" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"ServiceClass",
    "metadata":{"name":"datastore.postgresql","displayName":"PostgreSQL"},
    "spec":{"category":"datastore","description":"Managed PostgreSQL datastore","requiredCapabilities":["Provision","Bind","Observe","Delete"]}
  }'
echo

echo "Registering ServicePlan..."
curl -fsS -X POST "$BASE_URL/v1/service-plans" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"ServicePlan",
    "metadata":{"name":"postgres-small-ha"},
    "spec":{"serviceClassRef":"datastore.postgresql","description":"Small HA PostgreSQL plan","tier":"small","highAvailability":true}
  }'
echo

echo "Registering Plugin..."
curl -fsS -X POST "$BASE_URL/v1/plugins" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"Plugin",
    "metadata":{"name":"postgres.dstoreops.basic"},
    "spec":{"pluginType":"dStoreOps","version":"0.1.0","serviceClassRefs":["datastore.postgresql"],"deploymentMode":"compiled-in"}
  }'
echo

echo "Registering Capability..."
curl -fsS -X POST "$BASE_URL/v1/capabilities" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"Capability",
    "metadata":{"name":"postgres-basic-provision"},
    "spec":{"pluginRef":"postgres.dstoreops.basic","serviceClassRef":"datastore.postgresql","operation":"Provision","supported":true}
  }'
echo

echo "Creating ServiceInstance..."
curl -fsS -X POST "$BASE_URL/v1/service-instances" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"ServiceInstance",
    "metadata":{"name":"nhm-prod-postgres"},
    "spec":{"organizationRef":"nic","organizationUnitRef":"ministry-health","tenantRef":"national-health-mission","projectRef":"production","serviceClassRef":"datastore.postgresql","servicePlanRef":"postgres-small-ha","parameters":{"databaseName":"nhm"}}
  }'
echo

echo "Creating ServiceBinding..."
curl -fsS -X POST "$BASE_URL/v1/service-bindings" \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"platform.sovrunn.io/v1alpha1",
    "kind":"ServiceBinding",
    "metadata":{"name":"nhm-app-postgres-binding"},
    "spec":{"serviceInstanceRef":"nhm-prod-postgres","consumerRef":{"kind":"Application","name":"nhm-app"},"bindingType":"credentials"}
  }'
echo

echo "Listing Operations..."
curl -fsS "$BASE_URL/v1/operations"
echo
