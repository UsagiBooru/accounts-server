package utils

import (
	"github.com/UsagiBooru/accounts-server/gen"
)

func VerifyServerError(err error) gen.ImplResponse {
	if err != nil {
		Error(err.Error())
		return gen.Response(500, gen.GeneralMessageResponse{Message: "Internal server error."})
	}
	return gen.Response(200, nil)
}
