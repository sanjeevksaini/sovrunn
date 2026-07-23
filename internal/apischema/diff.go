package apischema

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

// ChangeClass is the closed set of schema-evolution classifications applied by
// ClassifyChange (D-11, F12-EVOLVE-002).
type ChangeClass string

const (
	// ChangeCompatible is an additive, backward-compatible schema delta.
	ChangeCompatible ChangeClass = "Compatible"
	// ChangeBreaking is a backward-incompatible schema delta.
	ChangeBreaking ChangeClass = "Breaking"
	// ChangeReviewRequired is a delta that needs compatibility or security review.
	ChangeReviewRequired ChangeClass = "ReviewRequired"
)

// ChangeKind identifies the kind of schema delta that produced a classification.
// Values align with the architecture change-classification table (F12-EVOLVE-002).
type ChangeKind string

const (
	KindAddOptionalField          ChangeKind = "add_optional_field"
	KindAddRequiredField          ChangeKind = "add_required_field"
	KindRemoveField               ChangeKind = "remove_field"
	KindPromoteOptionalToRequired ChangeKind = "promote_optional_to_required"
	KindDemoteRequiredToOptional  ChangeKind = "demote_required_to_optional"
	KindChangeFieldMeaning        ChangeKind = "change_field_meaning"
	KindChangeOwnerOrMutability   ChangeKind = "change_owner_or_mutability"
	KindNarrowEnum                ChangeKind = "narrow_enum"
	KindNarrowValidationRange     ChangeKind = "narrow_validation_range"
	KindWidenValidationRange      ChangeKind = "widen_validation_range"
	KindAddEnumValue              ChangeKind = "add_enum_value"
	KindChangeReferenceTarget     ChangeKind = "change_reference_target"
	KindChangeAllowedScopes       ChangeKind = "change_allowed_scopes"
	KindExposeInternalPublicly    ChangeKind = "expose_internal_publicly"
	KindAddRegisteredExtension    ChangeKind = "add_registered_extension"
	KindRemoveRegisteredExtension ChangeKind = "remove_registered_extension"
	KindMalformedSchema           ChangeKind = "malformed_schema"
)

// Change is one classified schema delta produced by ClassifyChange.
type Change struct {
	Class   ChangeClass
	Kind    ChangeKind
	Path    string // RFC 6901 JSON Pointer to the affected location
	Message string
}

// ClassifyChange compares old and new JSON Schema documents and returns one
// Change per detected delta classified as Compatible, Breaking, or
// ReviewRequired per the change-classification table (D-11, F12-EVOLVE-002):
//
//	add optional field                         → Compatible
//	add required field / promote to required   → Breaking
//	remove field                               → Breaking
//	change field meaning / owner / mutability  → Breaking
//	narrow enum or validation range            → Breaking
//	add enum value                             → ReviewRequired
//	change reference target kind/scope         → ReviewRequired
//	expose internal data publicly              → ReviewRequired
//	add registered extension                   → Compatible
//
// Identical schemas return an empty slice. Results are sorted by Path, Kind,
// then Class for deterministic gate output. Baseline integrity and approval
// gates are VerifyBaselineIntegrity and VerifyBaselineApproval (D-11).
func ClassifyChange(oldSchema, newSchema []byte) []Change {
	oldRoot, errOld := parseSchemaObject(oldSchema)
	newRoot, errNew := parseSchemaObject(newSchema)
	if errOld != nil || errNew != nil {
		msg := "schema document is not valid JSON object"
		if errOld != nil && errNew != nil {
			msg = fmt.Sprintf("old and new schemas are not valid JSON objects: old=%v; new=%v", errOld, errNew)
		} else if errOld != nil {
			msg = fmt.Sprintf("old schema is not a valid JSON object: %v", errOld)
		} else {
			msg = fmt.Sprintf("new schema is not a valid JSON object: %v", errNew)
		}
		return []Change{{
			Class:   ChangeReviewRequired,
			Kind:    KindMalformedSchema,
			Path:    "/",
			Message: msg,
		}}
	}

	var changes []Change
	diffSchemaNode("", oldRoot, newRoot, &changes)
	sortChanges(changes)
	if changes == nil {
		return []Change{}
	}
	return changes
}

func parseSchemaObject(raw []byte) (map[string]any, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty schema document")
	}
	var root any
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, err
	}
	obj, ok := root.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema root must be a JSON object")
	}
	return obj, nil
}

