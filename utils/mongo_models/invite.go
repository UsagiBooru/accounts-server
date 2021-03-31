package mongo_models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoInvite struct {
	// MongoのユニークID
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// 招待の発行者ID
	Inviter int32 `json:"inviter" bson:"inviter"`

	// 招待の利用者ID
	Invitee int32 `json:"invitee" bson:"invitee"`

	// 招待コード
	Code string `json:"code" bson:"code"`
}