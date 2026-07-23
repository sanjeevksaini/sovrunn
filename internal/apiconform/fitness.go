package apiconform

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

// Fitness check ID for the D-12 / F12-LEDGER-001 boundary-ledger fitness
// function. This is not one of the numbered F12-VERIFY-001 checks 1–15;
// those are aggregated separately via RegisteredFitnessChecks.
const FitnessCheckBoundaryLedger = "ledger"

// Stable fitness finding codes for boundary-ledger aggregation (task 16.5).
const (
	CodeFitnessCheckRunMissing             = "FITNESS_CHECK_RUN_MISSING"
	CodeFitnessLedgerLoadFailed            = "FITNESS_LEDGER_LOAD_FAILED"
	CodeFitnessLedgerParseFailed           = "FITNESS_LEDGER_PARSE_FAILED"
	CodeFitnessLedgerCategoryMissing       = "FITNESS_LEDGER_CATEGORY_MISSING"
	CodeFitnessLedgerInvalidBoundary       = "FITNESS_LEDGER_INVALID_BOUNDARY"
	CodeFitnessLedgerDuplicateBoundary     = "FITNESS_LEDGER_DUPLICATE_BOUNDARY"
	CodeFitnessLedgerSchemaBoundaryMissing = "FITNESS_LEDGER_SCHEMA_BOUNDARY_MISSING"
)

// FitnessCheckRegistration is one executable F12-VERIFY-001 fitness check
// registered for aggregation (architecture §8.2; task 16.5).
type FitnessCheckRegistration struct {
	// ID is the stable check identifier ("1".."15", plus "1a" for field-policy).
	ID string
	// Requirement is the F12-VERIFY-001 / architecture §8.2 requirement text.
	Requirement string
	// Run executes the check against the repository module root and returns
	// fail-closed findings (empty slice means pass).
	Run func(moduleRoot string) []FitnessFinding
}

// fitnessCheckRequirementMap is the authoritative F12-VERIFY-001 check →
// requirement enumeration asserted by aggregation tests (design: fitness.go).
var fitnessCheckRequirementMap = map[string]string{
	FitnessCheckExternalSchemaAnnotations:         "every external schema declares profile, boundary, stability, and allowed scopes",
	FitnessCheckFieldPolicyCoverage:               "every API-boundary property declares complete explicit x-sovrunn-field-policy",
	FitnessCheckNoProviderSDKInCoreCustomer:       "no core/customer schema imports or embeds provider SDK/native types",
	FitnessCheckMutableFieldOwnership:             "every mutable field and condition has one owner",
	FitnessCheckUnknownAndDuplicateFields:         "unknown and duplicate fields fail",
	FitnessCheckReferencesConstrainKindsAndScopes: "references constrain kinds and scopes",
	FitnessCheckCrossTenantNoExistenceDisclosure:  "cross-tenant access fails without existence disclosure",
	FitnessCheckNoRawSecretLikeValues:             "raw secret-like values are prohibited from metadata/status/errors",
	FitnessCheckObservationProvenanceAndFreshness: "externally sourced observations include provenance and freshness",
	FitnessCheckPublishedDefinitionsImmutable:     "published definitions are immutable",
	FitnessCheckSchemaCompatibility:               "schema compatibility detects breaking changes",
	FitnessCheckSizesBounded:                      "object, metadata, condition, violation, reference, and page sizes are bounded",
	FitnessCheckStableCodesAndJSONPointers:        "errors use stable codes and JSON Pointer paths",
	FitnessCheckGeneratedArtifactsMatchSchema:     "generated artifacts match the canonical schema",
	FitnessCheckLaterFeatureRuntimeAbsent:         "later-feature runtime behavior is absent",
	FitnessCheckExceptionsRequireApprovedHandoff:  "exceptions require an approved architecture handoff",
}

// FitnessCheckRequirementMap returns a copy of the F12-VERIFY-001 check →
// requirement enumeration owned by this package.
func FitnessCheckRequirementMap() map[string]string {
	out := make(map[string]string, len(fitnessCheckRequirementMap))
	for k, v := range fitnessCheckRequirementMap {
		out[k] = v
	}
	return out
}

