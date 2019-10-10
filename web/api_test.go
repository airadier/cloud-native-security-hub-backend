package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("RESOURCES_PATH", "../test/fixtures/resources")
	os.Setenv("VENDOR_PATH", "../test/fixtures/vendors")

	m.Run()
}

func TestRetrieveAllResourcesHandlerReturnsHTTPOk(t *testing.T) {
	testRetrieveAllReturnsHTTPOk(t, "/resources")
}

func TestRetrieveOneResourceHandlerReturnsHTTPOk(t *testing.T) {
	apacheID := "apache"
	testRetrieveAllReturnsHTTPOk(t, "/resources/"+apacheID)
}

func TestRetrieveFalcoRulesForHelmChartHandlerReturnsHTTPOk(t *testing.T) {
	apacheID := "apache"
	testRetrieveAllReturnsHTTPOk(t, "/resources/"+apacheID+"/custom-rules.yaml")
}

func TestRetrieveFalcoRulesForHelmChartVersionHandlerReturnsHTTPOk(t *testing.T) {
	apacheID := "apache"
	testRetrieveAllReturnsHTTPOk(t, "/resources/"+apacheID+"/1.0.0/custom-rules.yaml")
}

func TestRetrieveAllVendorsHandlerReturnsHTTPOk(t *testing.T) {
	testRetrieveAllReturnsHTTPOk(t, "/vendors")
}

func TestRetrieveOneVendorsHandlerReturnsHTTPOk(t *testing.T) {
	testRetrieveAllReturnsHTTPOk(t, "/vendors/apache")
}

func TestRetrieveAllResourcesFromVendorHandlerReturnsHTTPOk(t *testing.T) {
	testRetrieveAllReturnsHTTPOk(t, "/vendors/apache/resources")
}

func testRetrieveAllReturnsHTTPOk(t *testing.T, path string) {
	request, _ := http.NewRequest("GET", path, nil)
	recorder := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestRetrieveAllResourcesHandlerReturnsResourcesSerializedAsJSON(t *testing.T) {
	testRetrieveallSerializedAsJSON(t, "/resources", "../test/fixtures/resources")
}

func TestRetrieveAllVendorsHandlerReturnsResourcesSerializedAsJSON(t *testing.T) {
	testRetrieveallSerializedAsJSON(t, "/vendors", "../test/fixtures/vendors")
}

func testRetrieveallSerializedAsJSON(t *testing.T, urlPath, fixturesPath string) {
	repo, _ := resource.FromPath(fixturesPath)
	resources, _ := repo.FindAll()

	request, _ := http.NewRequest("GET", urlPath, nil)
	recorder := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(recorder, request)

	var result []*resource.Resource
	body, _ := ioutil.ReadAll(recorder.Body)
	json.Unmarshal([]byte(body), &result)
	assert.Equal(t, resources, result)
}

func TestRetrieveAllResourcesHandlerReturnsAJSONResponse(t *testing.T) {
	testRetrieveAllHandlerReturnsAJSONResponse(t, "/resources")
}
func TestRetrieveAllVendorHandlerReturnsAJSONResponse(t *testing.T) {
	testRetrieveAllHandlerReturnsAJSONResponse(t, "/vendors")
}

func testRetrieveAllHandlerReturnsAJSONResponse(t *testing.T, urlPath string) {
	request, _ := http.NewRequest("GET", urlPath, nil)
	recorder := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, "application/json", recorder.HeaderMap["Content-Type"][0])
}

func TestRetrieveFalcoRulesForHelmChartReturnsContent(t *testing.T) {
	apacheID := "apache"
	request, _ := http.NewRequest("GET", "/resources/"+apacheID+"/custom-rules.yaml", nil)

	recorder := httptest.NewRecorder()
	os.Setenv("RESOURCES_PATH", "../test/fixtures/resources")
	os.Setenv("VENDOR_PATH", "../test/fixtures/vendors")
	router := NewRouter()
	router.ServeHTTP(recorder, request)

	expectedResult := `customRules:
  rules-apache.yaml: |-
    - macro: apache_consider_syscalls
      condition: (evt.num < 0)
`
	assert.Equal(t, expectedResult, string(recorder.Body.Bytes()))
}

func TestRetrieveFalcoRulesForHelmChartVersionReturnsContent(t *testing.T) {
	apacheID := "apache"
	request, _ := http.NewRequest("GET", "/resources/"+apacheID+"/1.0.0/custom-rules.yaml", nil)

	recorder := httptest.NewRecorder()
	os.Setenv("RESOURCES_PATH", "../test/fixtures/resources")
	os.Setenv("VENDOR_PATH", "../test/fixtures/vendors")
	router := NewRouter()
	router.ServeHTTP(recorder, request)

	expectedResult := `customRules:
  rules-apache.yaml: |-
    - macro: apache_consider_syscalls
      condition: (evt.num < 0)
`
	assert.Equal(t, expectedResult, string(recorder.Body.Bytes()))
}

func TestRetrieveFalcoRulesForHelmChartReturnsAYAMLResponse(t *testing.T) {
	apacheID := "apache"
	request, _ := http.NewRequest("GET", "/resources/"+apacheID+"/custom-rules.yaml", nil)
	recorder := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, "application/x-yaml", recorder.HeaderMap["Content-Type"][0])
}

func TestRetrieveFalcoRulesForHelmChartVersionReturnsAYAMLResponse(t *testing.T) {
	apacheID := "apache"
	request, _ := http.NewRequest("GET", "/resources/"+apacheID+"/1.0.0/custom-rules.yaml", nil)
	recorder := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, "application/x-yaml", recorder.HeaderMap["Content-Type"][0])
}

func TestLoggerIsLogging(t *testing.T) {
	apacheID := "apache"
	url := "/resources/" + apacheID + "/custom-rules.yaml"
	request, _ := http.NewRequest("GET", url, nil)
	recorder := httptest.NewRecorder()

	buff := &bytes.Buffer{}
	router := NewRouterWithLogger(log.New(buff, "", 0))
	router.ServeHTTP(recorder, request)

	expectedLog := fmt.Sprintf("200 [] GET %s\n", url)
	assert.Equal(t, expectedLog, buff.String())
}

func TestHealthCheckEndpoint(t *testing.T) {
	request, _ := http.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
