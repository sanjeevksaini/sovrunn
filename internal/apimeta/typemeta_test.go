package apimeta

import "testing"

func TestParseAPIVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		apiVersion  string
		wantGroup   string
		wantVersion string
		wantOK      bool
	}{
		{
			name:        "valid domain version",
			apiVersion:  "fabric.sovrunn.io/v1alpha1",
			wantGroup:   "fabric.sovrunn.io",
			wantVersion: "v1alpha1",
			wantOK:      true,
		},
		{
			name:        "stable v1",
			apiVersion:  "core.sovrunn.io/v1",
			wantGroup:   "core.sovrunn.io",
			wantVersion: "v1",
			wantOK:      true,
		},
		{
			name:       "empty",
			apiVersion: "",
			wantOK:     false,
		},
		{
			name:       "missing slash",
			apiVersion: "fabric.sovrunn.io",
			wantOK:     false,
		},
		{
			name:       "empty group",
			apiVersion: "/v1alpha1",
			wantOK:     false,
		},
		{
			name:       "empty version",
			apiVersion: "fabric.sovrunn.io/",
			wantOK:     false,
		},
		{
			name:       "extra segment",
			apiVersion: "fabric.sovrunn.io/v1alpha1/extra",
			wantOK:     false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			group, version, ok := ParseAPIVersion(tc.apiVersion)
			if ok != tc.wantOK {
				t.Fatalf("ok=%v, want %v", ok, tc.wantOK)
			}
			if !tc.wantOK {
				return
			}
			if group != tc.wantGroup || version != tc.wantVersion {
				t.Fatalf("got (%q, %q), want (%q, %q)", group, version, tc.wantGroup, tc.wantVersion)
			}
		})
	}
}

func TestTypeMetaGroupVersion(t *testing.T) {
	t.Parallel()

	tm := TypeMeta{APIVersion: "fabric.sovrunn.io/v1beta1", Kind: "ResourcePool"}
	if got := tm.Group(); got != "fabric.sovrunn.io" {
		t.Fatalf("Group()=%q, want fabric.sovrunn.io", got)
	}
	if got := tm.Version(); got != "v1beta1" {
		t.Fatalf("Version()=%q, want v1beta1", got)
	}
	if !IsKnownVersion(tm.Version()) {
		t.Fatal("v1beta1 should be a known version")
	}
	if IsKnownVersion("v2") {
		t.Fatal("v2 must not be a known version")
	}

	bad := TypeMeta{APIVersion: "not-a-version"}
	if bad.Group() != "" || bad.Version() != "" {
		t.Fatalf("unparsable TypeMeta should yield empty group/version, got %q/%q", bad.Group(), bad.Version())
	}
}
