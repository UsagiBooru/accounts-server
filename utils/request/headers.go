package request

import (
	"context"
	"errors"
	"strconv"
)

type key int

const CtxUserId key = 1
const CtxUserPermission key = 2

// GetUserPermission gets a requested user's permission from context
func GetUserPermission(ctx context.Context) (int32, error) {
	v := ctx.Value(CtxUserPermission)
	permission, ok := v.(string)
	if !ok {
		return 0, errors.New("could not parse permission header")
	}
	permissionNumber, err := strconv.Atoi(permission)
	if err != nil {
		return 0, errors.New("could not parse permission header")
	}
	return int32(permissionNumber), nil
}

// GetUserID gets a requested user's id from context
func GetUserID(ctx context.Context) (int32, error) {
	v := ctx.Value(CtxUserId)
	userID, ok := v.(string)
	if !ok {
		return 0, errors.New("could not parse user id header")
	}
	userIDNumber, err := strconv.Atoi(userID)
	if err != nil {
		return 0, errors.New("could not parse user id header")
	}
	return int32(userIDNumber), nil
}

// GetHeaders gets requested user's permission and id from context
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
