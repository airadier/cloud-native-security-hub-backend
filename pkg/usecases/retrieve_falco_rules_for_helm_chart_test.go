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
				ID:      "nginx",
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx v1",
				Vendor:  "Nginx",
				Version: "1.0.0",
				Rules: []*resource.FalcoRuleData{
					{Raw: "nginxRuleV1"},
				},
			},
			{
				ID:      "nginx",
				Kind:    resource.FALCO_RULE,
				Name:    "Falco profile for Nginx v2",
				Vendor:  "Nginx",
				Version: "2.0.0",
				Rules: []*resource.FalcoRuleData{
					{Raw: "nginxRuleV2"},
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

func TestReturnsFalcoRulesForHelmChart(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChart{
		ResourceRepository: memoryResourceRepositoryWithRules(),
		ResourceID:         "nginx",
	}

	result, _ := useCase.Execute()
	expected := `customRules:
  rules-nginx.yaml: nginxRuleV2
`

	assert.Equal(t, expected, string(result))
}

func TestFalcoRulesForHelmChartReturnsNotFound(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChart{
		ResourceRepository: memoryResourceRepositoryWithRules(),
		ResourceID:         "notFound",
	}

	_, err := useCase.Execute()

	assert.Error(t, err)
}