func diffSchemaNode(path string, oldNode, newNode map[string]any, out *[]Change) {
	diffProperties(path, oldNode, newNode, out)
	diffRequired(path, oldNode, newNode, out)
	diffEnum(path, oldNode, newNode, out)
	diffValidationRange(path, oldNode, newNode, out)
	diffType(path, oldNode, newNode, out)
	diffRef(path, oldNode, newNode, out)
	diffAllowedScopes(path, oldNode, newNode, out)
	diffFieldPolicy(path, oldNode, newNode, out)
	diffRegisteredExtensions(path, oldNode, newNode, out)
	diffItems(path, oldNode, newNode, out)
	diffAdditionalProperties(path, oldNode, newNode, out)
}

func diffProperties(path string, oldNode, newNode map[string]any, out *[]Change) {
	oldProps := propertiesMap(oldNode)
	newProps := propertiesMap(newNode)
	oldRequired := requiredSet(oldNode)
	newRequired := requiredSet(newNode)

	allNames := make(map[string]struct{}, len(oldProps)+len(newProps))
	for name := range oldProps {
		allNames[name] = struct{}{}
	}
	for name := range newProps {
		allNames[name] = struct{}{}
	}
	names := sortedKeys(allNames)

	for _, name := range names {
		propPath := joinPointer(joinPointer(path, "properties"), name)
		oldProp, oldOK := oldProps[name]
		newProp, newOK := newProps[name]

		switch {
		case !oldOK && newOK:
			if _, req := newRequired[name]; req {
				*out = append(*out, Change{
					Class:   ChangeBreaking,
					Kind:    KindAddRequiredField,
					Path:    propPath,
					Message: fmt.Sprintf("added required field %q", name),
				})
			} else {
				*out = append(*out, Change{
					Class:   ChangeCompatible,
					Kind:    KindAddOptionalField,
					Path:    propPath,
					Message: fmt.Sprintf("added optional field %q", name),
				})
			}
			continue
		case oldOK && !newOK:
			*out = append(*out, Change{
				Class:   ChangeBreaking,
				Kind:    KindRemoveField,
				Path:    propPath,
				Message: fmt.Sprintf("removed field %q", name),
			})
			continue
		}

		oldObj, oldIsObj := oldProp.(map[string]any)
		newObj, newIsObj := newProp.(map[string]any)
		if oldIsObj && newIsObj {
			diffSchemaNode(propPath, oldObj, newObj, out)
		} else if !reflect.DeepEqual(oldProp, newProp) {
			*out = append(*out, Change{
				Class:   ChangeBreaking,
				Kind:    KindChangeFieldMeaning,
				Path:    propPath,
				Message: fmt.Sprintf("changed schema shape for field %q", name),
			})
		}

		_, wasRequired := oldRequired[name]
		_, isRequired := newRequired[name]
		switch {
		case !wasRequired && isRequired:
			*out = append(*out, Change{
				Class:   ChangeBreaking,
				Kind:    KindPromoteOptionalToRequired,
				Path:    propPath,
				Message: fmt.Sprintf("promoted optional field %q to required", name),
			})
		case wasRequired && !isRequired:
			*out = append(*out, Change{
				Class:   ChangeCompatible,
				Kind:    KindDemoteRequiredToOptional,
				Path:    propPath,
				Message: fmt.Sprintf("demoted required field %q to optional", name),
			})
		}
	}
}

func diffRequired(path string, oldNode, newNode map[string]any, out *[]Change) {
	// Property-level required transitions are handled in diffProperties.
	// Detect required arrays that mention names absent from both property maps
	// (malformed-but-present required entries) so they are not silently ignored.
	oldRequired := requiredSet(oldNode)
	newRequired := requiredSet(newNode)
	oldProps := propertiesMap(oldNode)
	newProps := propertiesMap(newNode)

	for _, name := range sortedKeys(newRequired) {
		if _, inOld := oldRequired[name]; inOld {
			continue
		}
		if _, inNewProps := newProps[name]; inNewProps {
			continue // already classified as add-required or promote
		}
		if _, inOldProps := oldProps[name]; inOldProps {
			continue
		}
		reqPath := joinPointer(joinPointer(path, "required"), name)
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindAddRequiredField,
			Path:    reqPath,
			Message: fmt.Sprintf("added required constraint for unknown field %q", name),
		})
	}
}

