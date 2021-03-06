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

import (
	"time"
)

// MylistStruct - マイリスト情報の構造体
type MylistStruct struct {

	// イラストID一覧
	Arts []LightArtStruct `json:"arts,omitempty"`

	// マイリスト作成日時
	CreatedDate time.Time `json:"createdDate,omitempty"`

	// マイリスト説明文
	Description string `json:"description,omitempty"`

	// マイリストID
	MylistID int32 `json:"mylistID,omitempty"`

	// マイリスト名
	Name string `json:"name,omitempty"`

	Owner LightAccountStruct `json:"owner,omitempty"`

	// 公開/非公開
	Private bool `json:"private,omitempty"`

	// マイリスト更新日時
	UpdatedDate time.Time `json:"updatedDate,omitempty"`
}
