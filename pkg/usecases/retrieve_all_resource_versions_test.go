package usecases

import (
	"testing"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestReturnsAllResourceVersions(t *testing.T) {
	resourceRepository := resource.NewMemoryRepository(
		[]*resource.Resource{
			&resource.Resource{ID: "nginx", Name: "Falco profile for Nginx v1", Version: "1.0.0"},
			&resource.Resource{ID: "nginx", Name: "Falco profile for Nginx v2", Version: "2.0.0"},
			&resource.Resource{ID: "grafana", Name: "Falco profile for Grafana", Version: "1.0.0"},
		},
	)

	useCase := RetrieveAllResourceVersions{ResourceRepository: resourceRepository, ResourceID: "nginx"}

	resources, _ := useCase.Execute()

	assert.Equal(t, []*resource.Resource{
		{ID: "nginx", Name: "Falco profile for Nginx v2", Version: "2.0.0"},
		{ID: "nginx", Name: "Falco profile for Nginx v1", Version: "1.0.0"},
	}, resources)
}