func diffEnum(path string, oldNode, newNode map[string]any, out *[]Change) {
	oldEnum, oldOK := asSlice(oldNode["enum"])
	newEnum, newOK := asSlice(newNode["enum"])
	enumPath := joinPointer(path, "enum")

	switch {
	case !oldOK && !newOK:
		return
	case !oldOK && newOK:
		// Introducing an enum constrains previously open values → Breaking.
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindNarrowEnum,
			Path:    enumPath,
			Message: "added enum constraint (narrows accepted values)",
		})
		return
	case oldOK && !newOK:
		// Removing an enum widens acceptance → Compatible.
		*out = append(*out, Change{
			Class:   ChangeCompatible,
			Kind:    KindWidenValidationRange,
			Path:    enumPath,
			Message: "removed enum constraint (widens accepted values)",
		})
		return
	}

	oldSet := enumValueSet(oldEnum)
	newSet := enumValueSet(newEnum)

	added := false
	removed := false
	for k := range newSet {
		if _, ok := oldSet[k]; !ok {
			added = true
			break
		}
	}
	for k := range oldSet {
		if _, ok := newSet[k]; !ok {
			removed = true
			break
		}
	}

	switch {
	case removed && !added:
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindNarrowEnum,
			Path:    enumPath,
			Message: "narrowed enum by removing one or more values",
		})
	case added && !removed:
		*out = append(*out, Change{
			Class:   ChangeReviewRequired,
			Kind:    KindAddEnumValue,
			Path:    enumPath,
			Message: "added one or more enum values (compatibility review required)",
		})
	case added && removed:
		// Replace/reorder with different membership: treat as narrow + add.
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindNarrowEnum,
			Path:    enumPath,
			Message: "changed enum membership (removed values; breaking)",
		})
		*out = append(*out, Change{
			Class:   ChangeReviewRequired,
			Kind:    KindAddEnumValue,
			Path:    enumPath,
			Message: "changed enum membership (added values; compatibility review required)",
		})
	}
}

func diffValidationRange(path string, oldNode, newNode map[string]any, out *[]Change) {
	// String length bounds.
	diffNumericBound(path, "minLength", oldNode, newNode, true /* higher is narrower */, out)
	diffNumericBound(path, "maxLength", oldNode, newNode, false /* lower is narrower */, out)
	// Numeric bounds.
	diffNumericBound(path, "minimum", oldNode, newNode, true, out)
	diffNumericBound(path, "maximum", oldNode, newNode, false, out)

	oldPattern, oldHas := oldNode["pattern"].(string)
	newPattern, newHas := newNode["pattern"].(string)
	patternPath := joinPointer(path, "pattern")
	switch {
	case !oldHas && newHas:
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindNarrowValidationRange,
			Path:    patternPath,
			Message: "added pattern constraint (narrows accepted values)",
		})
	case oldHas && !newHas:
		*out = append(*out, Change{
			Class:   ChangeCompatible,
			Kind:    KindWidenValidationRange,
			Path:    patternPath,
			Message: "removed pattern constraint (widens accepted values)",
		})
	case oldHas && newHas && oldPattern != newPattern:
		// Pattern changes can exclude previously valid values → Breaking.
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindNarrowValidationRange,
			Path:    patternPath,
			Message: "changed pattern constraint",
		})
	}
}

func diffNumericBound(path, key string, oldNode, newNode map[string]any, higherIsNarrower bool, out *[]Change) {
	oldVal, oldOK := asFloat(oldNode[key])
	newVal, newOK := asFloat(newNode[key])
	boundPath := joinPointer(path, key)

	switch {
	case !oldOK && !newOK:
		return
	case !oldOK && newOK:
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindNarrowValidationRange,
			Path:    boundPath,
			Message: fmt.Sprintf("added %s constraint (narrows accepted values)", key),
		})
		return
	case oldOK && !newOK:
		*out = append(*out, Change{
			Class:   ChangeCompatible,
			Kind:    KindWidenValidationRange,
			Path:    boundPath,
			Message: fmt.Sprintf("removed %s constraint (widens accepted values)", key),
		})
		return
	}

	if oldVal == newVal {
		return
	}

	narrowed := false
	if higherIsNarrower {
		narrowed = newVal > oldVal
	} else {
		narrowed = newVal < oldVal
	}

	if narrowed {
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindNarrowValidationRange,
			Path:    boundPath,
			Message: fmt.Sprintf("narrowed %s from %v to %v", key, oldVal, newVal),
		})
		return
	}
	*out = append(*out, Change{
		Class:   ChangeCompatible,
		Kind:    KindWidenValidationRange,
		Path:    boundPath,
		Message: fmt.Sprintf("widened %s from %v to %v", key, oldVal, newVal),
	})
}

