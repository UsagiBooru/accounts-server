package tests

import (
	"context"

	"github.com/UsagiBooru/accounts-server/models/mongomodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitMuteDatabase(m *mongo.Client) error {
	// Create mute
	col := m.Database("accounts").Collection("mutes")
	newMute := mongomodels.MongoMuteStruct{
		ID:         primitive.NewObjectID(),
		AccountID:  1,
		MuteID:     1,
		TargetType: "artist",
		TargetID:   1,
	}
	if _, err := col.InsertOne(context.Background(), newMute); err != nil {
		return err
	}
	// Create sequence
	col = m.Database("accounts").Collection("sequence")
	seq := mongomodels.MongoSequence{
		ID:    primitive.NewObjectID(),
		Key:   "muteID",
		Value: 1,
	}
	if _, err := col.InsertOne(context.Background(), seq); err != nil {
		return err
	}
	return nil
}
