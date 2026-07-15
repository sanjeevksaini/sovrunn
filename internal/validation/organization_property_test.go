package validation

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"testing/quick"
	"unicode"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: organization-resource-registry, Property 1: ValidateOrganization rejects invalid names
func TestProperty_ValidateOrganization_InvalidNames(t *testing.T) {
	f := func(name string) bool {
		if isValidDNSLabel(name) {
			return true
		}
		errs := ValidateOrganization(context.Background(), resources.Organization{
			Metadata: resources.Metadata{Name: name},
		})
		for _, e := range errs {
			if e.Field == "metadata.name" {
				return true
			}
		}
		return false
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: organization-resource-registry, Property 2: ValidateOrganization accepts all valid DNS-label names
func TestProperty_ValidateOrganization_ValidNames(t *testing.T) {
	f := func(name string) bool {
		if !isValidDNSLabel(name) {
			return true
		}
		errs := ValidateOrganization(context.Background(), resources.Organization{
			Metadata: resources.Metadata{Name: name},
		})
		for _, e := range errs {
			if e.Field == "metadata.name" {
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

var dnsLabelPattern = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

func isValidDNSLabel(name string) bool {
	if name == "" || len(name) > 63 {
		return false
	}
	if !dnsLabelPattern.MatchString(name) {
		return false
	}
	for _, r := range name {
		if r > unicode.MaxASCII || (!unicode.IsLower(r) && !unicode.IsDigit(r) && r != '-') {
			return false
		}
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return false
	}
	return true
}
