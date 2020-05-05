package rest

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"net/http"
)

func StartRestServer(errsCh chan error) {
	router := chi.NewRouter()

	router.Use(getMiddlewareList()...)
	router.Route("/", func(r chi.Router) {
		r.Mount("/", newHealthRoutes())
	})

	walkThroughRoutes := func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
		return nil
	}

	if err := chi.Walk(router, walkThroughRoutes); err != nil {
		logger.Error("error on verify http routes", err)
		errsCh <- errors.New("an error occurred on try to start http server")
	} else {
		port := config.Get().Presentation.Web.Port
		logger.Info("started rest server at port", port)
		errsCh <- http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	}

	close(errsCh)
}
