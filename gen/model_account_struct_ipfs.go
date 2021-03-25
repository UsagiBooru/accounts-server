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

// AccountStructIpfs - IPFS設定
type AccountStructIpfs struct {

	// 使用する任意のゲートウェイアドレス
	GatewayUrl string `json:"gatewayUrl,omitempty"`

	// 使用する任意のノードアドレス
	NodeUrl string `json:"nodeUrl,omitempty"`

	// IPFSゲートウェイを使用するか否か
	GatewayEnabled bool `json:"gatewayEnabled,omitempty"`

	// IPFSノードを使用するか否か
	NodeEnabled bool `json:"nodeEnabled,omitempty"`

	// マイリストを自動Pinningするか
	PinEnabled bool `json:"pinEnabled,omitempty"`
}
