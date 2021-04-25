package mongo_models

import (
	"context"
	"errors"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AccountID int32

// MongoAccountStructNotify - 通知クライアントを設定済みか
type MongoAccountStructNotify struct {
	HasLineNotify bool `bson:"hasLineNotify,omitempty"`

	HasWebNotify bool `bson:"hasWebNotify,omitempty"`
}

// MongoAccountStructInvite - 招待情報
type MongoAccountStructInvite struct {

	// 招待通し番号
	InviteID int32 `bson:"inviteID,omitempty" validate:"omitempty,gte=0"`

	// 招待コード(shortuuid)
	Code string `bson:"code,omitempty" validate:"omitempty,alphanum,min=4,max=12"`

	// 招待した人数の累計(誰を招待したかは表示されない)
	InvitedCount int32 `bson:"invitedCount,omitempty" validate:"omitempty,gte=0"`
}

type MongoAccountStructAccess struct {

	// 招待できるか
	CanInvite bool `bson:"canInvite,omitempty"`

	// いいねできるか
	CanLike bool `bson:"canLike,omitempty"`

	// コメントできるか
	CanComment bool `bson:"canComment,omitempty"`

	// 投稿できるか
	CanCreatePost bool `bson:"canCreatePost,omitempty"`

	// 投稿を編集できるか
	CanEditPost bool `bson:"canEditPost,omitempty"`

	// 投稿を承認できるか
	CanApprovePost bool `bson:"canApprovePost,omitempty"`
}

type MongoAccountStructIpfs struct {

	// 使用する任意のゲートウェイアドレス
	GatewayUrl string `bson:"gatewayUrl,omitempty" validate:"omitempty,url,max=100"`

	// 使用する任意のノードアドレス
	NodeUrl string `bson:"nodeUrl,omitempty" validate:"omitempty,url,max=100"`

	// IPFSゲートウェイを使用するか否か
	GatewayEnabled bool `bson:"gatewayEnabled,omitempty"`

	// IPFSノードを使用するか否か
	NodeEnabled bool `bson:"nodeEnabled,omitempty"`

	// マイリストを自動Pinningするか
	PinEnabled bool `bson:"pinEnabled,omitempty"`
}

type LightMongoAccountStruct struct {

	// アカウントID
	AccountID AccountID `bson:"accountID,omitempty"`

	// アカウント名
	Name string `bson:"name,omitempty"`
}

// MongoAccountStruct - アカウントの取得/編集に使うリクエスト
type MongoAccountStruct struct {
	// MongoのユニークID
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// アカウント状態 0:通常 1:ユーザー削除 2:管理者削除
	AccountStatus int32 `bson:"accountStatus,omitempty" validate:"omitempty,gte=0,lte=2"`

	// ユーザーID
	AccountID AccountID `json:"accountID,omitempty" bson:"accountID,omitempty" validate:"omitempty,gte=0"`

	// (Twitterのような)表示IDを指定します。ここで指定したIDがログインに使用されます。英数字のみ入力できます。
	DisplayID string `bson:"displayID,omitempty" validate:"omitempty,alphanum,min=3,max=20"`

	// APIキー
	ApiKey string `bson:"apiKey,omitempty" validate:"omitempty,min=0,max=500"`

	// 長期間有効トークン検証用シーケンス
	ApiSeq int32 `bson:"apiSeq,omitempty" validate:"omitempty,gte=0"`

	// 権限レベル 0:普通 5:Modelator 9:SysOp
	Permission int32 `bson:"permission,omitempty" validate:"omitempty,gte=0,lte=9"`

	// 新しいパスワードを入力します
	Password string `bson:"password,omitempty" validate:"omitempty,alphanum,min=6,max=100"`

	// ユーザーのメールアドレス(連絡用)
	Mail string `bson:"mail,omitempty" validate:"omitempty,email,max=80"`

	// TOTP認証用パスワード
	TotpCode string `bson:"totpCode,omitempty"`

	// TOTPが有効かが入ります
	TotpEnabled bool `bson:"totpEnabled,omitempty"`

	// 他のユーザーに表示されるユーザー名/投稿者名
	Name string `bson:"name,omitempty" validate:"omitempty,alphanumunicode,min=1,max=20"`

	// 他のユーザーに表示されるユーザー説明文/投稿者説明
	Description string `bson:"description,omitempty" validate:"omitempty,alphanumunicode,min=0,max=1000"`

	// ユーザーの推しキャラ(タグID)を選択します
	Favorite int32 `bson:"favorite,omitempty" validate:"omitempty,gte=0"`

	Access MongoAccountStructAccess `bson:"access,omitempty"`

	Inviter LightMongoAccountStruct `bson:"inviter,omitempty"`

	Invite MongoAccountStructInvite `bson:"invite,omitempty"`

	Notify MongoAccountStructNotify `bson:"notify,omitempty"`

	Ipfs MongoAccountStructIpfs `bson:"ipfs,omitempty"`
}

func (f *MongoAccountStruct) UpdateDisplayID(col *mongo.Collection, displayID string) (err error) {
	if displayID == "" || f.DisplayID == displayID {
		return nil
	}
	// Deny if displayID conflicted
	if err := col.FindOne(context.Background(), bson.M{"displayID": displayID}); err == nil {
		return errors.New("already used")
	}
	f.DisplayID = displayID
	return nil
}

func (f *MongoAccountStruct) UpdateName(col *mongo.Collection, name string) (err error) {
	if name == "" || f.Name == name {
		return nil
	}
	// Deny if displayID conflicted
	if err := col.FindOne(context.Background(), bson.M{"name": name}); err == nil {
		return errors.New("already used")
	}
	f.Name = name
	return nil
}

func (f *MongoAccountStruct) UpdateDescription(description string) (err error) {
	if description == "" || f.Description == description {
		return nil
	}
	f.Description = description
	return nil
}

func (f *MongoAccountStruct) UpdateMail(mail string) (err error) {
	if mail == "" {
		return nil
	}
	f.Mail = mail
	return nil
}

func (f *MongoAccountStruct) UpdateFavorite(favorite int32) (err error) {
	if favorite == 0 {
		return nil
	}
	f.Favorite = favorite
	return nil
}

func (f *MongoAccountStruct) UpdateAccess(access gen.AccountStructAccess) (err error) {
	if (access == gen.AccountStructAccess{}) {
		return nil
	}
	f.Access = MongoAccountStructAccess(access)
	return nil
}

func (f *MongoAccountStruct) UpdateIpfs(ipfs gen.AccountStructIpfs) (err error) {
	if (ipfs == gen.AccountStructIpfs{}) {
		return nil
	}
	f.Ipfs = MongoAccountStructIpfs(ipfs)
	return nil
}

func (f *MongoAccountStruct) UpdateApiSeq(update int32) (err error) {
	if update == 0 {
		return nil
	}
	f.ApiSeq += 1
	return nil
}

func (f *MongoAccountStruct) UpdatePermission(permission int32) (err error) {
	if f.Permission == permission {
		return nil
	}
	f.Permission = permission
	return nil
}

func (f *MongoAccountStruct) UpdatePassword(old_password string, new_password string) (err error) {
	if new_password == "" {
		return nil
	}
	// Validate old password hash
	if err := bcrypt.CompareHashAndPassword([]byte(f.Password), []byte(old_password)); err != nil {
		return errors.New("specified old password is incorrect")
	}
	// Get new password hash
	hashedNewPassword, err := bcrypt.GenerateFromPassword(
		[]byte(new_password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("internal password generation failed")
	}
	f.Password = string(hashedNewPassword)
	return nil
}

func (f *MongoAccountStruct) ValidatePassword(password string) (err error) {
	// Validate old password hash
	if err := bcrypt.CompareHashAndPassword([]byte(f.Password), []byte(password)); err != nil {
		return errors.New("password mismatched")
	}
	return nil
}

func (f *MongoAccountStruct) ToOpenApi(md *mongo.Client) (ac *gen.AccountStruct) {
	col := md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": f.Inviter.AccountID}
	var inviter LightMongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&inviter); err != nil {
		server.Debug(err.Error())
		return nil
	}
	inviterResp := gen.LightAccountStruct{
		AccountID: int32(inviter.AccountID),
		Name:      inviter.Name,
	}
	resp := gen.AccountStruct{
		AccountID:   int32(f.AccountID),
		DisplayID:   f.DisplayID,
		Permission:  f.Permission,
		ApiSeq:      f.ApiSeq,
		Favorite:    f.Favorite,
		Mail:        f.Mail,
		Name:        f.Name,
		Description: f.Description,
		Access:      gen.AccountStructAccess(f.Access),
		Inviter:     inviterResp,
		Invite:      gen.AccountStructInvite(f.Invite),
		Ipfs:        gen.AccountStructIpfs(f.Ipfs),
	}
	return &resp
}
