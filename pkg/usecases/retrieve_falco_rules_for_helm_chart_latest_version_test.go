package usecases

import (
	"testing"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func memoryResourceRepositoryWithRules() resource.Repository {
	return resource.NewMemoryRepository(
		[]*resource.Resource{
			{
				ID:     "nginx",
				Kind:   resource.FALCO_RULE,
				Name:   "Falco profile for Nginx",
				Vendor: "Nginx",
				Rules: []*resource.FalcoRuleData{
					{Version: "0.0.1", Raw: "nginxRule1"},
					{Version: "0.0.2", Raw: "nginxRule2"},
				},
			},
			{
				ID:            "traefik",
				Kind:          "GrafanaDashboard",
				Name:          "Grafana Dashboard for Traefik",
				Vendor:        "Traefik",
				LatestVersion: "1.0.1",
				Rules: []*resource.FalcoRuleData{
					{Version: "1.0.1", Raw: "traefikRule1.0.1"},
					{Version: "1.0.0", Raw: "traefikRule1.0.0"},
				},
			},
		},
	)
}

func TestReturnsFalcoRulesForHelmChartLatestVersion(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChartLatestVersion{
		ResourceRepository: memoryResourceRepositoryWithRules(),
		ResourceID:         "traefik",
	}

	result, _ := useCase.Execute()
	expected := `customRules:
  rules-traefik.yaml: traefikRule1.0.1
`

	assert.Equal(t, expected, string(result))
}

func TestReturnsFalcoRulesForHelmChartLatestVersionIfNotDefined(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChartLatestVersion{
		ResourceRepository: memoryResourceRepositoryWithRules(),
		ResourceID:         "nginx",
	}

	result, _ := useCase.Execute()
	expected := `customRules:
  rules-nginx.yaml: nginxRule2
`

	assert.Equal(t, expected, string(result))
}

func TestFalcoRulesForHelmChartLatestVersionReturnsNotFound(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChartLatestVersion{
		ResourceRepository: memoryResourceRepositoryWithRules(),
		ResourceID:         "notFound",
	}

	_, err := useCase.Execute()

	assert.Error(t, err)
}
