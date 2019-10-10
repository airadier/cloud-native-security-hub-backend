package usecases

import (
	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/falcosecurity/cloud-native-security-hub/pkg/vendor"
	"log"
	"os"
)

type Factory interface {
	NewRetrieveAllResourcesUseCase() *RetrieveAllResources
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

func NewFactory() Factory {
	factory := &factory{}
	factory.resourceRepository = factory.NewResourcesRepository()
	factory.vendorRepository = factory.NewVendorRepository()
	return factory
}

type factory struct {
	vendorRepository   vendor.Repository
	resourceRepository resource.Repository
}

func (f *factory) NewRetrieveAllResourcesUseCase() *RetrieveAllResources {
	return &RetrieveAllResources{
		ResourceRepository: f.resourceRepository,
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

func (f *factory) NewResourcesRepository() resource.Repository {
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

func (f *factory) NewVendorRepository() vendor.Repository {
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
