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

// UploadHistoryStruct - 投稿履歴の応答構造体
type UploadHistoryStruct struct {

	// 情報作成者のアカウントID
	AccountID int32 `json:"accountID"`

	// 投稿に成功した場合入るID
	ArtID int32 `json:"artID"`

	// データ登録完了時刻
	Finished string `json:"finished"`

	// データ登録処理開始時刻
	Started string `json:"started"`

	// 登録処理結果 5:成功 9:内部エラー
	Status int32 `json:"status"`

	// 通し投稿履歴番号(インデックス用)
	UploadID int32 `json:"uploadID"`
}
