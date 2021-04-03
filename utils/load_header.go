package utils

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/gorilla/mux"
)

type key int

const context_user_id key = 1
const context_user_permission key = 2

// load permission from context
func GetUserPermission(ctx context.Context) (int32, error) {
	v := ctx.Value(context_user_permission)
	permission, ok := v.(string)
	if !ok {
		return 0, errors.New("could not parse permission header")
	}
	permission_num, err := strconv.Atoi(permission)
	if err != nil {
		return 0, errors.New("could not parse permission header")
	}
	return int32(permission_num), nil
}

// load user id from context
func GetUserID(ctx context.Context) (int32, error) {
	v := ctx.Value(context_user_id)
	user_id, ok := v.(string)
	if !ok {
		return 0, errors.New("could not parse user id header")
	}
	user_id_num, err := strconv.Atoi(user_id)
	if err != nil {
		return 0, errors.New("could not parse user id header")
	}
	return int32(user_id_num), nil
}

// middleware to set context
func injectHeaderToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id := r.Header.Get("x-consumer-user-id")
		Debug("User id is: " + string(user_id))
		user_permission := r.Header.Get("x-consumer-user-permission")
		Debug("User permission is: " + string(user_permission))
		ctx := context.WithValue(r.Context(), context_user_id, user_id)
		ctx = context.WithValue(ctx, context_user_permission, user_permission)
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
