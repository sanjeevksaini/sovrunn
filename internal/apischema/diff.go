package apischema

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
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
// then Class for deterministic gate output. This task does not implement
// baseline integrity or approval checks (see VerifyBaselineIntegrity /
// VerifyBaselineApproval in a later task).
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
