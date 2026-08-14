package server

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
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
	defer func() {
		_ = ln.Close()
	}()

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
	serviceClassRegistry := registry.NewServiceClassRegistry()
	servicePlanRegistry := registry.NewServicePlanRegistry()
	serviceClassBlocker := registry.NewServicePlanChildBlockerChecker(servicePlanRegistry)
	orgHandler := api.NewOrgHandler(orgRegistry, registry.NoopChildBlockerChecker{}, nil)
	ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, nil, nil)
	tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry, nil, nil)
	projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry, nil)
	operationHandler := api.NewOperationHandler(operationRegistry)
	serviceClassHandler := api.NewServiceClassHandler(serviceClassRegistry, serviceClassBlocker, nil)
	servicePlanHandler := api.NewServicePlanHandler(servicePlanRegistry, serviceClassRegistry, nil)
	pluginRegistry := registry.NewPluginRegistry()
	capabilityRegistry := registry.NewCapabilityRegistry()
	pluginBlocker := registry.NewCapabilityChildBlockerChecker(capabilityRegistry)
	pluginHandler := api.NewPluginHandler(pluginRegistry, serviceClassRegistry, pluginBlocker, nil)
	capabilityHandler := api.NewCapabilityHandler(capabilityRegistry, pluginRegistry, serviceClassRegistry, nil)
	bootstrap := api.NewBootstrapHandler(cfg, readiness)
	srv := New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler, operationHandler, serviceClassHandler, servicePlanHandler, pluginHandler, capabilityHandler, nil, nil, bootstrap, readiness)

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
	return newTestServerWithConfig(cfg)
}

func newTestServerWithConfig(cfg config.Config) *Server {
	readiness := &health.ReadinessState{}
	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	tenantRegistry := registry.NewTenantRegistry()
	projectRegistry := registry.NewProjectRegistry()
	operationRegistry := registry.NewOperationRegistry()
	serviceClassRegistry := registry.NewServiceClassRegistry()
	servicePlanRegistry := registry.NewServicePlanRegistry()
	serviceClassBlocker := registry.NewServicePlanChildBlockerChecker(servicePlanRegistry)
	orgHandler := api.NewOrgHandler(orgRegistry, registry.NoopChildBlockerChecker{}, nil)
	ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, nil, nil)
	tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry, nil, nil)
	projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry, nil)
	operationHandler := api.NewOperationHandler(operationRegistry)
	serviceClassHandler := api.NewServiceClassHandler(serviceClassRegistry, serviceClassBlocker, nil)
	servicePlanHandler := api.NewServicePlanHandler(servicePlanRegistry, serviceClassRegistry, nil)
	pluginRegistry := registry.NewPluginRegistry()
	capabilityRegistry := registry.NewCapabilityRegistry()
	pluginBlocker := registry.NewCapabilityChildBlockerChecker(capabilityRegistry)
	pluginHandler := api.NewPluginHandler(pluginRegistry, serviceClassRegistry, pluginBlocker, nil)
	capabilityHandler := api.NewCapabilityHandler(capabilityRegistry, pluginRegistry, serviceClassRegistry, nil)
	bootstrap := api.NewBootstrapHandler(cfg, readiness)
	return New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler, operationHandler, serviceClassHandler, servicePlanHandler, pluginHandler, capabilityHandler, nil, nil, bootstrap, readiness)
}

func TestServer_BootstrapRoutes_RejectNonGET(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "POST healthz", method: http.MethodPost, path: "/healthz"},
		{name: "PUT readyz", method: http.MethodPut, path: "/readyz"},
		{name: "DELETE version", method: http.MethodDelete, path: "/version"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := newTestServer()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			srv.httpServer.Handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status = %d, want 405", tt.method, tt.path, rec.Code)
			}
			if allow := rec.Header().Get("Allow"); allow != http.MethodGet {
				t.Errorf("Allow = %q, want GET", allow)
			}
			if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", contentType)
			}

			var envelope resources.APIErrorEnvelope
			if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if envelope.Error.Code != resources.ErrCodeMethodNotAllowed {
				t.Errorf("error.code = %q, want %q", envelope.Error.Code, resources.ErrCodeMethodNotAllowed)
			}
			if envelope.Error.Message != "only GET is supported" {
				t.Errorf("error.message = %q, want %q", envelope.Error.Message, "only GET is supported")
			}
		})
	}
}