// RequiredFitnessCheckIDs returns the ordered F12-VERIFY-001 check IDs 1–15.
// Check "1a" (field-policy coverage) is registered in addition to these.
func RequiredFitnessCheckIDs() []string {
	return []string{
		FitnessCheckExternalSchemaAnnotations,         // 1
		FitnessCheckNoProviderSDKInCoreCustomer,       // 2
		FitnessCheckMutableFieldOwnership,             // 3
		FitnessCheckUnknownAndDuplicateFields,         // 4
		FitnessCheckReferencesConstrainKindsAndScopes, // 5
		FitnessCheckCrossTenantNoExistenceDisclosure,  // 6
		FitnessCheckNoRawSecretLikeValues,             // 7
		FitnessCheckObservationProvenanceAndFreshness, // 8
		FitnessCheckPublishedDefinitionsImmutable,     // 9
		FitnessCheckSchemaCompatibility,               // 10
		FitnessCheckSizesBounded,                      // 11
		FitnessCheckStableCodesAndJSONPointers,        // 12
		FitnessCheckGeneratedArtifactsMatchSchema,     // 13
		FitnessCheckLaterFeatureRuntimeAbsent,         // 14
		FitnessCheckExceptionsRequireApprovedHandoff,  // 15
	}
}

// RegisteredFitnessChecks returns the ordered registry of executable fitness
// checks. Includes F12-VERIFY-001 checks 1–15 plus check 1a (field-policy).
func RegisteredFitnessChecks() []FitnessCheckRegistration {
	schemasRoot := func(moduleRoot string) string {
		return filepath.Join(moduleRoot, CanonicalSchemasDir)
	}
	return []FitnessCheckRegistration{
		{
			ID:          FitnessCheckExternalSchemaAnnotations,
			Requirement: fitnessCheckRequirementMap[FitnessCheckExternalSchemaAnnotations],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckExternalSchemaAnnotations(schemasRoot(moduleRoot))
			},
		},
		{
			ID:          FitnessCheckFieldPolicyCoverage,
			Requirement: fitnessCheckRequirementMap[FitnessCheckFieldPolicyCoverage],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckFieldPolicyCoverage(schemasRoot(moduleRoot))
			},
		},
		{
			ID:          FitnessCheckNoProviderSDKInCoreCustomer,
			Requirement: fitnessCheckRequirementMap[FitnessCheckNoProviderSDKInCoreCustomer],
			Run:         CheckNoProviderSDKInCoreCustomer,
		},
		{
			ID:          FitnessCheckMutableFieldOwnership,
			Requirement: fitnessCheckRequirementMap[FitnessCheckMutableFieldOwnership],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckMutableFieldOwnership(schemasRoot(moduleRoot))
			},
		},
		{
			ID:          FitnessCheckUnknownAndDuplicateFields,
			Requirement: fitnessCheckRequirementMap[FitnessCheckUnknownAndDuplicateFields],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckUnknownAndDuplicateFieldsFail()
			},
		},
		{
			ID:          FitnessCheckReferencesConstrainKindsAndScopes,
			Requirement: fitnessCheckRequirementMap[FitnessCheckReferencesConstrainKindsAndScopes],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckReferencesConstrainKindsAndScopes(schemasRoot(moduleRoot))
			},
		},
		{
			ID:          FitnessCheckCrossTenantNoExistenceDisclosure,
			Requirement: fitnessCheckRequirementMap[FitnessCheckCrossTenantNoExistenceDisclosure],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckCrossTenantAccessNoExistenceDisclosure()
			},
		},
		{
			ID:          FitnessCheckNoRawSecretLikeValues,
			Requirement: fitnessCheckRequirementMap[FitnessCheckNoRawSecretLikeValues],
			Run:         CheckNoRawSecretLikeValues,
		},
		{
			ID:          FitnessCheckObservationProvenanceAndFreshness,
			Requirement: fitnessCheckRequirementMap[FitnessCheckObservationProvenanceAndFreshness],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckObservationProvenanceAndFreshness(schemasRoot(moduleRoot))
			},
		},
		{
			ID:          FitnessCheckPublishedDefinitionsImmutable,
			Requirement: fitnessCheckRequirementMap[FitnessCheckPublishedDefinitionsImmutable],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckPublishedDefinitionsImmutable(schemasRoot(moduleRoot))
			},
		},
		{
			ID:          FitnessCheckSchemaCompatibility,
			Requirement: fitnessCheckRequirementMap[FitnessCheckSchemaCompatibility],
			Run:         CheckSchemaCompatibilityDetectsBreaking,
		},
		{
			ID:          FitnessCheckSizesBounded,
			Requirement: fitnessCheckRequirementMap[FitnessCheckSizesBounded],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckSizesBounded()
			},
		},
		{
			ID:          FitnessCheckStableCodesAndJSONPointers,
			Requirement: fitnessCheckRequirementMap[FitnessCheckStableCodesAndJSONPointers],
			Run: func(moduleRoot string) []FitnessFinding {
				return CheckErrorsUseStableCodesAndJSONPointers()
			},
		},
		{
			ID:          FitnessCheckGeneratedArtifactsMatchSchema,
			Requirement: fitnessCheckRequirementMap[FitnessCheckGeneratedArtifactsMatchSchema],
			Run:         CheckGeneratedArtifactsMatchCanonicalSchema,
		},
		{
			ID:          FitnessCheckLaterFeatureRuntimeAbsent,
			Requirement: fitnessCheckRequirementMap[FitnessCheckLaterFeatureRuntimeAbsent],
			Run:         CheckLaterFeatureRuntimeAbsent,
		},
		{
			ID:          FitnessCheckExceptionsRequireApprovedHandoff,
			Requirement: fitnessCheckRequirementMap[FitnessCheckExceptionsRequireApprovedHandoff],
			Run:         CheckExceptionsRequireApprovedHandoff,
		},
	}
}

