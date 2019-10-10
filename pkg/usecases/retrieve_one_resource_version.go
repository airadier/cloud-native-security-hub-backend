package usecases

import (
	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
)

type RetrieveOneResourceVersion struct {
	ResourceRepository resource.Repository
	ResourceID         string
	ResourceVersion    string
}

func (useCase *RetrieveOneResourceVersion) Execute() (res *resource.Resource, err error) {
	return useCase.ResourceRepository.FindByIdAndVersion(useCase.ResourceID, useCase.ResourceVersion)
}
