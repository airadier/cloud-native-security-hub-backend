package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func NewRouter() http.Handler {
	router := httprouter.New()
	registerOn(router, nil)

	return cors.Default().Handler(router)
}

func NewRouterWithLogger(logger *log.Logger) http.Handler {
	router := httprouter.New()
	registerOn(router, logger)

	return cors.Default().Handler(router)
}

func registerOn(router *httprouter.Router, logger *log.Logger) {
	h := NewHandlerRepository(logger)
	router.GET("/resources", h.retrieveAllResourcesHandler)
	router.GET("/resources/:resource", h.retrieveOneResourceHandler)
	router.GET("/resources/:resource/:version/custom-rules.yaml", h.retrieveFalcoRulesForHelmChartVersionHandler)
	router.GET("/resources/:resource/:version", h.retrieveOneResourceVersionHandler)
	//Collision with previous routes!! Hack: we handle specific case when version == "custom-rules-yaml"?
	//router.GET("/resources/:resource/custom-rules.yaml", h.retrieveFalcoRulesForHelmChartHandler)
	router.GET("/vendors", h.retrieveAllVendorsHandler)
	router.GET("/vendors/:vendor", h.retrieveOneVendorsHandler)
	router.GET("/vendors/:vendor/resources", h.retrieveAllResourcesFromVendorHandler)
	router.GET("/health", h.healthCheckHandler)
	router.NotFound = h.notFound()
}
