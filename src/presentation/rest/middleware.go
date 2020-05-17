package rest

import "net/http"

func responseHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getMiddlewareList() (middlewareList []func(http.Handler) http.Handler) {
	middlewareList = append(middlewareList, responseHeadersMiddleware)

	return
}