func TestServer_BootstrapRoutes_RejectHEADWithoutBody(t *testing.T) {
	srv := newTestServer()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodHead, "/healthz", nil)

	srv.httpServer.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("HEAD /healthz status = %d, want 405", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow != http.MethodGet {
		t.Errorf("Allow = %q, want GET", allow)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("HEAD /healthz body length = %d, want 0", rec.Body.Len())
	}
}

func TestServer_LiveReadinessAndShutdown(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		_ = listener.Close()
		t.Fatal("expected TCP address")
	}
	port := tcpAddr.Port
	if err := listener.Close(); err != nil {
		t.Fatalf("listener.Close() error = %v", err)
	}

	cfg := config.Config{
		Server: config.ServerConfig{
			Host:            "127.0.0.1",
			Port:            port,
			ShutdownTimeout: 1,
		},
	}
	srv := newTestServerWithConfig(cfg)
	startErr := make(chan error, 1)
	go func() {
		startErr <- srv.Start()
	}()
	t.Cleanup(func() {
		_ = srv.Shutdown(time.Second)
	})

	client := &http.Client{Timeout: time.Second}
	baseURL := "http://" + cfg.Addr()
	deadline := time.Now().Add(5 * time.Second)
	for {
		select {
		case err := <-startErr:
			t.Fatalf("Start() returned before readiness: %v", err)
		default:
		}

		resp, err := client.Get(baseURL + "/readyz")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("GET /readyz did not return 200 within 5s; last error = %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	resp, err := client.Get(baseURL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz error = %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /healthz status = %d, want 200", resp.StatusCode)
	}

	if err := srv.Shutdown(time.Second); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	srv.httpServer.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("GET /readyz after Shutdown status = %d, want 503", rec.Code)
	}
	var readinessResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&readinessResponse); err != nil {
		t.Fatalf("decode readiness response: %v", err)
	}
	if readinessResponse.Status != "not_ready" {
		t.Errorf("status = %q, want not_ready", readinessResponse.Status)
	}
	if readinessResponse.Message != health.ReasonShuttingDown {
		t.Errorf("message = %q, want %q", readinessResponse.Message, health.ReasonShuttingDown)
	}
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

func TestServer_ServiceClassRoutes_Registered(t *testing.T) {
	srv := newTestServer()

	t.Run("collection GET returns empty list", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/service-classes", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /v1/service-classes status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
		if rec.Body.String() != "{\"items\":[]}\n" {
			t.Fatalf("GET /v1/service-classes body = %q, want {\"items\":[]}", rec.Body.String())
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		if reqID := rec.Header().Get("X-Sovrunn-Request-ID"); reqID == "" {
			t.Error("X-Sovrunn-Request-ID missing; middleware chain should set it")
		}
	})

	t.Run("bare item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/service-classes/", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/service-classes/ status = %d, want 404", rec.Code)
		}
	})

	t.Run("item path reachable returns 404 for missing resource", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/service-classes/postgres", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/service-classes/postgres status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
		if reqID := rec.Header().Get("X-Sovrunn-Request-ID"); reqID == "" {
			t.Error("X-Sovrunn-Request-ID missing; middleware chain should set it")
		}
	})

	t.Run("wrong-shape item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/service-classes/postgres/extra", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/service-classes/postgres/extra status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
	})

	t.Run("collection unsupported method returns 405", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/service-classes", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("DELETE /v1/service-classes status = %d, want 405", rec.Code)
		}
	})
}

func TestServer_ServicePlanRoutes_Registered(t *testing.T) {
	srv := newTestServer()

	t.Run("collection GET returns empty list", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/service-plans", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /v1/service-plans status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
		if rec.Body.String() != "{\"items\":[]}\n" {
			t.Fatalf("GET /v1/service-plans body = %q, want {\"items\":[]}", rec.Body.String())
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		if reqID := rec.Header().Get("X-Sovrunn-Request-ID"); reqID == "" {
			t.Error("X-Sovrunn-Request-ID missing; middleware chain should set it")
		}
	})

	t.Run("bare item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/service-plans/", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/service-plans/ status = %d, want 404", rec.Code)
		}
	})

	t.Run("item path reachable returns 404 for missing resource", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/service-plans/postgres/small", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/service-plans/postgres/small status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
		if reqID := rec.Header().Get("X-Sovrunn-Request-ID"); reqID == "" {
			t.Error("X-Sovrunn-Request-ID missing; middleware chain should set it")
		}
	})

	t.Run("wrong-shape item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/service-plans/postgres", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/service-plans/postgres status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
	})

	t.Run("collection unsupported method returns 405", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/service-plans", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("DELETE /v1/service-plans status = %d, want 405", rec.Code)
		}
	})
}

