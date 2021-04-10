package request

import (
	"encoding/json"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils/response"
)

func ValidateRequiredFields(req interface{}, fields []string) gen.ImplResponse {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return response.NewRequestErrorWithMessage("unknown request format error")
	}
	var mapData map[string]interface{}
	if err := json.Unmarshal([]byte(reqJson), &mapData); err != nil {
		return response.NewRequestErrorWithMessage("unknown request format error")
	}
	for _, field := range fields {
		if _, ok := mapData[field]; !ok {
			return response.NewRequestErrorWithMessage("request parameter " + field + " was not satisfied")
		}
	}
	return gen.Response(200, nil)
}

func ValidatePermission(issuerPermission int32, issuerID int32, targetID int32) gen.ImplResponse {
	notMod := issuerPermission < PermissionModerator
	if targetID != issuerID && notMod {
		return response.NewPermissionError()
	}
	return gen.Response(200, nil)
}
