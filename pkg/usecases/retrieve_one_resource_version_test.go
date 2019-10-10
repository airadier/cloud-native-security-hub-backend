package usecases

import (
	"testing"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func memoryResourceRepositoryWithVersions() resource.Repository {
	return resource.NewMemoryRepository(
		[]*resource.Resource{
			&resource.Resource{
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx v1.0.0",
				Vendor:  "Nginx",
				ID:      "nginx",
				Version: "1.0.0",
			},
			&resource.Resource{
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx v1.0.1",
				Vendor:  "Nginx",
				ID:      "nginx",
				Version: "1.0.1",
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

func TestReturnsOneResourceVersion(t *testing.T) {
	useCase := RetrieveOneResourceVersion{
		ResourceRepository: memoryResourceRepositoryWithVersions(),
		ResourceID:         "nginx",
		ResourceVersion:    "1.0.0",
	}

	res, _ := useCase.Execute()

	assert.Equal(t, &resource.Resource{
		Kind:    resource.FALCO_RULE,
		Name:    "Falco profile for Nginx v1.0.0",
		Vendor:  "Nginx",
		ID:      "nginx",
		Version: "1.0.0",
	}, res)
}

func TestReturnsResourceVersionNotFound(t *testing.T) {
	useCase := RetrieveOneResourceVersion{
		ResourceRepository: memoryResourceRepositoryWithVersions(),
		ResourceID:         "nginx",
		ResourceVersion:    "notFound",
	}

	_, err := useCase.Execute()

	assert.Error(t, err)
}
