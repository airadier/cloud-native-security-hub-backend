package usecases

import (
	"testing"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func memoryResourceRepositoryWithRulesVersions() resource.Repository {
	return resource.NewMemoryRepository(
		[]*resource.Resource{
			{
				ID:      "nginx",
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx",
				Vendor:  "Nginx",
				Version: "1.0.0",
				Rules: []*resource.FalcoRuleData{
					{Raw: "nginxRule0"},
				},
			},
			{
				ID:      "nginx",
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx",
				Vendor:  "Nginx",
				Version: "1.0.2",
				Rules: []*resource.FalcoRuleData{
					{Raw: "nginxRule2"},
				},
			},
			{
				ID:      "nginx",
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx",
				Vendor:  "Nginx",
				Version: "1.0.1",
				Rules: []*resource.FalcoRuleData{
					{Raw: "nginxRule1"},
				},
			},
			{
				ID:      "traefik",
				Kind:    "GrafanaDashboard",
				Name:    "Grafana Dashboard for Traefik",
				Vendor:  "Traefik",
				Version: "1.0.0",
				Rules: []*resource.FalcoRuleData{
					{Raw: "traefikRule"},
				},
			},
		},
	)
}

func TestReturnsFalcoRulesForHelmChartVersion(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChartVersion{
		ResourceRepository: memoryResourceRepositoryWithRulesVersions(),
		ResourceID:         "nginx",
		ResourceVersion:    "1.0.1",
	}

	result, _ := useCase.Execute()
	expected := `customRules:
  rules-nginx.yaml: nginxRule1
`

	assert.Equal(t, expected, string(result))
}

func TestFalcoRulesForHelmChartVersionReturnsNotFound(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChartVersion{
		ResourceRepository: memoryResourceRepositoryWithRulesVersions(),
		ResourceID:         "nginx",
		ResourceVersion:    "notfound",
	}

	_, err := useCase.Execute()

	assert.Error(t, err)
}
