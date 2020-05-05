package rest

import (
	"github.com/go-chi/chi"
	"net/http"
)

func newHealthRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/live", getLive)
	router.Get("/health", getHealth)
	router.Get("/ready", getReady)

	return router
}

func getLive(w http.ResponseWriter, r *http.Request) {
	writeOkResponse(w, getHealthResponseData())
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	writeOkResponse(w, getHealthResponseData())
}

func getReady(w http.ResponseWriter, r *http.Request) {
	writeOkResponse(w, getHealthResponseData())
}