func TestServer_PluginRoutes_Registered(t *testing.T) {
	srv := newTestServer()

	t.Run("collection GET returns empty list", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/plugins", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /v1/plugins status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
		if rec.Body.String() != "{\"items\":[]}\n" {
			t.Fatalf("GET /v1/plugins body = %q, want {\"items\":[]}", rec.Body.String())
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		if reqID := rec.Header().Get("X-Sovrunn-Request-ID"); reqID == "" {
			t.Error("X-Sovrunn-Request-ID missing; middleware chain should set it")
		}
	})

	t.Run("bare item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/plugins/", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/plugins/ status = %d, want 404", rec.Code)
		}
	})

	t.Run("item path reachable returns 404 for missing resource", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/plugins/postgres-ops", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/plugins/postgres-ops status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
		if reqID := rec.Header().Get("X-Sovrunn-Request-ID"); reqID == "" {
			t.Error("X-Sovrunn-Request-ID missing; middleware chain should set it")
		}
	})

	t.Run("wrong-shape item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/plugins/postgres-ops/extra", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/plugins/postgres-ops/extra status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
	})

	t.Run("collection unsupported method returns 405", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/plugins", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("DELETE /v1/plugins status = %d, want 405", rec.Code)
		}
	})
}

func TestServer_CapabilityRoutes_Registered(t *testing.T) {
	srv := newTestServer()

	t.Run("collection GET returns empty list", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/capabilities", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /v1/capabilities status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
		if rec.Body.String() != "{\"items\":[]}\n" {
			t.Fatalf("GET /v1/capabilities body = %q, want {\"items\":[]}", rec.Body.String())
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		if reqID := rec.Header().Get("X-Sovrunn-Request-ID"); reqID == "" {
			t.Error("X-Sovrunn-Request-ID missing; middleware chain should set it")
		}
	})

	t.Run("bare item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/capabilities/", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/capabilities/ status = %d, want 404", rec.Code)
		}
	})

	t.Run("item path reachable returns 404 for missing resource", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/capabilities/postgres-provision", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/capabilities/postgres-provision status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
		if reqID := rec.Header().Get("X-Sovrunn-Request-ID"); reqID == "" {
			t.Error("X-Sovrunn-Request-ID missing; middleware chain should set it")
		}
	})

	t.Run("wrong-shape item path returns 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/capabilities/postgres-provision/extra", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET /v1/capabilities/postgres-provision/extra status = %d, want 404", rec.Code)
		}
		var envelope resources.APIErrorEnvelope
		if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if envelope.Error.Code != resources.ErrCodeResourceNotFound {
			t.Fatalf("error.code = %q, want RESOURCE_NOT_FOUND", envelope.Error.Code)
		}
	})

	t.Run("collection unsupported method returns 405", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/capabilities", nil)
		srv.httpServer.Handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("DELETE /v1/capabilities status = %d, want 405", rec.Code)
		}
	})
}

