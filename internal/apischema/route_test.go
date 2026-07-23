package apischema

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateRouteValid(t *testing.T) {
	t.Parallel()

	valid := []string{
		"/apis/fabric.sovrunn.io/v1alpha1/resource-pools",
		"/apis/core.sovrunn.io/v1alpha1/tenants/{tenant}/projects",
		"/apis/core.sovrunn.io/v1alpha1/tenants/acme/projects",
		"/apis/core.sovrunn.io/v1alpha1/tenants/acme/projects/demo",
		"/apis/platform.sovrunn.io/v1beta1/plugins",
		"/apis/platform.sovrunn.io/v1/operations",
		"/apis/sovrunn.io/v1alpha1/projects",
	}
	for _, path := range valid {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			if err := ValidateRoute(path); err != nil {
				t.Fatalf("ValidateRoute(%q) = %v, want nil", path, err)
			}
		})
	}
}

func TestValidateRouteUnversionedRejected(t *testing.T) {
	t.Parallel()

	unversioned := []string{
		"/organizations",
		"/tenants",
		"/projects",
		"/healthz",
		"/readyz",
		"/v1alpha1/projects",
		"/apis/core.sovrunn.io/projects",
		"/apis/core.sovrunn.io/tenants/acme/projects",
	}
	for _, path := range unversioned {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			err := ValidateRoute(path)
			if err == nil {
				t.Fatalf("ValidateRoute(%q) = nil, want unversioned error", path)
			}
			var re *RouteError
			if !errors.As(err, &re) {
				t.Fatalf("ValidateRoute(%q) type=%T, want *RouteError", path, err)
			}
			if re.Code != CodeRouteUnversioned {
				t.Fatalf("ValidateRoute(%q) code=%q, want %q", path, re.Code, CodeRouteUnversioned)
			}
		})
	}
}

func TestValidateRouteMalformedGroupRejected(t *testing.T) {
	t.Parallel()

	cases := []struct {
		path string
		want string
	}{
		{path: "/apis/CORE.SOVRUNN.IO/v1alpha1/projects", want: CodeRouteInvalidGroup},
		{path: "/apis/core_sovrunn_io/v1alpha1/projects", want: CodeRouteInvalidGroup},
		{path: "/apis/core..sovrunn.io/v1alpha1/projects", want: CodeRouteInvalidGroup},
		{path: "/apis/-core.sovrunn.io/v1alpha1/projects", want: CodeRouteInvalidGroup},
		{path: "/apis/core.sovrunn.io-/v1alpha1/projects", want: CodeRouteInvalidGroup},
		{path: "/apis/amazonaws.com/v1alpha1/projects", want: CodeRouteInvalidGroup},
		{path: "/apis/fabric.example.com/v1alpha1/resource-pools", want: CodeRouteInvalidGroup},
		{path: "/apis/not a group/v1alpha1/projects", want: CodeRouteInvalidGroup},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			err := ValidateRoute(tc.path)
			if err == nil {
				t.Fatalf("ValidateRoute(%q) = nil, want group error", tc.path)
			}
			var re *RouteError
			if !errors.As(err, &re) {
				t.Fatalf("ValidateRoute(%q) type=%T, want *RouteError", tc.path, err)
			}
			if re.Code != tc.want {
				t.Fatalf("ValidateRoute(%q) code=%q, want %q (msg=%q)", tc.path, re.Code, tc.want, re.Message)
			}
		})
	}
}

func TestValidateRouteOtherFailures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		path string
		want string
	}{
		{name: "empty", path: "", want: CodeRouteEmpty},
		{name: "trailing-slash", path: "/apis/core.sovrunn.io/v1alpha1/projects/", want: CodeRouteMalformed},
		{name: "query", path: "/apis/core.sovrunn.io/v1alpha1/projects?x=1", want: CodeRouteMalformed},
		{name: "incomplete", path: "/apis/core.sovrunn.io/v1alpha1", want: CodeRouteUnversioned},
		{name: "bad-version", path: "/apis/core.sovrunn.io/v2/projects", want: CodeRouteInvalidVersion},
		{name: "bad-plural", path: "/apis/core.sovrunn.io/v1alpha1/Projects", want: CodeRouteInvalidPlural},
		{name: "bad-param", path: "/apis/core.sovrunn.io/v1alpha1/tenants/{Tenant}/projects", want: CodeRouteInvalidSegment},
		{name: "bad-name", path: "/apis/core.sovrunn.io/v1alpha1/tenants/Acme/projects", want: CodeRouteInvalidSegment},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateRoute(tc.path)
			if err == nil {
				t.Fatalf("ValidateRoute(%q) = nil, want code %q", tc.path, tc.want)
			}
			var re *RouteError
			if !errors.As(err, &re) {
				t.Fatalf("ValidateRoute(%q) type=%T, want *RouteError", tc.path, err)
			}
			if re.Code != tc.want {
				t.Fatalf("ValidateRoute(%q) code=%q, want %q (msg=%q)", tc.path, re.Code, tc.want, re.Message)
			}
			if strings.TrimSpace(re.Message) == "" {
				t.Fatalf("ValidateRoute(%q) message is empty", tc.path)
			}
		})
	}
}
