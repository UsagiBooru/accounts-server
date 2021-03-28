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

// AccountStruct - アカウントの取得/編集に使うリクエスト
type AccountStruct struct {

	// ユーザーID
	AccountID int32 `json:"accountID,omitempty"`

	// (Twitterのような)表示IDを指定します。ここで指定したIDがログインに使用されます。英数字のみ入力できます。
	DisplayID string `json:"displayID,omitempty"`

	// APIキー
	ApiKey string `json:"apiKey,omitempty"`

	// 長期間有効トークン検証用シーケンス
	ApiSeq int32 `json:"apiSeq,omitempty"`

	// 権限レベル 0:普通 5:Modelator 9:SysOp
	Permission int32 `json:"permission,omitempty"`

	// 新しいパスワードを入力します
	Password string `json:"password,omitempty"`

	// 現時点のパスワードを入力します。 userPasswordを変更する場合に必要となります。
	OldPassword string `json:"oldPassword,omitempty"`

	// TOTPが有効かが入ります
	TotpEnabled bool `json:"totpEnabled,omitempty"`

	// 他のユーザーに表示されるユーザー名/投稿者名
	Name string `json:"name,omitempty"`

	// 他のユーザーに表示されるユーザー説明文/投稿者説明
	Description string `json:"description,omitempty"`

	// ユーザーの推しキャラ(タグID)を選択します
	Favorite int32 `json:"favorite,omitempty"`

	Access AccountStructAccess `json:"access,omitempty"`

	Inviter LightAccountStruct `json:"inviter,omitempty"`

	Invite AccountStructInvite `json:"invite,omitempty"`

	Notify AccountStructNotify `json:"notify,omitempty"`

	Ipfs AccountStructIpfs `json:"ipfs,omitempty"`
}