// RunAllFitnessChecks executes every registered fitness check (1–15 plus 1a)
// against moduleRoot and returns the concatenated findings. An empty result
// means all registered checks passed.
func RunAllFitnessChecks(moduleRoot string) []FitnessFinding {
	regs := RegisteredFitnessChecks()
	var findings []FitnessFinding
	for _, reg := range regs {
		if reg.Run == nil {
			findings = append(findings, FitnessFinding{
				Check:   reg.ID,
				Schema:  "fitness-registry",
				Path:    "/" + reg.ID,
				Code:    CodeFitnessCheckRunMissing,
				Message: fmt.Sprintf("fitness check %q has nil Run function", reg.ID),
			})
			continue
		}
		findings = append(findings, reg.Run(moduleRoot)...)
	}
	return findings
}

// CheckBoundaryLedger implements the D-12 / F12-LEDGER-001 fitness function:
// strictly parse docs/api/boundary-ledger.yaml; require every declared
// boundary to carry all F12-LEDGER-001 categories; require every boundary
// declared on a canonical schema to have a ledger entry.
func CheckBoundaryLedger(moduleRoot string) []FitnessFinding {
	ledgerPath := filepath.Join(moduleRoot, BoundaryLedgerPath)
	schemasRoot := filepath.Join(moduleRoot, CanonicalSchemasDir)
	return checkBoundaryLedgerAt(ledgerPath, schemasRoot)
}

