package tests

import (
	"net/http"
	"strconv"

	"github.com/UsagiBooru/accounts-server/utils/request"
)

func SetAdminUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "1")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(request.PermissionAdmin))
	return req
}

func SetModUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "2")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(request.PermissionModerator))
	return req
}

func SetNormalUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "3")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(request.PermissionUser))
	return req
}