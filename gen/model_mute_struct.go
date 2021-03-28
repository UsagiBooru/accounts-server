/*
 * UsagiBooru Accounts API
 *
 * アカウント関連API
 *
 * API version: 2.0
 * Contact: dsgamer777@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type MuteStruct struct {

	// ミュートID
	MuteID int32 `json:"muteID,omitempty"`

	// ミュート種別
	TargetType string `json:"targetType,omitempty"`

	// 対象のタグ/絵師ID
	TargetID int32 `json:"targetID,omitempty"`
}