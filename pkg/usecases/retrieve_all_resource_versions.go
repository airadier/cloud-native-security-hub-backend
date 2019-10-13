package usecases

import (
	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
)

type RetrieveAllResourceVersions struct {
	ResourceRepository resource.Repository
	ResourceID         string
}

func (useCase *RetrieveAllResourceVersions) Execute() ([]*resource.Resource, error) {
	return useCase.ResourceRepository.FindById(useCase.ResourceID)
}
