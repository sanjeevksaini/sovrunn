package apimeta

import (
	"encoding/json"
	"testing"
)

type listItem struct {
	Name string `json:"name"`
}

func TestListEnvelopeJSONPromotesTypeMeta(t *testing.T) {
	t.Parallel()

	env := ListEnvelope[listItem]{
		TypeMeta: TypeMeta{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "ProjectList",
		},
		Items: []listItem{{Name: "alpha"}, {Name: "beta"}},
		Page: Page{
			NextPageToken: "opaque-token-1",
		},
	}

	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("json.Unmarshal top-level: %v", err)
	}

	for _, key := range []string{"apiVersion", "kind", "items", "page"} {
		if _, ok := top[key]; !ok {
			t.Fatalf("expected top-level key %q in %s", key, raw)
		}
	}
	if _, nested := top["typeMeta"]; nested {
		t.Fatalf("TypeMeta must not nest under typeMeta; got %s", raw)
	}
	if _, nested := top["TypeMeta"]; nested {
		t.Fatalf("TypeMeta must not nest under TypeMeta; got %s", raw)
	}

	var decoded ListEnvelope[listItem]
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal envelope: %v", err)
	}
	if decoded.APIVersion != env.APIVersion {
		t.Fatalf("APIVersion = %q, want %q", decoded.APIVersion, env.APIVersion)
	}
	if decoded.Kind != env.Kind {
		t.Fatalf("Kind = %q, want %q", decoded.Kind, env.Kind)
	}
	if len(decoded.Items) != 2 {
		t.Fatalf("Items len = %d, want 2", len(decoded.Items))
	}
	if decoded.Page.NextPageToken != "opaque-token-1" {
		t.Fatalf("NextPageToken = %q, want opaque-token-1", decoded.Page.NextPageToken)
	}
}

func TestPageOmitsEmptyNextPageToken(t *testing.T) {
	t.Parallel()

	env := ListEnvelope[listItem]{
		TypeMeta: TypeMeta{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "ProjectList",
		},
		Items: []listItem{},
		Page:  Page{},
	}

	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	var page map[string]json.RawMessage
	if err := json.Unmarshal(top["page"], &page); err != nil {
		t.Fatalf("json.Unmarshal page: %v", err)
	}
	if _, ok := page["nextPageToken"]; ok {
		t.Fatalf("empty nextPageToken must be omitted; got %s", raw)
	}
}
