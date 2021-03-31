package mongo_models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoSequence struct {
	// MongoのユニークID
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// 招待の発行者ID
	Key string `json:"key" bson:"key"`

	// 招待の利用者ID
	Value int32 `json:"value" bson:"value"`
}