func diffType(path string, oldNode, newNode map[string]any, out *[]Change) {
	oldType, oldOK := oldNode["type"]
	newType, newOK := newNode["type"]
	if !oldOK && !newOK {
		return
	}
	if reflect.DeepEqual(oldType, newType) {
		return
	}
	*out = append(*out, Change{
		Class:   ChangeBreaking,
		Kind:    KindChangeFieldMeaning,
		Path:    joinPointer(path, "type"),
		Message: fmt.Sprintf("changed type from %v to %v", oldType, newType),
	})
}

func diffRef(path string, oldNode, newNode map[string]any, out *[]Change) {
	oldRef, oldOK := oldNode["$ref"].(string)
	newRef, newOK := newNode["$ref"].(string)
	switch {
	case !oldOK && !newOK:
		return
	case oldOK && newOK && oldRef == newRef:
		return
	default:
		*out = append(*out, Change{
			Class:   ChangeReviewRequired,
			Kind:    KindChangeReferenceTarget,
			Path:    joinPointer(path, "$ref"),
			Message: "changed reference target ($ref); compatibility review required",
		})
	}
}

func diffAllowedScopes(path string, oldNode, newNode map[string]any, out *[]Change) {
	oldScopes, oldOK := asStringSlice(oldNode[ExtAllowedScopes])
	newScopes, newOK := asStringSlice(newNode[ExtAllowedScopes])
	switch {
	case !oldOK && !newOK:
		return
	case oldOK && newOK && stringSliceEqual(oldScopes, newScopes):
		return
	default:
		*out = append(*out, Change{
			Class:   ChangeReviewRequired,
			Kind:    KindChangeAllowedScopes,
			Path:    joinPointer(path, ExtAllowedScopes),
			Message: "changed x-sovrunn-allowed-scopes (reference/scope target review required)",
		})
	}
}

func diffFieldPolicy(path string, oldNode, newNode map[string]any, out *[]Change) {
	oldPol, oldOK := oldNode[ExtFieldPolicy].(map[string]any)
	newPol, newOK := newNode[ExtFieldPolicy].(map[string]any)
	polPath := joinPointer(path, ExtFieldPolicy)

	switch {
	case !oldOK && !newOK:
		return
	case !oldOK && newOK:
		// Adding explicit field policy is additive documentation of ownership.
		*out = append(*out, Change{
			Class:   ChangeCompatible,
			Kind:    KindAddRegisteredExtension,
			Path:    polPath,
			Message: "added x-sovrunn-field-policy",
		})
		return
	case oldOK && !newOK:
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindRemoveRegisteredExtension,
			Path:    polPath,
			Message: "removed x-sovrunn-field-policy",
		})
		return
	}

	oldClass, _ := oldPol["classification"].(string)
	newClass, _ := newPol["classification"].(string)
	if oldClass != newClass {
		if isInternalClassification(oldClass) && isPublicFacingClassification(newClass) {
			*out = append(*out, Change{
				Class:   ChangeReviewRequired,
				Kind:    KindExposeInternalPublicly,
				Path:    joinPointer(polPath, "classification"),
				Message: fmt.Sprintf("exposes internal classification %q as %q (security/boundary review required)", oldClass, newClass),
			})
		} else {
			*out = append(*out, Change{
				Class:   ChangeBreaking,
				Kind:    KindChangeFieldMeaning,
				Path:    joinPointer(polPath, "classification"),
				Message: fmt.Sprintf("changed field-policy classification from %q to %q", oldClass, newClass),
			})
		}
	}

	oldWriter, _ := oldPol["authorizedWriter"].(string)
	newWriter, _ := newPol["authorizedWriter"].(string)
	if oldWriter != newWriter {
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindChangeOwnerOrMutability,
			Path:    joinPointer(polPath, "authorizedWriter"),
			Message: fmt.Sprintf("changed authorizedWriter from %q to %q", oldWriter, newWriter),
		})
	}

	oldMut, _ := oldPol["mutability"].(string)
	newMut, _ := newPol["mutability"].(string)
	if oldMut != newMut {
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindChangeOwnerOrMutability,
			Path:    joinPointer(polPath, "mutability"),
			Message: fmt.Sprintf("changed mutability from %q to %q", oldMut, newMut),
		})
	}
}

