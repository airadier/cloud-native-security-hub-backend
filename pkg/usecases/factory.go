package usecases

import (
	"log"
	"os"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/falcosecurity/cloud-native-security-hub/pkg/vendor"
)

type Factory interface {
	NewRetrieveAllResourcesUseCase() *RetrieveAllResources
	NewRetrieveAllResourcesLatestVersionsUseCase() *RetrieveAllResourcesLatestVersions
	NewRetrieveAllResourceVersionsUseCase(resourceID string) *RetrieveAllResourceVersions
	NewRetrieveOneResourceUseCase(resourceID string) *RetrieveOneResource
	NewRetrieveOneResourceVersionUseCase(resourceID, resourceVersion string) *RetrieveOneResourceVersion
	NewRetrieveFalcoRulesForHelmChartUseCase(resourceID string) *RetrieveFalcoRulesForHelmChart
	NewRetrieveFalcoRulesForHelmChartVersionUseCase(resourceID, resourceVersion string) *RetrieveFalcoRulesForHelmChartVersion
	NewRetrieveAllVendorsUseCase() *RetrieveAllVendors
	NewRetrieveOneVendorUseCase(vendorID string) *RetrieveOneVendor
	NewRetrieveAllResourcesFromVendorUseCase(vendorID string) *RetrieveAllResourcesFromVendor

	NewResourcesRepository() resource.Repository
	NewVendorRepository() vendor.Repository
}

func NewFileFactory() Factory {
	factory := &fileFactory{}
	factory.resourceRepository = factory.NewResourcesRepository()
	factory.vendorRepository = factory.NewVendorRepository()
	return factory
}

func NewDBFactory() Factory {
	factory := &dbFactory{}
	factory.resourceRepository = factory.NewResourcesRepository()
	factory.vendorRepository = factory.NewVendorRepository()
	return factory
}

type factory struct {
	vendorRepository   vendor.Repository
	resourceRepository resource.Repository
}

type fileFactory struct {
	factory
}

type dbFactory struct {
	factory
}

func (f *factory) NewRetrieveAllResourcesUseCase() *RetrieveAllResources {
	return &RetrieveAllResources{
		ResourceRepository: f.resourceRepository,
	}
}

func (f *factory) NewRetrieveAllResourcesLatestVersionsUseCase() *RetrieveAllResourcesLatestVersions {
	return &RetrieveAllResourcesLatestVersions{
		ResourceRepository: f.resourceRepository,
	}
}

func (f *factory) NewRetrieveAllResourceVersionsUseCase(resourceID string) *RetrieveAllResourceVersions {
	return &RetrieveAllResourceVersions{
		ResourceRepository: f.resourceRepository,
		ResourceID:         resourceID,
	}
}

func (f *factory) NewRetrieveOneResourceUseCase(resourceID string) *RetrieveOneResource {
	return &RetrieveOneResource{
		ResourceRepository: f.resourceRepository,
		ResourceID:         resourceID,
	}
}

func (f *factory) NewRetrieveOneResourceVersionUseCase(resourceID, resourceVersion string) *RetrieveOneResourceVersion {
	return &RetrieveOneResourceVersion{
		ResourceRepository: f.resourceRepository,
		ResourceID:         resourceID,
		ResourceVersion:    resourceVersion,
	}
}

func (f *factory) NewRetrieveFalcoRulesForHelmChartUseCase(resourceID string) *RetrieveFalcoRulesForHelmChart {
	return &RetrieveFalcoRulesForHelmChart{
		ResourceRepository: f.resourceRepository,
		ResourceID:         resourceID,
	}
}

func (f *factory) NewRetrieveFalcoRulesForHelmChartVersionUseCase(resourceID, resourceVersion string) *RetrieveFalcoRulesForHelmChartVersion {
	return &RetrieveFalcoRulesForHelmChartVersion{
		ResourceRepository: f.resourceRepository,
		ResourceID:         resourceID,
		ResourceVersion:    resourceVersion,
	}
}

func (f *factory) NewRetrieveAllVendorsUseCase() *RetrieveAllVendors {
	return &RetrieveAllVendors{
		VendorRepository: f.vendorRepository,
	}
}

func (f *factory) NewRetrieveOneVendorUseCase(vendorID string) *RetrieveOneVendor {
	return &RetrieveOneVendor{
		VendorRepository: f.vendorRepository,
		VendorID:         vendorID,
	}
}

func (f *factory) NewRetrieveAllResourcesFromVendorUseCase(vendorID string) *RetrieveAllResourcesFromVendor {
	return &RetrieveAllResourcesFromVendor{
		VendorID:           vendorID,
		VendorRepository:   f.vendorRepository,
		ResourceRepository: f.resourceRepository,
	}
}

func (f *fileFactory) NewResourcesRepository() resource.Repository {
	resourcesPath, ok := os.LookupEnv("RESOURCES_PATH")
	if !ok {
		log.Println("The RESOURCES_PATH env var is not set")
		os.Exit(1)
	}
	repo, err := resource.FromPath(resourcesPath)
	if err != nil {
		log.Println("the resource repository of type file does not exist")
		os.Exit(1)
	}
	return repo
}

func (f *fileFactory) NewVendorRepository() vendor.Repository {
	vendorPath, ok := os.LookupEnv("VENDOR_PATH")
	if !ok {
		log.Println("The VENDOR_PATH env var is not set")
		os.Exit(1)
	}
	repo, err := vendor.FromPath(vendorPath)
	if err != nil {
		log.Println("the resource repository of type file does not exist")
		os.Exit(1)
	}
	return repo
}

func (f *dbFactory) NewResourcesRepository() resource.Repository {
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Println("The DB_HOST env var is not set")
		os.Exit(1)
	}

	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		log.Println("The DB_PORT env var is not set")
		os.Exit(1)
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Println("The DB_NAME env var is not set")
		os.Exit(1)
	}

	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Println("The DB_USER env var is not set")
		os.Exit(1)
	}

	dbPass, ok := os.LookupEnv("DB_PASS")
	if !ok {
		log.Println("The DB_PASS env var is not set")
		os.Exit(1)
	}

	repo, err := resource.FromPostgres(dbHost, dbPort, dbName, dbUser, dbPass)
	if err != nil {
		log.Printf("error creating DB resources repo: %s\n", err)
		os.Exit(1)
	}

	return repo
}

func (f *dbFactory) NewVendorRepository() vendor.Repository {
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Println("The DB_HOST env var is not set")
		os.Exit(1)
	}

	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		log.Println("The DB_PORT env var is not set")
		os.Exit(1)
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Println("The DB_NAME env var is not set")
		os.Exit(1)
	}

	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Println("The DB_USER env var is not set")
		os.Exit(1)
	}

	dbPass, ok := os.LookupEnv("DB_PASS")
	if !ok {
		log.Println("The DB_PASS env var is not set")
		os.Exit(1)
	}

	repo, err := vendor.FromPostgres(dbHost, dbPort, dbName, dbUser, dbPass)
	if err != nil {
		log.Printf("error creating DB vendors repo: %s\n", err)
		os.Exit(1)
	}

	return repo
}
