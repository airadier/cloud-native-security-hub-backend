package usecases

import (
	"testing"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestReturnsAllResources(t *testing.T) {
	resourceRepository := resource.NewMemoryRepository(
		[]*resource.Resource{
			&resource.Resource{Name: "Falco profile for Nginx v1", Version: "1.0.0"},
			&resource.Resource{Name: "Falco profile for Nginx v2", Version: "2.0.0"},
			&resource.Resource{Name: "Falco profile for Grafana", Version: "1.0.0"},
		},
	)

	useCase := RetrieveAllResources{ResourceRepository: resourceRepository}

	resources, _ := useCase.Execute()

	assert.Equal(t, []*resource.Resource{
		{Name: "Falco profile for Nginx v2", Version: "2.0.0"},
		{Name: "Falco profile for Nginx v1", Version: "1.0.0"},
		{Name: "Falco profile for Grafana", Version: "1.0.0"},
	}, resources)
}
