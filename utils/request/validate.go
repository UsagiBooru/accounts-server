package request

import (
	"encoding/json"
	"errors"

	"github.com/UsagiBooru/accounts-server/models/const_models/account_const"
)

// ValidateRequiredFields validates required fields are not empty
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

// ValidatePermission validates issuer account can change specified account
func ValidatePermission(issuerPermission int32, issuerID int32, targetID int32) error {
	notMod := issuerPermission < account_const.PERMISSION_MOD
	if targetID != issuerID && notMod {
		return errors.New("not enough permissions")
	}
	return nil
}
