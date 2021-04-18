package mongo_models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoMuteHelper struct {
	col *mongo.Collection
}

func NewMongoMuteHelper(md *mongo.Client) MongoMuteHelper {
	return MongoMuteHelper{md.Database("accounts").Collection("mutes")}
}

func (h *MongoMuteHelper) CreateMute(muteID int32, targetType string, targetID int32) (*MongoMuteStruct, error) {
	newMuteForNew := MongoMuteStruct{
		ID:         primitive.NewObjectID(),
		MuteID:     muteID,
		TargetType: targetType,
		TargetID:   targetID,
	}
	if _, err := h.col.InsertOne(context.Background(), newMuteForNew); err != nil {
		return nil, errors.New("insert mute failed")
	}
	return &newMuteForNew, nil
}

func (h *MongoMuteHelper) FindMute(muteID int32) (*MongoMuteStruct, error) {
	filter := bson.M{
		"muteID": muteID,
	}
	var Mute MongoMuteStruct
	if err := h.col.FindOne(context.Background(), filter).Decode(&Mute); err != nil {
		return nil, errors.New("mute was not found")
	}
	return &Mute, nil
}

func (h *MongoMuteHelper) FindMuteUsingFilter(filter bson.M) (*MongoMuteStruct, error) {
	var Mute MongoMuteStruct
	if err := h.col.FindOne(context.Background(), filter).Decode(&Mute); err != nil {
		return nil, errors.New("mute was not found")
	}
	return &Mute, nil
}

func (h *MongoMuteHelper) DeleteMute(muteID int32) error {
	filter := bson.M{
		"muteID": muteID,
	}
	if _, err := h.col.DeleteOne(context.Background(), filter); err != nil {
		return errors.New("delete Mute failed")
	}
	return nil
}
