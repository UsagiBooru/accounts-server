package internal

import (
	"context"
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/gorilla/mux"
)

// middleware to set context
func injectHeaderToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id := r.Header.Get("x-consumer-user-id")
		Debug("User id is: " + string(user_id))
		user_permission := r.Header.Get("x-consumer-user-permission")
		Debug("User permission is: " + string(user_permission))
		ctx := context.WithValue(r.Context(), request.context_user_id, user_id)
		ctx = context.WithValue(ctx, request.context_user_permission, user_permission)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

// router creater for use middleware
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
