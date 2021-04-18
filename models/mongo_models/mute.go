package mongo_models

import (
	"github.com/UsagiBooru/accounts-server/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoMuteStruct struct {
	// MongoのユニークID
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// ユーザーID
	AccountID AccountID `json:"accountID,omitempty" bson:"accountID,omitempty" validate:"gte=0"`

	// ミュートID
	MuteID int32 `json:"muteID,omitempty" validate:"gte=0"`

	// ミュート種別
	TargetType string `json:"targetType,omitempty" validate:"gte=0,lte=9"`

	// 対象のタグ/絵師ID
	TargetID int32 `json:"targetID,omitempty" validate:"gte=0"`
}

func (f *MongoMuteStruct) ToOpenApi() *gen.MuteStruct {
	resp := gen.MuteStruct{
		MuteID:     f.MuteID,
		AccountID:  int32(f.AccountID),
		TargetType: f.TargetType,
		TargetID:   f.TargetID,
	}
	return &resp
}
