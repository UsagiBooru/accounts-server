package request

import (
	"context"
	"errors"
	"strconv"
)

type key int

const Context_user_id key = 1
const Context_user_permission key = 2

// load permission from context
func GetUserPermission(ctx context.Context) (int32, error) {
	v := ctx.Value(Context_user_permission)
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
	v := ctx.Value(Context_user_id)
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

// load context shorthand
func GetHeaders(ctx context.Context) (int32, int32, error) {
	issuerID, err := GetUserID(ctx)
	if err != nil {
		return 0, 0, err
	}
	issuerPermission, err := GetUserPermission(ctx)
	if err != nil {
		return 0, 0, err
	}
	return issuerID, issuerPermission, nil
}
