package usecases

import "github.com/falcosecurity/cloud-native-security-hub/pkg/resource"

type RetrieveFalcoRulesForHelmChartVersion struct {
	ResourceRepository resource.Repository
	ResourceID         string
	ResourceVersion    string
}

func (useCase *RetrieveFalcoRulesForHelmChartVersion) Execute() ([]byte, error) {
	res, err := useCase.ResourceRepository.FindByIdAndVersion(useCase.ResourceID, useCase.ResourceVersion)
	if err != nil {
		return nil, err
	}
	return res.GenerateRulesForHelmChart(), nil
}
