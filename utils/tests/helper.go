package tests

import (
	"net/http"
	"strconv"

	"github.com/UsagiBooru/accounts-server/models/const_models/account_const"
)

// SetAdminUserHeader set requested user as ID:1 and permission:9
func SetAdminUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "1")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(int(account_const.PERMISSION_ADMIN)))
	return req
}

// SetModUserHeader set requested user as ID:2 and permission:5
func SetModUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "2")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(int(account_const.PERMISSION_MOD)))
	return req
}

// SetNormalUserHeader set requested user as ID:3 and permission:0
func SetNormalUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "3")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(int(account_const.PERMISSION_USER)))
	return req
}
