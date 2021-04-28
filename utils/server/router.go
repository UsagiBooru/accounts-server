package server

import (
	"context"
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils/request"
	"github.com/gorilla/mux"
)

// middleware to set context
func injectHeaderToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("x-consumer-user-id")
		userPermission := r.Header.Get("x-consumer-user-permission")
		Debug("User id is: " + string(userID))
		Debug("User permission is: " + string(userPermission))
		ctx := context.WithValue(r.Context(), request.CtxUserId, userID)
		ctx = context.WithValue(ctx, request.CtxUserPermission, userPermission)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

// NewRouterWithInject creates a new router with inject header middleware
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
