/*
 * UsagiBooru Accounts API
 *
 * Accounts related api (required)
 *
 * API version: 2.0
 * Contact: dsgamer777@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package gen

// PostLoginWithFormRequest - ログインする際に利用される要求構造体
type PostLoginWithFormRequest struct {

	// ログインID
	Id string `json:"id"`

	// ログインパスワード
	Password string `json:"password"`

	// ログインTOTPトークン
	TotpCode string `json:"totpCode,omitempty"`
}
