package server

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestServer_Start_FailsWhenPortInUse_ReadinessFalse(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}
	defer ln.Close()

	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("expected TCP address")
	}

	cfg := config.Config{
		Server: config.ServerConfig{
			Host:            "127.0.0.1",
			Port:            tcpAddr.Port,
			ShutdownTimeout: 30,
		},
	}

	readiness := &health.ReadinessState{}
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	tenantRegistry := registry.NewTenantRegistry()
	projectRegistry := registry.NewProjectRegistry()
	operationRegistry := registry.NewOperationRegistry()
	orgHandler := api.NewOrgHandler(orgRegistry, registry.NoopChildBlockerChecker{}, nil)
	ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, nil, nil)
	tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry, nil, nil)
	projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry, nil)
	operationHandler := api.NewOperationHandler(operationRegistry)
	bootstrap := api.NewBootstrapHandler(cfg, readiness)
	srv := New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler, operationHandler, bootstrap, readiness)

	if err := srv.Start(); err == nil {
		t.Fatal("Start() expected error when port is already in use")
	}
	if readiness.IsReady() {
		t.Fatal("readiness should remain false when listener bind fails")
	}
}

// newTestServer builds a Server with all handlers wired for route testing.
func newTestServer() *Server {
	cfg := config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1", Port: 0, ShutdownTimeout: 30},
	}
	readiness := &health.ReadinessState{}
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	tenantRegistry := registry.NewTenantRegistry()
	projectRegistry := registry.NewProjectRegistry()
	operationRegistry := registry.NewOperationRegistry()
	orgHandler := api.NewOrgHandler(orgRegistry, registry.NoopChildBlockerChecker{}, nil)
	ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, nil, nil)
	tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry, nil, nil)
	projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry, nil)
	operationHandler := api.NewOperationHandler(operationRegistry)
	bootstrap := api.NewBootstrapHandler(cfg, readiness)
	return New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler, operationHandler, bootstrap, readiness)
}

func TestServer_TenantRoutes_Registered(t *testing.T) {
	srv := newTestServer()

	t.Run("collection GET returns empty list", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/tenants", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /v1/tenants status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("item bad path shape returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/tenants/only-one-segment", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/tenants/only-one-segment status = %d, want 404", rec.Code)
		}
	})

	t.Run("collection unsupported method returns 405", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/tenants", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("DELETE /v1/tenants status = %d, want 405", rec.Code)
		}
	})
}

func TestServer_ProjectRoutes_Registered(t *testing.T) {
	srv := newTestServer()

	t.Run("collection GET returns empty list", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/projects", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /v1/projects status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
		if rec.Body.String() != "{\"items\":[]}\n" {
			t.Fatalf("GET /v1/projects body = %q, want {\"items\":[]}", rec.Body.String())
		}
	})

	t.Run("item bad path shape returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/projects/only-one-segment", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/projects/only-one-segment status = %d, want 404", rec.Code)
		}
	})

	t.Run("collection unsupported method returns 405", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/projects", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("DELETE /v1/projects status = %d, want 405", rec.Code)
		}
	})
}

func TestServer_OperationRoutes_Registered(t *testing.T) {
	srv := newTestServer()

	t.Run("collection GET returns empty list", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/operations", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /v1/operations status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
		if rec.Body.String() != "{\"items\":[]}\n" {
			t.Fatalf("GET /v1/operations body = %q, want {\"items\":[]}", rec.Body.String())
		}
	})

	t.Run("bare item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/operations/", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/operations/ status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
	})

	t.Run("collection POST returns 405", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/operations", nil)
		req.Header.Set("Content-Type", "application/json")
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("POST /v1/operations status = %d, want 405", rec.Code)
		}
	})
}
