package mongomodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoInvite struct {
	// MongoのユニークID
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// 招待の発行者ID
	Inviter AccountID `json:"inviter" bson:"inviter" validate:"gte=0"`

	// 招待の利用者ID
	Invitee AccountID `json:"invitee" bson:"invitee" validate:"gte=0"`

	// 招待コード
	Code string `json:"code" bson:"code" validate:"alphanum,min=4,max=12"`
}
