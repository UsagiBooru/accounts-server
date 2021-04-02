package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/gorilla/mux"
)

func injectHeaderToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[START] middleware1")
		ctx := context.WithValue(r.Context(), "user-id", "123")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		fmt.Println("[END] middleware1")
	}
}

func NewRouterWithInject(routers ...gen.Router) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, api := range routers {
		for _, route := range api.Routes() {
			var handler http.Handler
			handler = injectHeaderToContext(route.HandlerFunc)
			handler = gen.Logger(handler, route.Name)

			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}

	return router
}
