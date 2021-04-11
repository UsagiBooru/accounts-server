package request

import (
	"encoding/json"
	"errors"
)

func ValidateRequiredFields(req interface{}, fields []string) error {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return errors.New("unknown request format error")
	}
	var mapData map[string]interface{}
	if err := json.Unmarshal([]byte(reqJson), &mapData); err != nil {
		return errors.New("unknown request format error")
	}
	for _, field := range fields {
		if _, ok := mapData[field]; !ok {
			return errors.New("request parameter " + field + " was not satisfied")
		}
	}
	return nil
}

func ValidatePermission(issuerPermission int32, issuerID int32, targetID int32) error {
	notMod := issuerPermission < PermissionModerator
	if targetID != issuerID && notMod {
		return errors.New("not enough permissions")
	}
	return nil
}
