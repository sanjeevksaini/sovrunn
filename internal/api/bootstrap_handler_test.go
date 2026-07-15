package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
)

func TestBootstrapHandler_Healthz(t *testing.T) {
	h := NewBootstrapHandler(config.Config{}, &health.ReadinessState{})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.Healthz(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want ok", body["status"])
	}
}

func TestBootstrapHandler_Readyz_Ready(t *testing.T) {
	rs := &health.ReadinessState{}
	rs.SetReady(true)
	h := NewBootstrapHandler(config.Config{}, rs)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	h.Readyz(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["status"] != "ready" {
		t.Errorf("status = %q, want ready", body["status"])
	}
}

func TestBootstrapHandler_Readyz_NotReady(t *testing.T) {
	h := NewBootstrapHandler(config.Config{}, &health.ReadinessState{})
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	h.Readyz(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["status"] != "not-ready" {
		t.Errorf("status = %q, want not-ready", body["status"])
	}
}

func TestBootstrapHandler_Version(t *testing.T) {
	h := NewBootstrapHandler(config.Config{}, &health.ReadinessState{})
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	h.Version(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	for _, key := range []string{"name", "version", "phase", "status"} {
		if body[key] == "" {
			t.Errorf("missing field %q", key)
		}
	}
	if body["name"] != "sovrunn-api" {
		t.Errorf("name = %q, want sovrunn-api", body["name"])
	}
}