func diffRegisteredExtensions(path string, oldNode, newNode map[string]any, out *[]Change) {
	for _, ext := range RegisteredExtensionKeywords {
		if ext == ExtFieldPolicy || ext == ExtAllowedScopes {
			// Handled by dedicated comparators that apply table-specific rules.
			continue
		}
		_, oldOK := oldNode[ext]
		_, newOK := newNode[ext]
		extPath := joinPointer(path, ext)
		switch {
		case !oldOK && newOK:
			*out = append(*out, Change{
				Class:   ChangeCompatible,
				Kind:    KindAddRegisteredExtension,
				Path:    extPath,
				Message: fmt.Sprintf("added registered extension %s", ext),
			})
		case oldOK && !newOK:
			*out = append(*out, Change{
				Class:   ChangeBreaking,
				Kind:    KindRemoveRegisteredExtension,
				Path:    extPath,
				Message: fmt.Sprintf("removed registered extension %s", ext),
			})
		case oldOK && newOK && !reflect.DeepEqual(oldNode[ext], newNode[ext]):
			// Profile/boundary/stability vocabulary changes affect meaning.
			*out = append(*out, Change{
				Class:   ChangeBreaking,
				Kind:    KindChangeFieldMeaning,
				Path:    extPath,
				Message: fmt.Sprintf("changed registered extension %s", ext),
			})
		}
	}
}

func diffItems(path string, oldNode, newNode map[string]any, out *[]Change) {
	oldItems, oldOK := oldNode["items"].(map[string]any)
	newItems, newOK := newNode["items"].(map[string]any)
	switch {
	case oldOK && newOK:
		diffSchemaNode(joinPointer(path, "items"), oldItems, newItems, out)
	case !oldOK && newOK:
		*out = append(*out, Change{
			Class:   ChangeBreaking,
			Kind:    KindNarrowValidationRange,
			Path:    joinPointer(path, "items"),
			Message: "added items schema constraint",
		})
	case oldOK && !newOK:
		*out = append(*out, Change{
			Class:   ChangeCompatible,
			Kind:    KindWidenValidationRange,
			Path:    joinPointer(path, "items"),
			Message: "removed items schema constraint",
		})
	case reflect.DeepEqual(oldNode["items"], newNode["items"]):
		return
	default:
		if oldNode["items"] != nil || newNode["items"] != nil {
			*out = append(*out, Change{
				Class:   ChangeBreaking,
				Kind:    KindChangeFieldMeaning,
				Path:    joinPointer(path, "items"),
				Message: "changed items schema",
			})
		}
	}
}

func diffAdditionalProperties(path string, oldNode, newNode map[string]any, out *[]Change) {
	oldAP, oldOK := oldNode["additionalProperties"]
	newAP, newOK := newNode["additionalProperties"]
	if !oldOK && !newOK {
		return
	}
	if reflect.DeepEqual(oldAP, newAP) {
		return
	}

	apPath := joinPointer(path, "additionalProperties")

	// true → false / schema: narrows; false → true: widens.
	oldBool, oldIsBool := oldAP.(bool)
	newBool, newIsBool := newAP.(bool)
	if oldIsBool && newIsBool {
		if oldBool && !newBool {
			*out = append(*out, Change{
				Class:   ChangeBreaking,
				Kind:    KindNarrowValidationRange,
				Path:    apPath,
				Message: "narrowed additionalProperties from true to false",
			})
			return
		}
		if !oldBool && newBool {
			*out = append(*out, Change{
				Class:   ChangeCompatible,
				Kind:    KindWidenValidationRange,
				Path:    apPath,
				Message: "widened additionalProperties from false to true",
			})
			return
		}
	}

	oldObj, oldIsObj := oldAP.(map[string]any)
	newObj, newIsObj := newAP.(map[string]any)
	if oldIsObj && newIsObj {
		diffSchemaNode(apPath, oldObj, newObj, out)
		return
	}

	*out = append(*out, Change{
		Class:   ChangeBreaking,
		Kind:    KindChangeFieldMeaning,
		Path:    apPath,
		Message: "changed additionalProperties constraint",
	})
}

