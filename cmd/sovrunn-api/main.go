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

	logger := log.New(os.Stdout, "", log.LstdFlags)

	orgRegistry := registry.NewOrganizationRegistry()
	ouRegistry := registry.NewOrganizationUnitRegistry()
	ouBlocker := registry.NewOUChildBlockerChecker(ouRegistry)

	tenantRegistry := registry.NewTenantRegistry()
	tenantBlocker := registry.NewTenantChildBlockerChecker(tenantRegistry)
	projectRegistry := registry.NewProjectRegistry()
	projectBlocker := registry.NewProjectChildBlockerChecker(projectRegistry)
	operationRegistry := registry.NewOperationRegistry()
	emitter := api.NewRegistryEmitter(operationRegistry, nil)

	serviceClassRegistry := registry.NewServiceClassRegistry()
	servicePlanRegistry := registry.NewServicePlanRegistry()
	serviceClassBlocker := registry.NewServicePlanChildBlockerChecker(servicePlanRegistry)

	pluginRegistry := registry.NewPluginRegistry()
	capabilityRegistry := registry.NewCapabilityRegistry()
	pluginBlocker := registry.NewCapabilityChildBlockerChecker(capabilityRegistry)

	serviceInstanceRegistry := registry.NewServiceInstanceRegistry()
	serviceBindingRegistry := registry.NewServiceBindingRegistry()
	capabilityLookup := registry.NewCapabilityLookup(capabilityRegistry)
	planInstanceBlocker := registry.NewServiceInstancePlanBlockerChecker(serviceInstanceRegistry)
	projectInstanceBlocker := registry.NewServiceInstanceProjectBlockerChecker(serviceInstanceRegistry)

	orgHandler := api.NewOrgHandler(orgRegistry, ouBlocker, emitter)
	ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, tenantBlocker, emitter)
	tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry, projectBlocker, emitter)
	projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry, emitter, projectInstanceBlocker)
	operationHandler := api.NewOperationHandler(operationRegistry)
	serviceClassHandler := api.NewServiceClassHandler(serviceClassRegistry, serviceClassBlocker, emitter)
	servicePlanHandler := api.NewServicePlanHandler(servicePlanRegistry, serviceClassRegistry, emitter, planInstanceBlocker)
	pluginHandler := api.NewPluginHandler(pluginRegistry, serviceClassRegistry, pluginBlocker, emitter)
	capabilityHandler := api.NewCapabilityHandler(capabilityRegistry, pluginRegistry, serviceClassRegistry, emitter)
	serviceInstanceHandler := api.NewServiceInstanceHandler(
		serviceInstanceRegistry,
		orgRegistry,
		ouRegistry,
		tenantRegistry,
		projectRegistry,
		serviceClassRegistry,
		servicePlanRegistry,
		capabilityLookup,
		serviceBindingRegistry,
		emitter,
		logger,
	)
	serviceBindingHandler := api.NewServiceBindingHandler(
		serviceBindingRegistry,
		serviceInstanceRegistry,
		emitter,
	)

	readiness := &health.ReadinessState{}
	bootstrapHandler := api.NewBootstrapHandler(cfg, readiness)
	srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler, operationHandler, serviceClassHandler, servicePlanHandler, pluginHandler, capabilityHandler, serviceInstanceHandler, serviceBindingHandler, bootstrapHandler, readiness)

	if err := srv.Start(); err != nil {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
