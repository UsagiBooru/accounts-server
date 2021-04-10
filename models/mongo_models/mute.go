package mongo_models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoMuteStruct struct {
	// MongoのユニークID
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// ミュートID
	MuteID int32 `json:"muteID,omitempty"`

	// ミュート種別
	TargetType string `json:"targetType,omitempty"`

	// 対象のタグ/絵師ID
	TargetID int32 `json:"targetID,omitempty"`
}