func propertiesMap(node map[string]any) map[string]any {
	props, ok := node["properties"].(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return props
}

func requiredSet(node map[string]any) map[string]struct{} {
	raw, ok := node["required"].([]any)
	if !ok {
		return map[string]struct{}{}
	}
	out := make(map[string]struct{}, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		out[s] = struct{}{}
	}
	return out
}

func asSlice(v any) ([]any, bool) {
	s, ok := v.([]any)
	return s, ok
}

func asFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		if math.IsNaN(n) || math.IsInf(n, 0) {
			return 0, false
		}
		return n, true
	case json.Number:
		f, err := n.Float64()
		if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
			return 0, false
		}
		return f, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func enumValueSet(values []any) map[string]struct{} {
	out := make(map[string]struct{}, len(values))
	for _, v := range values {
		out[canonicalJSON(v)] = struct{}{}
	}
	return out
}

func canonicalJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%T:%v", v, v)
	}
	return string(b)
}

func asStringSlice(v any) ([]string, bool) {
	raw, ok := v.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	// Order-insensitive for allowed-scopes vocabulary sets.
	as := append([]string(nil), a...)
	bs := append([]string(nil), b...)
	sort.Strings(as)
	sort.Strings(bs)
	for i := range as {
		if as[i] != bs[i] {
			return false
		}
	}
	return true
}

func isInternalClassification(c string) bool {
	return c == "Internal" || c == "Operator-confidential" || c == "Sensitive" || c == "Secret-reference-only"
}

func isPublicFacingClassification(c string) bool {
	return c == "Public" || c == "Customer-visible"
}

func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortChanges(changes []Change) {
	sort.SliceStable(changes, func(i, j int) bool {
		if changes[i].Path != changes[j].Path {
			return changes[i].Path < changes[j].Path
		}
		if changes[i].Kind != changes[j].Kind {
			return changes[i].Kind < changes[j].Kind
		}
		if changes[i].Class != changes[j].Class {
			return changes[i].Class < changes[j].Class
		}
		return changes[i].Message < changes[j].Message
	})
}

// Baseline governance filenames under api/schemas/baseline/ (D-11).
const (
	BaselineManifestFileName  = "BASELINE_MANIFEST.json"
	BaselineApprovalsFileName = "BASELINE_APPROVALS.json"
)

// baselineManifest is the on-disk shape of BASELINE_MANIFEST.json.
// Digests are lowercase hex-encoded SHA-256 of the baseline file bytes
// (optionally prefixed with "sha256:" in the file; comparisons normalize).
type baselineManifest struct {
	Files map[string]string `json:"files"`
}

// baselineApprovalsFile is the on-disk shape of BASELINE_APPROVALS.json.
//
// RecordedDigests holds the last approved digest per baseline-relative path.
// An absent/empty map is the initial-bootstrap case: no prior-approval
// evidence is required (task 10.3). When RecordedDigests is populated,
// any current digest that differs from the recorded value is a baseline
// change and MUST have matching approval evidence — co-editing the baseline
// file and BASELINE_MANIFEST.json alone is never sufficient (D-11).
type baselineApprovalsFile struct {
	RecordedDigests map[string]string  `json:"recordedDigests"`
	Approvals       []baselineApproval `json:"approvals"`
}

// baselineApproval is one recorded baseline-change approval evidence entry.
// Exactly one of ADH or ApprovalToken must be non-empty, together with
// Path, OldDigest, NewDigest, Reviewer, and Date.
type baselineApproval struct {
	Path          string `json:"path"`
	OldDigest     string `json:"oldDigest"`
	NewDigest     string `json:"newDigest"`
	ADH           string `json:"adh"`
	ApprovalToken string `json:"approvalToken"`
	Reviewer      string `json:"reviewer"`
	Date          string `json:"date"`
}

// VerifyBaselineIntegrity recomputes SHA-256 digests of every baseline schema
// file under baselineDir and compares them to BASELINE_MANIFEST.json.
// A mismatch fails so a silent baseline edit is detected. This is an
// INTEGRITY check only — the manifest is not an independently unforgeable
// approval, since a committer can change a baseline and its digest together
// (D-11, F12-EVOLVE-002, F12-VERIFY-001(10)).
func VerifyBaselineIntegrity(manifestPath, baselineDir string) error {
	manifest, err := loadBaselineManifest(manifestPath)
	if err != nil {
		return err
	}
	actual, err := computeBaselineDigests(baselineDir)
	if err != nil {
		return err
	}
	return compareBaselineDigests(manifest.Files, actual)
}