func TestServer_CloudModelSchedulerStartStop(t *testing.T) {
	srv := newTestServer()
	audit := &cloudmodel.MemoryAuditAppender{}
	rt := NewCloudModelRuntime(audit)
	srv.AttachCloudModel(rt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rt.Scheduler.Start(ctx)
	if err := srv.Shutdown(2 * time.Second); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	// A second Stop after Shutdown must be safe and prove the scheduler halted.
	rt.Scheduler.Stop()
}

func TestServer_AuditBeforeExpiryPublication(t *testing.T) {
	audit := &cloudmodel.MemoryAuditAppender{}
	rt := NewCloudModelRuntime(audit)
	seedPendingParticipation(t, rt.Store)

	if err := rt.Publication.ExpireParticipation(context.Background(), "cccccccccccccccccccccccccccccccc", "1", "sched-1"); err != nil {
		t.Fatalf("ExpireParticipation: %v", err)
	}
	got, ok := rt.Store.GetParticipation("cccccccccccccccccccccccccccccccc")
	if !ok || got.Status.Phase != model.ParticipationPhaseExpired {
		t.Fatalf("want Expired, got %#v ok=%v", got.Status, ok)
	}
	if audit.Len() != 1 {
		t.Fatalf("audit len = %d, want 1", audit.Len())
	}
	if audit.Events()[0].Record.Action != cloudmodel.AuditActionParticipationExpiry {
		t.Fatalf("action = %s", audit.Events()[0].Record.Action)
	}
}

func TestServer_ExpiryAppendFailurePreventsPublication(t *testing.T) {
	audit := &cloudmodel.MemoryAuditAppender{}
	audit.SetFail(errors.New("append boom"))
	rt := NewCloudModelRuntime(audit)
	seedPendingParticipation(t, rt.Store)

	err := rt.Publication.ExpireParticipation(context.Background(), "cccccccccccccccccccccccccccccccc", "1", "sched-fail")
	if err == nil || err.Error() != string(apiproblem.CodeInternalError) {
		t.Fatalf("want INTERNAL_ERROR, got %v", err)
	}
	got, _ := rt.Store.GetParticipation("cccccccccccccccccccccccccccccccc")
	if got.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("phase = %s, want Pending", got.Status.Phase)
	}
	if audit.Len() != 0 {
		t.Fatal("append failure must not publish AuditEvent or expiry")
	}
}

func TestServer_GracefulShutdownOrdering(t *testing.T) {
	srv := newTestServer()
	audit := &cloudmodel.MemoryAuditAppender{}
	rt := NewCloudModelRuntime(audit)
	srv.AttachCloudModel(rt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rt.Scheduler.Start(ctx)

	ns := cloudmodel.IdempotencyNamespace{
		PrincipalID: "p1",
		Pattern:     "POST /v1/cloud-platforms",
		Key:         "shutdown-key",
	}
	digest, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{
		Body:    json.RawMessage(`{"metadata":{"name":"x"}}`),
		Pattern: ns.Pattern,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res := rt.Idempotency.Reserve(context.Background(), ns, digest, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("owner reserve: %#v", res)
	}

	waiterDone := make(chan cloudmodel.ReserveResult, 1)
	go func() {
		waiterDone <- rt.Idempotency.Reserve(context.Background(), ns, digest, nil)
	}()
	time.Sleep(20 * time.Millisecond)

	// Ordered shutdown: scheduler stops, then in-flight reservations abort and
	// waiters wake; no completed replay record is published.
	if err := srv.Shutdown(2 * time.Second); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	res := <-waiterDone
	if res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("waiter after shutdown abort: %#v", res)
	}
	if st, ok := rt.Idempotency.StateForTest(ns); !ok || st != cloudmodel.ReservationInFlight {
		t.Fatalf("post-shutdown state = %s ok=%v", st, ok)
	}
	if rt.Idempotency.CompletedCountForTest() != 0 {
		t.Fatal("shutdown must not complete idempotency replay records")
	}
	// Scheduler Stop is idempotent after Shutdown; no further transitions occur.
	rt.Scheduler.Stop()
}

func seedPendingParticipation(t *testing.T, store *cloudmodel.Store) {
	t.Helper()
	cp := model.CloudPlatform{
		Metadata: apimeta.ObjectMeta{Name: "plat", UID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		Spec: model.CloudPlatformSpec{OwnerRegistration: model.OwnerRegistration{
			LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
		}},
		Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
	}
	prov := model.CloudProvider{
		Metadata: apimeta.ObjectMeta{Name: "prov", UID: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US"}},
		Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
	}
	if _, prob := store.CreateCloudPlatform(cp); prob != nil {
		t.Fatalf("platform: %#v", prob)
	}
	if _, prob := store.CreateCloudProvider(prov); prob != nil {
		t.Fatalf("provider: %#v", prob)
	}
	part := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part", UID: "cccccccccccccccccccccccccccccccc", ResourceVersion: "1",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "plat", UID: cp.Metadata.UID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform, Name: "plat", UID: cp.Metadata.UID},
			CloudProviderRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider, Name: "prov", UID: prov.Metadata.UID},
			Environment:      model.ParticipationEnvironmentDevelopment,
		},
		Status: cloudmodel.InitialParticipationStatus(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
	}
	if _, prob := store.CreateParticipation(part); prob != nil {
		t.Fatalf("participation: %#v", prob)
	}
}
