package server

import (
	"net"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
	"github.com/sanjeevksaini/sovrunn/internal/registry"
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
	orgHandler := api.NewOrgHandler(orgRegistry, registry.NoopChildBlockerChecker{})
	ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, nil)
	bootstrap := api.NewBootstrapHandler(cfg, readiness)
	srv := New(cfg, orgHandler, ouHandler, bootstrap, readiness)

	if err := srv.Start(); err == nil {
		t.Fatal("Start() expected error when port is already in use")
	}
	if readiness.IsReady() {
		t.Fatal("readiness should remain false when listener bind fails")
	}
}
