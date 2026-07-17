package main

import (
	"flag"
	"log"
	"os"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/server"
)

func main() {
	configPath := flag.String("config", "configs/sovrunn-api.local.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("startup error: %v", err)
		os.Exit(1)
	}

	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	ouBlocker := registry.NewOUChildBlockerChecker(ouRegistry)

	tenantRegistry := registry.NewTenantRegistry()
	tenantBlocker := registry.NewTenantChildBlockerChecker(tenantRegistry)

	orgHandler := api.NewOrgHandler(orgRegistry, ouBlocker)
	ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, tenantBlocker)
	tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry)

	readiness := &health.ReadinessState{}
	bootstrapHandler := api.NewBootstrapHandler(cfg, readiness)
	srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, bootstrapHandler, readiness)

	if err := srv.Start(); err != nil {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