// VerifyBaselineApproval enforces that any baseline change is APPROVED, not
// merely digest-consistent. When a baseline file's digest differs from the
// prior recorded digest in BASELINE_APPROVALS.json, it requires matching
// approval evidence containing the exact old/new digests, the approving ADH
// or approval token, the reviewer, and the date. Changing the baseline and
// its manifest in one commit without that evidence is NOT sufficient.
// The human governance boundary remains protected review / CODEOWNERS on the
// baseline and its approval record (D-11, F12-EVOLVE-002, F12-VERIFY-001(10)).
func VerifyBaselineApproval(approvalsPath, manifestPath, baselineDir string) error {
	if err := VerifyBaselineIntegrity(manifestPath, baselineDir); err != nil {
		return fmt.Errorf("baseline integrity prerequisite failed: %w", err)
	}

	approvals, err := loadBaselineApprovals(approvalsPath)
	if err != nil {
		return err
	}

	actual, err := computeBaselineDigests(baselineDir)
	if err != nil {
		return err
	}

	recorded := approvals.RecordedDigests
	if len(recorded) == 0 {
		// Initial baseline bootstrap: no prior-approval evidence required.
		return nil
	}

	// Every recorded path that disappeared, and every current path whose
	// digest differs from the recorded value, needs matching evidence.
	paths := make(map[string]struct{}, len(recorded)+len(actual))
	for p := range recorded {
		paths[p] = struct{}{}
	}
	for p := range actual {
		paths[p] = struct{}{}
	}

	for _, path := range sortedKeys(paths) {
		oldDigest, hadOld := recorded[path]
		newDigest, hasNew := actual[path]
		oldDigest = normalizeDigest(oldDigest)
		newDigest = normalizeDigest(newDigest)

		switch {
		case hadOld && hasNew && oldDigest == newDigest:
			continue // unchanged
		case !hadOld && hasNew:
			// New baseline file relative to recorded set.
			if err := requireApprovalEvidence(approvals.Approvals, path, "", newDigest); err != nil {
				return err
			}
		case hadOld && !hasNew:
			return fmt.Errorf("baseline approval: path %q was recorded but is missing from baseline directory (deletion requires approval evidence workflow)", path)
		default:
			// Digest changed: require evidence with exact old and new digests.
			if err := requireApprovalEvidence(approvals.Approvals, path, oldDigest, newDigest); err != nil {
				return err
			}
		}
	}
	return nil
}

func requireApprovalEvidence(entries []baselineApproval, path, oldDigest, newDigest string) error {
	oldDigest = normalizeDigest(oldDigest)
	newDigest = normalizeDigest(newDigest)
	for _, e := range entries {
		if e.Path != path {
			continue
		}
		if normalizeDigest(e.OldDigest) != oldDigest {
			continue
		}
		if normalizeDigest(e.NewDigest) != newDigest {
			continue
		}
		if err := validateApprovalEvidenceFields(e); err != nil {
			return fmt.Errorf("baseline approval: path %q: %w", path, err)
		}
		return nil
	}
	return fmt.Errorf("baseline approval: path %q changed from digest %q to %q without recorded approval evidence (oldDigest/newDigest/adh-or-token/reviewer/date); co-editing baseline and manifest is not sufficient", path, oldDigest, newDigest)
}

func validateApprovalEvidenceFields(e baselineApproval) error {
	adh := strings.TrimSpace(e.ADH)
	token := strings.TrimSpace(e.ApprovalToken)
	if adh == "" && token == "" {
		return fmt.Errorf("approval evidence missing approving ADH or approval token")
	}
	if strings.TrimSpace(e.Reviewer) == "" {
		return fmt.Errorf("approval evidence missing reviewer identity")
	}
	if strings.TrimSpace(e.Date) == "" {
		return fmt.Errorf("approval evidence missing date")
	}
	if strings.TrimSpace(e.Path) == "" {
		return fmt.Errorf("approval evidence missing path")
	}
	if normalizeDigest(e.NewDigest) == "" {
		return fmt.Errorf("approval evidence missing newDigest")
	}
	// oldDigest may be empty for a newly introduced baseline file.
	return nil
}

