package usecases;

import (
	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dummyResourcesRepositoryForOne struct{}

func (resources *dummyResourcesRepositoryForOne) FindAll() ([]*resource.Resource, error) {
	return []*resource.Resource{
		{
			Kind:   resource.FALCO_RULE,
			Name:   "Falco profile for Nginx",
			Vendor: "Nginx",
		},
		{
			Kind:   "GrafanaDashboard",
			Name:   "Grafana Dashboard for Traefik",
			Vendor: "Traefik",
		},
	}, nil
}

func (resources *dummyResourcesRepositoryForOne) FindById(id string) (*resource.Resource, error) {
	return &resource.Resource{
		Kind:   resource.FALCO_RULE,
		Name:   "Falco profile for Nginx",
		Vendor: "Nginx",
	}, nil
}

func TestReturnsOneResource(t *testing.T) {
	useCase := RetrieveOneResource{
		ResourceRepository: &dummyResourcesRepositoryForOne{},
		Hash:               "bekiisotdwhvmetchrwp",
	}

	res, _ := useCase.Execute()

	assert.Equal(t, &resource.Resource{
		Kind:   resource.FALCO_RULE,
		Name:   "Falco profile for Nginx",
		Vendor: "Nginx",
	}, res)
}

func TestReturnsResourceNotFound(t *testing.T) {
	useCase := RetrieveOneResource{
		ResourceRepository: &dummyResourcesRepositoryForOne{},
		Hash:               "notFound",
	}

	_, err := useCase.Execute()

	assert.Error(t, err)
}
