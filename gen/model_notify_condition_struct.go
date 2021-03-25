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

type NotifyConditionStruct struct {

	// NotifyConditionID
	NotifyConditionID int32 `json:"notifyConditionID,omitempty"`

	// 条件種別
	TargetType string `json:"targetType,omitempty"`

	// 条件ID 全通知なら0/タグID/絵師ID
	TargetID int32 `json:"targetID,omitempty"`

	// 通知方法
	TargetMethod string `json:"targetMethod,omitempty"`

	// 対象の通知クライアント(ターゲットが全てなら-1)
	TargetClient int32 `json:"targetClient,omitempty"`
}