// checkBoundaryLedgerAt is the path-injectable implementation used by
// CheckBoundaryLedger and negative fitness tests.
func checkBoundaryLedgerAt(ledgerPath, schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding

	raw, err := readRepoFile(ledgerPath)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckBoundaryLedger,
			Schema:  BoundaryLedgerPath,
			Path:    "/",
			Code:    CodeFitnessLedgerLoadFailed,
			Message: err.Error(),
		}}
	}

	doc, err := ParseBoundaryLedgerYAML(raw)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckBoundaryLedger,
			Schema:  BoundaryLedgerPath,
			Path:    "/",
			Code:    CodeFitnessLedgerParseFailed,
			Message: err.Error(),
		}}
	}

	seen := make(map[string]int, len(doc.Boundaries))
	for i, entry := range doc.Boundaries {
		id := strings.TrimSpace(entry.ID)
		path := fmt.Sprintf("/boundaries/%d", i)
		if id == "" || !apimeta.Boundary(id).Valid() {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckBoundaryLedger,
				Schema:  BoundaryLedgerPath,
				Path:    path + "/id",
				Code:    CodeFitnessLedgerInvalidBoundary,
				Message: fmt.Sprintf("ledger boundary id %q is not a Matrix C1 boundary", entry.ID),
			})
			continue
		}
		seen[id]++
		if seen[id] > 1 {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckBoundaryLedger,
				Schema:  BoundaryLedgerPath,
				Path:    path + "/id",
				Code:    CodeFitnessLedgerDuplicateBoundary,
				Message: fmt.Sprintf("duplicate ledger entry for boundary %q", id),
			})
		}
		for _, cat := range LedgerEntryCategoryGaps(entry) {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckBoundaryLedger,
				Schema:  BoundaryLedgerPath,
				Path:    path + "/" + cat,
				Code:    CodeFitnessLedgerCategoryMissing,
				Message: fmt.Sprintf("F12-LEDGER-001 category %q is missing or empty on boundary %q", cat, id),
			})
		}
	}

	schemaBoundaries, schemaFindings := collectCanonicalSchemaBoundaries(schemasRoot)
	findings = append(findings, schemaFindings...)
	for _, boundary := range sortedKeys(schemaBoundaries) {
		if seen[boundary] == 0 {
			schemas := schemaBoundaries[boundary]
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckBoundaryLedger,
				Schema:  BoundaryLedgerPath,
				Path:    "/boundaries",
				Code:    CodeFitnessLedgerSchemaBoundaryMissing,
				Message: fmt.Sprintf("canonical schema boundary %q (schemas: %s) has no ledger entry", boundary, strings.Join(schemas, ", ")),
			})
		}
	}

	return findings
}

// LedgerEntryCategoryGaps returns the F12-LEDGER-001 category field names that
// are missing or empty on entry. An empty result means the entry is complete.
func LedgerEntryCategoryGaps(entry BoundaryLedgerEntry) []string {
	checks := map[string]bool{
		"purpose":              strings.TrimSpace(entry.Purpose) != "",
		"owner":                strings.TrimSpace(entry.Owner) != "",
		"producers":            len(entry.Producers) > 0 && allNonEmptyStrings(entry.Producers),
		"consumers":            len(entry.Consumers) > 0 && allNonEmptyStrings(entry.Consumers),
		"allowed_data":         len(entry.AllowedData) > 0 && allNonEmptyStrings(entry.AllowedData),
		"prohibited_data":      len(entry.ProhibitedData) > 0 && allNonEmptyStrings(entry.ProhibitedData),
		"authorization":        strings.TrimSpace(entry.Authorization) != "",
		"audit":                strings.TrimSpace(entry.Audit) != "",
		"observability":        strings.TrimSpace(entry.Observability) != "",
		"failure_behavior":     strings.TrimSpace(entry.FailureBehavior) != "",
		"versioning":           strings.TrimSpace(entry.Versioning) != "",
		"replacement_path":     strings.TrimSpace(entry.ReplacementPath) != "",
		"migration_path":       strings.TrimSpace(entry.MigrationPath) != "",
		"reassessment_trigger": strings.TrimSpace(entry.ReassessmentTrigger) != "",
	}
	var missing []string
	for _, cat := range requiredLedgerCategories {
		if !checks[cat] {
			missing = append(missing, cat)
		}
	}
	return missing
}

func collectCanonicalSchemaBoundaries(schemasRoot string) (map[string][]string, []FitnessFinding) {
	out := make(map[string][]string)
	var findings []FitnessFinding
	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := readRepoFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckBoundaryLedger,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessLedgerLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		meta, issues := apischema.ReadAnnotations(body)
		if len(issues) > 0 || !meta.Boundary.Valid() {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckBoundaryLedger,
				Schema:  schemaID,
				Path:    "/" + apischema.ExtBoundary,
				Code:    CodeFitnessLedgerInvalidBoundary,
				Message: fmt.Sprintf("canonical schema %s missing or invalid x-sovrunn-boundary", name),
			})
			continue
		}
		b := string(meta.Boundary)
		out[b] = append(out[b], schemaID)
	}
	return out, findings
}

func allNonEmptyStrings(values []string) bool {
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	return true
}

func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
