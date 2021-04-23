package tests

import (
	"context"
	"errors"

	"github.com/UsagiBooru/accounts-server/models/mongo_models"
	"github.com/UsagiBooru/accounts-server/utils/server"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func ReGenerateTestDatabase() error {
	conf := server.GetConfig()
	m := server.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass)
	// Drop database
	drops := []string{"users", "invites", "sequence"}
	for _, d := range drops {
		col := m.Database("accounts").Collection(d)
		err := col.Drop(context.Background())
		if err != nil {
			return err
		}
	}
	// Get password hash
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(PASSWORD),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("password hash create failed")
	}
	// Create account
	col := m.Database("accounts").Collection("users")
	user := mongo_models.MongoAccountStruct{
		ID:            primitive.NewObjectID(),
		TotpCode:      "Hogehoge",
		AccountStatus: 0,
		AccountID:     1,
		DisplayID:     "domao",
		ApiSeq:        0,
		Permission:    0,
		Password:      string(hashedPassword),
		Mail:          "debug@example.com",
		TotpEnabled:   false,
		Name:          "ドマオー",
		Description:   "",
		Favorite:      0,
		Access: mongo_models.MongoAccountStructAccess{
			CanInvite:      true,
			CanLike:        true,
			CanComment:     true,
			CanCreatePost:  true,
			CanEditPost:    true,
			CanApprovePost: true,
		},
		Inviter: mongo_models.LightMongoAccountStruct{
			AccountID: 1,
		},
		Invite: mongo_models.MongoAccountStructInvite{
			Code:         "dev",
			InvitedCount: -1,
		},
		Notify: mongo_models.MongoAccountStructNotify{
			HasLineNotify: false,
			HasWebNotify:  false,
		},
		Ipfs: mongo_models.MongoAccountStructIpfs{
			GatewayUrl:     "https://cloudflare-ipfs.com",
			NodeUrl:        "",
			GatewayEnabled: false,
			NodeEnabled:    false,
			PinEnabled:     false,
		},
	}
	if _, err := col.InsertOne(context.Background(), user); err != nil {
		return err
	}
	// Create invite
	col = m.Database("accounts").Collection("invites")
	invite := mongo_models.MongoInvite{
		ID:      primitive.NewObjectID(),
		Code:    "devcode1",
		Inviter: 1,
		Invitee: 0,
	}
	if _, err := col.InsertOne(context.Background(), invite); err != nil {
		return err
	}
	// Create sequence
	col = m.Database("accounts").Collection("sequence")
	seq := mongo_models.MongoSequence{
		ID:    primitive.NewObjectID(),
		Key:   "accountID",
		Value: 1,
	}
	if _, err := col.InsertOne(context.Background(), seq); err != nil {
		return err
	}
	return nil
}
