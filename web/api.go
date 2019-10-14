package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/usecases"
	"github.com/julienschmidt/httprouter"
)

type HandlerRepository interface {
	notFound() http.HandlerFunc
	retrieveAllResourcesLatestVersionsHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params)
	retrieveAllResourceVersionsHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params)
	retrieveOneResourceHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	retrieveOneResourceVersionHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	retrieveFalcoRulesForHelmChartHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	retrieveFalcoRulesForHelmChartVersionHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	retrieveAllVendorsHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params)
	retrieveOneVendorsHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	retrieveAllResourcesFromVendorHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	healthCheckHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type handlerRepository struct {
	factory usecases.Factory
	logger  *log.Logger
}

func NewHandlerRepository(logger *log.Logger) HandlerRepository {
	return &handlerRepository{
		factory: usecases.NewDBFactory(),
		logger:  logger,
	}
}

func (h *handlerRepository) notFound() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		h.logRequest(request, 404)
		http.NotFound(writer, request)
	}

}
func (h *handlerRepository) logRequest(request *http.Request, statusCode int) {
	if h.logger == nil {
		return
	}

	line := fmt.Sprintf("%d [%s] %s %s", statusCode, request.RemoteAddr, request.Method, request.URL)
	h.logger.Println(line)
}

func (h *handlerRepository) retrieveAllResourcesLatestVersionsHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	useCase := h.factory.NewRetrieveAllResourcesLatestVersionsUseCase()
	resources, err := useCase.Execute()
	if err != nil {
		h.logRequest(request, 500)
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")

	h.logRequest(request, 200)
	json.NewEncoder(writer).Encode(resources)
}

func (h *handlerRepository) retrieveAllResourceVersionsHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	useCase := h.factory.NewRetrieveAllResourceVersionsUseCase(params.ByName("resource"))
	resources, err := useCase.Execute()
	if err != nil {
		h.logRequest(request, 500)
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")

	h.logRequest(request, 200)
	json.NewEncoder(writer).Encode(resources)
}

func (h *handlerRepository) retrieveOneResourceHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	useCase := h.factory.NewRetrieveOneResourceUseCase(params.ByName("resource"))
	resources, err := useCase.Execute()
	if err != nil {
		writer.WriteHeader(500)
		h.logRequest(request, 500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	h.logRequest(request, 200)
	json.NewEncoder(writer).Encode(resources)
}

func (h *handlerRepository) retrieveOneResourceVersionHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	version := params.ByName("version")
	if version == "custom-rules.yaml" {
		// Hack for route collision in router.go
		h.retrieveFalcoRulesForHelmChartHandler(writer, request, params)
		return
	}

	useCase := h.factory.NewRetrieveOneResourceVersionUseCase(params.ByName("resource"), version)
	resources, err := useCase.Execute()
	if err != nil {
		writer.WriteHeader(500)
		h.logRequest(request, 500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	h.logRequest(request, 200)
	json.NewEncoder(writer).Encode(resources)
}

func (h *handlerRepository) retrieveFalcoRulesForHelmChartHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	useCase := h.factory.NewRetrieveFalcoRulesForHelmChartUseCase(params.ByName("resource"))
	content, err := useCase.Execute()
	if err != nil {
		writer.WriteHeader(500)
		h.logRequest(request, 500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/x-yaml")
	h.logRequest(request, 200)
	writer.Write(content)
}

func (h *handlerRepository) retrieveFalcoRulesForHelmChartVersionHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	useCase := h.factory.NewRetrieveFalcoRulesForHelmChartVersionUseCase(params.ByName("resource"), params.ByName("version"))
	content, err := useCase.Execute()
	if err != nil {
		writer.WriteHeader(500)
		h.logRequest(request, 500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/x-yaml")
	h.logRequest(request, 200)
	writer.Write(content)
}

func (h *handlerRepository) retrieveAllVendorsHandler(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	useCase := h.factory.NewRetrieveAllVendorsUseCase()
	resources, err := useCase.Execute()
	if err != nil {
		writer.WriteHeader(500)
		h.logRequest(request, 500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	h.logRequest(request, 200)
	json.NewEncoder(writer).Encode(resources)
}

func (h *handlerRepository) retrieveOneVendorsHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	useCase := h.factory.NewRetrieveOneVendorUseCase(params.ByName("vendor"))
	resources, err := useCase.Execute()
	if err != nil {
		writer.WriteHeader(500)
		h.logRequest(request, 500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	h.logRequest(request, 200)
	json.NewEncoder(writer).Encode(resources)
}

func (h *handlerRepository) retrieveAllResourcesFromVendorHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	useCase := h.factory.NewRetrieveAllResourcesFromVendorUseCase(params.ByName("vendor"))
	resources, err := useCase.Execute()
	if err != nil {
		writer.WriteHeader(500)
		h.logRequest(request, 500)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	h.logRequest(request, 200)
	json.NewEncoder(writer).Encode(resources)
}

func (h *handlerRepository) healthCheckHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	h.logRequest(request, 200)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte("OK"))
}
