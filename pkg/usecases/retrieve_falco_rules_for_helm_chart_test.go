package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReturnsFalcoRulesForHelmChart(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChart{
		ResourceRepository: memoryResourceRepositoryWithRules(),
		ResourceID:         "nginx",
		RuleVersion:        "0.0.1",
	}

	result, _ := useCase.Execute()
	expected := `customRules:
  rules-nginx.yaml: nginxRule1
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

func TestFalcoRulesForHelmChartWrongVersionReturnsNotFound(t *testing.T) {
	useCase := RetrieveFalcoRulesForHelmChart{
		ResourceRepository: memoryResourceRepositoryWithRules(),
		ResourceID:         "nginx",
		RuleVersion:        "notfound",
	}

	_, err := useCase.Execute()

	assert.Error(t, err)
}
