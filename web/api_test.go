package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
	"github.com/stretchr/testify/assert"
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

func TestRetrieveResourceVersionsHandlerReturnsHTTPOk(t *testing.T) {
	apacheID := "apache"
	testRetrieveAllReturnsHTTPOk(t, "/resources/"+apacheID+"/versions")
}

func TestRetrieveOneResourceVersionHandlerReturnsHTTPOk(t *testing.T) {
	apacheID := "apache"
	testRetrieveAllReturnsHTTPOk(t, "/resources/"+apacheID+"/versions/1.0.0")
}

func TestRetrieveFalcoRulesForHelmChartHandlerReturnsHTTPOk(t *testing.T) {
	apacheID := "apache"
	testRetrieveAllReturnsHTTPOk(t, "/resources/"+apacheID+"/custom-rules.yaml")
}

func TestRetrieveFalcoRulesForHelmChartVersionHandlerReturnsHTTPOk(t *testing.T) {
	apacheID := "apache"
	testRetrieveAllReturnsHTTPOk(t, "/resources/"+apacheID+"/versions/1.0.0/custom-rules.yaml")
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

	request, _ := http.NewRequest("GET", "/resources", nil)
	recorder := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(recorder, request)

	expectedResult := `[{"id":"apache","kind":"FalcoRules","vendor":"Apache","name":"Apache","shortDescription":"","description":"# Apache Falco Rules\n","keywords":["web"],"icon":"https://upload.wikimedia.org/wikipedia/commons/thumb/d/db/Apache_HTTP_server_logo_%282016%29.svg/300px-Apache_HTTP_server_logo_%282016%29.svg.png","website":"","maintainers":[{"name":"nestorsalceda","email":"nestor.salceda@sysdig.com"},{"name":"fedebarcelona","email":"fede.barcelona@sysdig.com"}],"version":"1.0.1","rules":[{"raw":"- macro: apache_consider_syscalls\n  condition: (evt.num \u003c 0)\n  version: 1.0.1"}]},{"id":"mongodb","kind":"FalcoRules","vendor":"Mongo","name":"MongoDB","shortDescription":"","description":"# MongoDB Falco Rules\n","keywords":["database"],"icon":"https://upload.wikimedia.org/wikipedia/en/thumb/4/45/MongoDB-Logo.svg/2560px-MongoDB-Logo.svg.png","website":"","maintainers":[{"name":"nestorsalceda","email":"nestor.salceda@sysdig.com"},{"name":"fedebarcelona","email":"fede.barcelona@sysdig.com"}],"version":"1.0.0","rules":[{"raw":"- macro: mongo_consider_syscalls\n  condition: (evt.num \u003c 0)"}]}]
`
	assert.Equal(t, expectedResult, string(recorder.Body.Bytes()))
}

func TestRetrieveAllResourceVersionsHandlerReturnsResourcesSerializedAsJSON(t *testing.T) {
	apacheID := "apache"
	request, _ := http.NewRequest("GET", "/resources/"+apacheID+"/versions", nil)
	recorder := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(recorder, request)

	expectedResult := `[{"id":"apache","kind":"FalcoRules","vendor":"Apache","name":"Apache","shortDescription":"","description":"# Apache Falco Rules\n","keywords":["web"],"icon":"https://upload.wikimedia.org/wikipedia/commons/thumb/d/db/Apache_HTTP_server_logo_%282016%29.svg/300px-Apache_HTTP_server_logo_%282016%29.svg.png","website":"","maintainers":[{"name":"nestorsalceda","email":"nestor.salceda@sysdig.com"},{"name":"fedebarcelona","email":"fede.barcelona@sysdig.com"}],"version":"1.0.1","rules":[{"raw":"- macro: apache_consider_syscalls\n  condition: (evt.num \u003c 0)\n  version: 1.0.1"}]},{"id":"apache","kind":"FalcoRules","vendor":"Apache","name":"Apache","shortDescription":"","description":"# Apache Falco Rules\n","keywords":["web"],"icon":"https://upload.wikimedia.org/wikipedia/commons/thumb/d/db/Apache_HTTP_server_logo_%282016%29.svg/300px-Apache_HTTP_server_logo_%282016%29.svg.png","website":"","maintainers":[{"name":"nestorsalceda","email":"nestor.salceda@sysdig.com"},{"name":"fedebarcelona","email":"fede.barcelona@sysdig.com"}],"version":"1.0.0","rules":[{"raw":"- macro: apache_consider_syscalls\n  condition: (evt.num \u003c 0)\n  version: 1.0.0"}]}]
`
	assert.Equal(t, expectedResult, string(recorder.Body.Bytes()))
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
      version: 1.0.1
`
	assert.Equal(t, expectedResult, string(recorder.Body.Bytes()))
}

func TestRetrieveFalcoRulesForHelmChartVersionReturnsContent(t *testing.T) {
	apacheID := "apache"
	request, _ := http.NewRequest("GET", "/resources/"+apacheID+"/versions/1.0.0/custom-rules.yaml", nil)

	recorder := httptest.NewRecorder()
	os.Setenv("RESOURCES_PATH", "../test/fixtures/resources")
	os.Setenv("VENDOR_PATH", "../test/fixtures/vendors")
	router := NewRouter()
	router.ServeHTTP(recorder, request)

	expectedResult := `customRules:
  rules-apache.yaml: |-
    - macro: apache_consider_syscalls
      condition: (evt.num < 0)
      version: 1.0.0
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
	request, _ := http.NewRequest("GET", "/resources/"+apacheID+"/versions/1.0.0/custom-rules.yaml", nil)
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
