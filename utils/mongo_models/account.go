package mongo_models

import (
	"github.com/UsagiBooru/accounts-server/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoAccount struct {
	// MongoのユニークID
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// TOTP認証用パスワード
	TotpKey string `json:"totpEnabled,omitempty"`

	gen.AccountStruct
}
