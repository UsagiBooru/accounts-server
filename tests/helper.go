package tests

import (
	"net/http"
	"strconv"

	"github.com/UsagiBooru/accounts-server/models/const_models/account_const"
)

func SetAdminUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "1")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(int(account_const.PERMISSION_ADMIN)))
	return req
}

func SetModUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "2")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(int(account_const.PERMISSION_MOD)))
	return req
}

func SetNormalUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "3")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(int(account_const.PERMISSION_USER)))
	return req
}