func loadBaselineManifest(path string) (baselineManifest, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return baselineManifest{}, fmt.Errorf("baseline manifest: read %s: %w", path, err)
	}
	var manifest baselineManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return baselineManifest{}, fmt.Errorf("baseline manifest: parse %s: %w", path, err)
	}
	if manifest.Files == nil {
		manifest.Files = map[string]string{}
	}
	normalized := make(map[string]string, len(manifest.Files))
	for p, d := range manifest.Files {
		p = filepath.ToSlash(p)
		if p == "" || isBaselineGovernanceFile(p) {
			return baselineManifest{}, fmt.Errorf("baseline manifest: invalid file entry %q", p)
		}
		nd := normalizeDigest(d)
		if nd == "" || !isSHA256Hex(nd) {
			return baselineManifest{}, fmt.Errorf("baseline manifest: invalid digest for %q", p)
		}
		normalized[p] = nd
	}
	manifest.Files = normalized
	return manifest, nil
}

func loadBaselineApprovals(path string) (baselineApprovalsFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return baselineApprovalsFile{}, fmt.Errorf("baseline approvals: read %s: %w", path, err)
	}
	// Allow a literally empty file or empty JSON object for the initial baseline.
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return baselineApprovalsFile{}, nil
	}
	var approvals baselineApprovalsFile
	if err := json.Unmarshal(raw, &approvals); err != nil {
		return baselineApprovalsFile{}, fmt.Errorf("baseline approvals: parse %s: %w", path, err)
	}
	if approvals.RecordedDigests == nil {
		approvals.RecordedDigests = map[string]string{}
	}
	normalized := make(map[string]string, len(approvals.RecordedDigests))
	for p, d := range approvals.RecordedDigests {
		p = filepath.ToSlash(p)
		nd := normalizeDigest(d)
		if p == "" || isBaselineGovernanceFile(p) {
			return baselineApprovalsFile{}, fmt.Errorf("baseline approvals: invalid recordedDigests path %q", p)
		}
		if nd == "" || !isSHA256Hex(nd) {
			return baselineApprovalsFile{}, fmt.Errorf("baseline approvals: invalid recorded digest for %q", p)
		}
		normalized[p] = nd
	}
	approvals.RecordedDigests = normalized
	if approvals.Approvals == nil {
		approvals.Approvals = []baselineApproval{}
	}
	return approvals, nil
}

func computeBaselineDigests(baselineDir string) (map[string]string, error) {
	info, err := os.Stat(baselineDir)
	if err != nil {
		return nil, fmt.Errorf("baseline directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("baseline directory: %s is not a directory", baselineDir)
	}

	out := make(map[string]string)
	err = filepath.WalkDir(baselineDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(baselineDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if isBaselineGovernanceFile(rel) {
			return nil
		}
		// Baseline snapshots are JSON schema documents.
		if !strings.HasSuffix(rel, ".json") {
			return nil
		}
		digest, err := sha256FileHex(path)
		if err != nil {
			return fmt.Errorf("baseline digest %s: %w", rel, err)
		}
		out[rel] = digest
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func compareBaselineDigests(manifest, actual map[string]string) error {
	paths := make(map[string]struct{}, len(manifest)+len(actual))
	for p := range manifest {
		paths[p] = struct{}{}
	}
	for p := range actual {
		paths[p] = struct{}{}
	}
	for _, path := range sortedKeys(paths) {
		want, inManifest := manifest[path]
		got, onDisk := actual[path]
		switch {
		case inManifest && !onDisk:
			return fmt.Errorf("baseline integrity: path %q listed in manifest but missing from baseline directory", path)
		case !inManifest && onDisk:
			return fmt.Errorf("baseline integrity: path %q present in baseline directory but missing from manifest", path)
		case want != got:
			return fmt.Errorf("baseline integrity: digest mismatch for %q: manifest=%s actual=%s (tampered or stale baseline)", path, want, got)
		}
	}
	return nil
}

func sha256FileHex(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

func normalizeDigest(d string) string {
	d = strings.TrimSpace(d)
	d = strings.TrimPrefix(d, "sha256:")
	d = strings.TrimPrefix(d, "SHA256:")
	return strings.ToLower(strings.TrimSpace(d))
}

func isSHA256Hex(d string) bool {
	if len(d) != sha256.Size*2 {
		return false
	}
	for _, c := range d {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

func isBaselineGovernanceFile(rel string) bool {
	base := filepath.Base(rel)
	return base == BaselineManifestFileName || base == BaselineApprovalsFileName
}
