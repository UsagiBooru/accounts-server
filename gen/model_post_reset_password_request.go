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

// PostResetPasswordRequest - パスワードをリセットする際に使う要求構造体
type PostResetPasswordRequest struct {

	// メールアドレス(のみ)
	Mail string `json:"mail,omitempty"`
}
