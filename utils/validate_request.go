package utils

import (
	"encoding/json"

	"github.com/UsagiBooru/accounts-server/gen"
)

func ValidateRequiredFields(req interface{}, fields []string) gen.ImplResponse {
	reqJson, err := json.Marshal(req)
	if err != nil {
		Error(err.Error())
		return gen.Response(400, gen.GeneralMessageResponse{Message: "unknown request format error"})
	}
	var mapData map[string]interface{}
	if err := json.Unmarshal([]byte(reqJson), &mapData); err != nil {
		Error(err.Error())
		return gen.Response(400, gen.GeneralMessageResponse{Message: "unknown request format error"})
	}
	for _, field := range fields {
		if _, ok := mapData[field]; !ok {
			Error(err.Error())
			return gen.Response(400, gen.GeneralMessageResponse{Message: "request parameter " + field + " was not satisfied"})
		}
	}
	return gen.Response(200, nil)
}
