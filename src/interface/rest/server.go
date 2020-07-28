package restinterface

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/vagner-nascimento/go-enriching-adp/config"
	"github.com/vagner-nascimento/go-enriching-adp/src/apperror"
	"github.com/vagner-nascimento/go-enriching-adp/src/infra/logger"
	"net/http"
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

		walkThroughRoutes := func(
			method string,
			route string,
			handler http.Handler,
			middleware ...func(http.Handler) http.Handler,
		) error {
			return nil
		}

		if err := chi.Walk(router, walkThroughRoutes); err != nil {
			logger.Error("error on verify rest routes", err)
			errsCh <- apperror.New("an error occurred on try to start rest server", err, nil)
		} else {
			port := config.Get().Presentation.Web.Port
			logger.Info("started rest server at port", port)

			errsCh <- http.ListenAndServe(fmt.Sprintf(":%d", port), router)
		}
	}()

	return errsCh
}
