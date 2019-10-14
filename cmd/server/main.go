package main

import (
	"log"
	"net/http"
	"os"

	"github.com/falcosecurity/cloud-native-security-hub/web"
)

func main() {
	router := web.NewDBRouterWithLogger(log.New(os.Stderr, "", log.Ltime|log.Ldate|log.LUTC))

	log.Fatal(http.ListenAndServe(":8080", router))
}
