package mongomodels

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSequenceHelper struct {
	md         *mongo.Client
	dbName     string
	seqName    string
	seqCurrent int32
}

func NewMongoSequenceHelper(md *mongo.Client, dbName string, seqName string) MongoSequenceHelper {
	return MongoSequenceHelper{md, dbName, seqName, 0}
}

func (m *MongoSequenceHelper) GetSeq() (resp int32, err error) {
	col := m.md.Database(m.dbName).Collection("sequence")
	filter := bson.M{"key": m.seqName}
	var seq MongoSequence
	if err := col.FindOne(context.Background(), filter).Decode(&seq); err != nil {
		return 0, errors.New("get " + m.seqName + " sequence failed")
	}
	m.seqCurrent = seq.Value
	return int32(seq.Value), nil
}

func (m *MongoSequenceHelper) UpdateSeq() (err error) {
	col := m.md.Database(m.dbName).Collection("sequence")
	filter := bson.M{"key": m.seqName}
	set := bson.M{"$set": bson.M{"value": m.seqCurrent + 1}}
	if _, err = col.UpdateOne(context.Background(), filter, set); err != nil {
		return errors.New("update sequence failed")
	}
	return nil
}
