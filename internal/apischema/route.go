package apischema

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Stable route-form diagnostic codes (D-09, F12-NAMING-004).
const (
	CodeRouteEmpty          = "ROUTE_EMPTY"
	CodeRouteUnversioned    = "ROUTE_UNVERSIONED"
	CodeRouteMalformed      = "ROUTE_MALFORMED"
	CodeRouteInvalidGroup   = "ROUTE_INVALID_GROUP"
	CodeRouteInvalidVersion = "ROUTE_INVALID_VERSION"
	CodeRouteInvalidPlural  = "ROUTE_INVALID_PLURAL"
	CodeRouteInvalidSegment = "ROUTE_INVALID_SEGMENT"
)

const (
	routePrefix        = "/apis/"
	maxAPIGroupChars   = 253
	maxRouteLabelChars = 63
)

// dnsSubdomainRe is a lowercase DNS-1123 subdomain (labels separated by dots).
var dnsSubdomainRe = regexp.MustCompile(
	`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)+$`,
)

// kebabLabelRe is a lowercase kebab-case DNS label (plural collection or name).
var kebabLabelRe = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

// pathParamRe is a route-template placeholder such as {tenant}.
var pathParamRe = regexp.MustCompile(`^\{[a-z][a-z0-9]*\}$`)

// versionishRe detects a version-shaped segment so unapproved versions can be
// distinguished from a missing version slot (unversioned under /apis/).
var versionishRe = regexp.MustCompile(`^v[0-9]`)

// RouteError is the typed failure returned by ValidateRoute.
type RouteError struct {
	Code    string
	Message string
}

// Error implements the error interface.
func (e *RouteError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code == "" {
		return e.Message
	}
	if e.Message == "" {
		return e.Code
	}
	return e.Code + ": " + e.Message
}

// ValidateRoute enforces the adopting API route form
// `/apis/<group>/<version>/<plural-kebab>` (D-09, F12-NAMING-004).
//
// Scoped collections may nest parent scope segments after the first plural,
// for example:
//
//	/apis/core.sovrunn.io/v1alpha1/tenants/{tenant}/projects
//
// Unversioned public endpoints (including Phase 1 compatibility paths such as
// `/organizations`) are rejected. FEATURE-0012 registers no runtime routes;
// this helper is grammar-only for adopters and conformance checks.
func ValidateRoute(path string) error {
	if path == "" {
		return &RouteError{
			Code:    CodeRouteEmpty,
			Message: "route path is required",
		}
	}
	if strings.ContainsAny(path, "?#") {
		return &RouteError{
			Code:    CodeRouteMalformed,
			Message: "route path must not include query or fragment",
		}
	}
	if strings.Contains(path, "//") {
		return &RouteError{
			Code:    CodeRouteMalformed,
			Message: "route path must not contain empty segments",
		}
	}
	if strings.HasSuffix(path, "/") {
		return &RouteError{
			Code:    CodeRouteMalformed,
			Message: "route path must not have a trailing slash",
		}
	}

	// Unversioned public endpoints (Phase 1 style and any non-/apis/ surface).
	if !strings.HasPrefix(path, routePrefix) {
		return &RouteError{
			Code:    CodeRouteUnversioned,
			Message: "route must use /apis/<group>/<version>/<plural-kebab>; unversioned public endpoints are rejected",
		}
	}

	// path = /apis/<group>/<version>/<plural>[/...]
	rest := strings.TrimPrefix(path, routePrefix)
	segments := strings.Split(rest, "/")
	for _, seg := range segments {
		if seg == "" {
			return &RouteError{
				Code:    CodeRouteMalformed,
				Message: "route path must not contain empty segments",
			}
		}
	}
	// /apis/<group>/<plural> (and shorter) omit the version slot.
	if len(segments) < 3 {
		return &RouteError{
			Code:    CodeRouteUnversioned,
			Message: "route must include an approved version segment (v1alpha1, v1beta1, or v1)",
		}
	}

	group := segments[0]
	version := segments[1]
	pluralAndRest := segments[2:]

	if err := validateAPIGroup(group); err != nil {
		return err
	}

	if !apimeta.IsKnownVersion(version) {
		if versionishRe.MatchString(version) {
			return &RouteError{
				Code:    CodeRouteInvalidVersion,
				Message: fmt.Sprintf("version %q must be one of v1alpha1, v1beta1, or v1", version),
			}
		}
		return &RouteError{
			Code:    CodeRouteUnversioned,
			Message: "route must include an approved version segment (v1alpha1, v1beta1, or v1)",
		}
	}

	// First resource segment after version must be a plural-kebab collection.
	if err := validatePluralSegment(pluralAndRest[0]); err != nil {
		return err
	}

	// Remaining segments alternate name-or-param, plural-kebab, ...
	for i := 1; i < len(pluralAndRest); i++ {
		seg := pluralAndRest[i]
		if i%2 == 1 {
			if err := validateNameOrParamSegment(seg); err != nil {
				return err
			}
			continue
		}
		if err := validatePluralSegment(seg); err != nil {
			return err
		}
	}
	return nil
}

func validateAPIGroup(group string) error {
	if utf8.RuneCountInString(group) > maxAPIGroupChars {
		return &RouteError{
			Code:    CodeRouteInvalidGroup,
			Message: fmt.Sprintf("API group exceeds %d characters", maxAPIGroupChars),
		}
	}
	if !isSovrunnRouteGroup(group) {
		return &RouteError{
			Code:    CodeRouteInvalidGroup,
			Message: "API group must be a lowercase DNS-style domain under sovrunn.io",
		}
	}
	if group != "sovrunn.io" && !dnsSubdomainRe.MatchString(group) {
		return &RouteError{
			Code:    CodeRouteInvalidGroup,
			Message: "API group must be a lowercase DNS-style domain under sovrunn.io",
		}
	}
	return nil
}

func isSovrunnRouteGroup(group string) bool {
	return group == "sovrunn.io" || strings.HasSuffix(group, ".sovrunn.io")
}

func validatePluralSegment(seg string) error {
	if utf8.RuneCountInString(seg) > maxRouteLabelChars || !kebabLabelRe.MatchString(seg) {
		return &RouteError{
			Code:    CodeRouteInvalidPlural,
			Message: "collection segment must be lowercase plural kebab-case DNS label, at most 63 characters",
		}
	}
	return nil
}

func validateNameOrParamSegment(seg string) error {
	if pathParamRe.MatchString(seg) {
		return nil
	}
	if utf8.RuneCountInString(seg) > maxRouteLabelChars || !kebabLabelRe.MatchString(seg) {
		return &RouteError{
			Code:    CodeRouteInvalidSegment,
			Message: "path segment must be a lowercase kebab-case name or {param} placeholder",
		}
	}
	return nil
}
