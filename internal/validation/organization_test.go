package validation

import (
	"context"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateOrganization_ValidNames(t *testing.T) {
	valid := []string{"a", "a1", "a-b", strings.Repeat("a", 63)}
	for _, name := range valid {
		errs := ValidateOrganization(context.Background(), resources.Organization{
			Metadata: resources.Metadata{Name: name},
		})
		if len(errs) != 0 {
			t.Errorf("name %q: got errors %v, want none", name, errs)
		}
	}
}

func TestValidateOrganization_EmptyName(t *testing.T) {
	errs := ValidateOrganization(context.Background(), resources.Organization{})
	if len(errs) != 1 || errs[0].Field != "metadata.name" {
		t.Fatalf("got %v, want single metadata.name error", errs)
	}
}

func TestValidateOrganization_Uppercase(t *testing.T) {
	errs := ValidateOrganization(context.Background(), resources.Organization{
		Metadata: resources.Metadata{Name: "ABC"},
	})
	if len(errs) != 1 || errs[0].Field != "metadata.name" {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganization_Spaces(t *testing.T) {
	errs := ValidateOrganization(context.Background(), resources.Organization{
		Metadata: resources.Metadata{Name: "a b"},
	})
	if len(errs) != 1 || errs[0].Field != "metadata.name" {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganization_TooLong(t *testing.T) {
	errs := ValidateOrganization(context.Background(), resources.Organization{
		Metadata: resources.Metadata{Name: strings.Repeat("a", 64)},
	})
	if len(errs) != 1 || errs[0].Field != "metadata.name" {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganization_LeadingHyphen(t *testing.T) {
	errs := ValidateOrganization(context.Background(), resources.Organization{
		Metadata: resources.Metadata{Name: "-abc"},
	})
	if len(errs) != 1 || errs[0].Field != "metadata.name" {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganization_TrailingHyphen(t *testing.T) {
	errs := ValidateOrganization(context.Background(), resources.Organization{
		Metadata: resources.Metadata{Name: "abc-"},
	})
	if len(errs) != 1 || errs[0].Field != "metadata.name" {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganization_SingleHyphen(t *testing.T) {
	errs := ValidateOrganization(context.Background(), resources.Organization{
		Metadata: resources.Metadata{Name: "-"},
	})
	if len(errs) != 1 || errs[0].Field != "metadata.name" {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganization_SpecialChars(t *testing.T) {
	errs := ValidateOrganization(context.Background(), resources.Organization{
		Metadata: resources.Metadata{Name: "a.b"},
	})
	if len(errs) != 1 || errs[0].Field != "metadata.name" {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateNamePath(t *testing.T) {
	errs := ValidateNamePath(context.Background(), "valid-name")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}
