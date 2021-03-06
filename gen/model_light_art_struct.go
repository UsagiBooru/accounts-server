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

// LightArtStruct - イラスト情報の簡易構造体(読み取り専用)(検索結果等で利用)
type LightArtStruct struct {

	// イラストID
	ArtID int32 `json:"artID,omitempty"`

	// 絵師情報(複数可)
	Artists []LightArtistStruct `json:"artists,omitempty"`

	// 説明文 NOTE: 通常出典記載の説明文と同じ物が入る
	Caption string `json:"caption,omitempty"`

	// 登録日(%Y-%m-%d %H:%M:%S)
	Datetime time.Time `json:"datetime,omitempty"`

	File LightArtStructFile `json:"file,omitempty"`

	// 累計いいね数
	Likes int64 `json:"likes,omitempty"`

	// リクエストしたユーザーがマイリストしているか
	Mylisted bool `json:"mylisted,omitempty"`

	// マイリスト済みのユーザー数
	Mylists int64 `json:"mylists,omitempty"`

	// アダルトコンテンツか否か
	Nsfw bool `json:"nsfw,omitempty"`

	// 出典のサービス名
	OriginService string `json:"originService,omitempty"`

	// 出典URL
	OriginUrl string `json:"originUrl,omitempty"`

	// グループになっている場合のページ番号
	Page int32 `json:"page,omitempty"`

	// 元画像との類似度(画像検索のみ)
	Similarity float32 `json:"similarity,omitempty"`

	// イラスト(作品)名
	Title string `json:"title,omitempty"`

	Uploader LightAccountStruct `json:"uploader,omitempty"`

	// 累計閲覧数
	Views int64 `json:"views,omitempty"`
}
