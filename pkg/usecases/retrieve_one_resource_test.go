package usecases

import (
	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func memoryResourceRepository() resource.Repository {
	return resource.NewMemoryRepository(
		[]*resource.Resource{
			&resource.Resource{
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx v0.9.0",
				Vendor:  "Nginx",
				ID:      "nginx",
				Version: "0.9.0",
			},
			&resource.Resource{
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx v1.0.1",
				Vendor:  "Nginx",
				ID:      "nginx",
				Version: "1.0.1",
			},
			&resource.Resource{
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx v1.0.0",
				Vendor:  "Nginx",
				ID:      "nginx",
				Version: "1.0.0",
			},
			&resource.Resource{
				Kind:   resource.FALCO_RULE,
				Name:   "Falco profile for Traefik",
				Vendor: "Traefik",
				ID:     "traefik",
			},
		},
	)
}

func TestReturnsOneResource(t *testing.T) {
	useCase := RetrieveOneResource{
		ResourceRepository: memoryResourceRepository(),
		ResourceID:         "nginx",
	}

	res, _ := useCase.Execute()

	assert.Equal(t, &resource.Resource{
		Kind:    resource.FALCO_RULE,
		Name:    "Falco profile for Nginx v1.0.1",
		Vendor:  "Nginx",
		ID:      "nginx",
		Version: "1.0.1",
	}, res)
}

func TestReturnsResourceNotFound(t *testing.T) {
	useCase := RetrieveOneResource{
		ResourceRepository: memoryResourceRepository(),
		ResourceID:         "notFound",
	}

	_, err := useCase.Execute()

	assert.Error(t, err)
}
