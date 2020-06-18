package rest

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/vagner-nascimento/go-adp-bridge/config"
	"github.com/vagner-nascimento/go-adp-bridge/src/apperror"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"net/http"
	"strings"
)

func StartRestServer() <-chan error {
	errsCh := make(chan error)

	go func() {
		defer close(errsCh)

		router := chi.NewRouter()
		router.Use(getMiddlewareList()...)
		router.Route("/", func(r chi.Router) {
			r.Mount("/", newHealthRoutes())
		})

		var availableRoutes string
		walkThroughRoutes := func(
			method string,
			route string,
			handler http.Handler,
			middleware ...func(http.Handler) http.Handler,
		) error {
			availableRoutes += fmt.Sprintf("\n%s %s", method, strings.Replace(route, "/*/", "/", -1))

			return nil
		}

		if err := chi.Walk(router, walkThroughRoutes); err != nil {
			logger.Error("error on verify rest routes", err)
			errsCh <- apperror.New("an error occurred on try to start rest server", err, nil)
		} else {
			logger.Info("available rest routes:", availableRoutes)

			port := config.Get().Presentation.Web.Port
			logger.Info("started rest server at port", port)

			errsCh <- http.ListenAndServe(fmt.Sprintf(":%d", port), router)
		}
	}()

	return errsCh
}
