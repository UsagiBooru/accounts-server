package mongomodels

import (
	"context"
	"errors"

	"github.com/UsagiBooru/accounts-server/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoMuteHelper struct {
	col *mongo.Collection
}

// NewMongoMuteHelper creates a helper for handle mutes endpoints
func NewMongoMuteHelper(md *mongo.Client) MongoMuteHelper {
	return MongoMuteHelper{md.Database("accounts").Collection("mutes")}
}

func (h *MongoMuteHelper) ToMongo(mt gen.MuteStruct) *MongoMuteStruct {
	resp := MongoMuteStruct{
		MuteID:     mt.MuteID,
		AccountID:  AccountID(mt.AccountID),
		TargetType: mt.TargetType,
		TargetID:   mt.TargetID,
	}
	return &resp
}

func (h *MongoMuteHelper) CreateMute(muteID int32, targetType string, targetID int32) (*MongoMuteStruct, error) {
	newMute := MongoMuteStruct{
		ID:         primitive.NewObjectID(),
		MuteID:     muteID,
		TargetType: targetType,
		TargetID:   targetID,
	}
	if _, err := h.col.InsertOne(context.Background(), newMute); err != nil {
		return nil, errors.New("insert mute failed")
	}
	return &newMute, nil
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

func (h *MongoMuteHelper) FindDuplicatedMute(targetType string, targetID int32, accountID AccountID) error {
	filter := bson.M{
		"targetType": targetType,
		"targetID":   targetID,
		"accountID":  accountID,
	}
	var Mute MongoMuteStruct
	if err := h.col.FindOne(context.Background(), filter).Decode(&Mute); err == nil {
		return errors.New("duplicated mute was found")
	}
	return nil
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
