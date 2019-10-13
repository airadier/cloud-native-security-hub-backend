package usecases

import (
	"testing"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/falcosecurity/cloud-native-security-hub/pkg/vendor"
	"github.com/stretchr/testify/assert"
)

func memoryResourceRepositoryFromVendor() resource.Repository {
	return resource.NewMemoryRepository(
		[]*resource.Resource{
			{
				Name:    "Falco profile for Nginx v1",
				Vendor:  "Nginx",
				Version: "1.0.0",
			},
			{
				Name:    "Falco profile for Nginx v2",
				Vendor:  "Nginx",
				Version: "2.0.0",
			},
			{
				Name:    "Grafana Dashboard for Traefik v1",
				Vendor:  "Traefik",
				Version: "1.0.0",
			},
		},
	)
}

func memoryVendorRepositoryFromVendor() vendor.Repository {
	return vendor.NewMemoryRepository(
		[]*vendor.Vendor{
			{
				ID:   "apache",
				Name: "Apache",
			},
			{
				ID:   "nginx",
				Name: "Nginx",
			},
			{
				ID:   "traefik",
				Name: "Traefik",
			},
		},
	)
}

func TestReturnsAllResourcesFromVendor(t *testing.T) {
	useCase := RetrieveAllResourcesFromVendor{
		VendorID:           "Nginx",
		ResourceRepository: memoryResourceRepositoryFromVendor(),
		VendorRepository:   memoryVendorRepositoryFromVendor(),
	}

	resources, _ := useCase.Execute()

	assert.Equal(t, []*resource.Resource{
		{
			Name:    "Falco profile for Nginx v2",
			Vendor:  "Nginx",
			Version: "2.0.0",
		},
	}, resources)
}

func TestReturnsVendorNotFoundResourcesFromVendor(t *testing.T) {
	useCase := RetrieveAllResourcesFromVendor{
		VendorID:           "not-found",
		ResourceRepository: memoryResourceRepositoryFromVendor(),
		VendorRepository:   memoryVendorRepositoryFromVendor(),
	}

	_, err := useCase.Execute()

	assert.Error(t, err)
}

func TestReturnsResourcesNotFoundResourcesFromVendor(t *testing.T) {
	useCase := RetrieveAllResourcesFromVendor{
		VendorID:           "apache",
		ResourceRepository: memoryResourceRepositoryFromVendor(),
		VendorRepository:   memoryVendorRepositoryFromVendor(),
	}

	_, err := useCase.Execute()

	assert.Error(t, err) //vendor exists but has no resources
}
