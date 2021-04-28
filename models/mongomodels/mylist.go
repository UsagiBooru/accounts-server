/*
 * UsagiBooru Accounts API
 *
 * アカウント関連API
 *
 * API version: 2.0
 * Contact: dsgamer777@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package mongomodels

import (
	"time"
)

// MongoLightArtStruct - 簡易イラスト情報(読み取り専用)
type MongoLightArtStruct struct {
	// イラストID
	ArtID int32 `json:"artID,omitempty" validate:"gte=0"`
}

// MongoMylistStruct - マイリスト情報
type MongoMylistStruct struct {

	// マイリストID
	MylistID int32 `json:"mylistID,omitempty"`

	// マイリスト名
	Name string `json:"name,omitempty" validate:"omitempty,alphanumunicode,min=1,max=50"`

	// マイリスト説明文
	Description string `json:"description,omitempty" validate:"omitempty,alphanumunicode,min=1,max=50"`

	// マイリスト作成日時
	CreatedDate time.Time `json:"createdDate,omitempty"`

	// マイリスト更新日時
	UpdatedDate time.Time `json:"updatedDate,omitempty"`

	// 公開/非公開
	Private bool `json:"private,omitempty"`

	// イラストID一覧
	Arts []MongoLightArtStruct `json:"arts,omitempty"`

	// マイリスト所有者の簡易アカウント情報
	Owner LightMongoAccountStruct `json:"owner,omitempty"`
}
