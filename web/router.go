package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func NewRouter() http.Handler {
	router := httprouter.New()
	registerOn(NewHandlerFileRepository(nil), router, nil)

	return cors.Default().Handler(router)
}

func NewRouterWithLogger(logger *log.Logger) http.Handler {
	router := httprouter.New()

	registerOn(NewHandlerFileRepository(logger), router, logger)

	return cors.Default().Handler(router)
}

func NewDBRouterWithLogger(logger *log.Logger) http.Handler {
	router := httprouter.New()

	registerOn(NewHandlerDBRepository(logger), router, logger)

	return cors.Default().Handler(router)
}

func registerOn(h HandlerRepository, router *httprouter.Router, logger *log.Logger) {
	router.GET("/resources", h.retrieveAllResourcesLatestVersionsHandler)
	router.GET("/resources/:resource", h.retrieveOneResourceHandler)
	router.GET("/resources/:resource/custom-rules.yaml", h.retrieveFalcoRulesForHelmChartHandler)
	router.GET("/resources/:resource/versions", h.retrieveAllResourceVersionsHandler)
	router.GET("/resources/:resource/versions/:version/custom-rules.yaml", h.retrieveFalcoRulesForHelmChartVersionHandler)
	router.GET("/resources/:resource/versions/:version", h.retrieveOneResourceVersionHandler)
	router.GET("/vendors", h.retrieveAllVendorsHandler)
	router.GET("/vendors/:vendor", h.retrieveOneVendorsHandler)
	router.GET("/vendors/:vendor/resources", h.retrieveAllResourcesFromVendorHandler)
	router.GET("/health", h.healthCheckHandler)
	router.NotFound = h.notFound()
}
