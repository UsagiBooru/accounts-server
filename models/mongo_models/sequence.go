package mongo_models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSequence struct {
	// MongoのユニークID
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	// 招待の発行者ID
	Key string `json:"key" bson:"key"`

	// 招待の利用者ID
	Value int32 `json:"value" bson:"value"`
}

func GetSeq(md *mongo.Client, dbName string, seqKey string) (resp int32, err error) {
	col := md.Database(dbName).Collection("sequence")
	filter := bson.M{"key": seqKey}
	var seq MongoSequence
	if err := col.FindOne(context.Background(), filter).Decode(&seq); err != nil {
		return 0, errors.New("get " + seqKey + " sequence failed")
	}
	return int32(seq.Value), nil
}

func UpdateSeq(md *mongo.Client, dbName string, seqKey string, seqCurrent int32) (err error) {
	col := md.Database(dbName).Collection("sequence")
	filter := bson.M{"key": seqKey}
	set := bson.M{"$set": bson.M{"value": seqCurrent + 1}}
	if _, err = col.UpdateOne(context.Background(), filter, set); err != nil {
		return errors.New("update sequence failed")
	}
	return nil
}
