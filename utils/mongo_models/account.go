package mongo_models

import (
	"github.com/UsagiBooru/accounts-server/gen"
)

type MongoAccount struct {
	// TOTP認証用パスワード
	TotpKey string `json:"totpEnabled,omitempty"`

	// パスワードのSALT(ユーザー別にUUIDを発行)
	PasswordSalt string `json:"passwordSalt,omitempty"`

	gen.AccountStruct
}
