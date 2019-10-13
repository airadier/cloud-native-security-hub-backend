package usecases

import (
	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
)

type RetrieveAllResourcesLatestVersions struct {
	ResourceRepository resource.Repository
}

func (useCase *RetrieveAllResourcesLatestVersions) Execute() ([]*resource.Resource, error) {
	return useCase.ResourceRepository.FindAllLatestVersions()
}
