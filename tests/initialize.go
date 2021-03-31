package tests

import (
	"context"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils"
	"github.com/UsagiBooru/accounts-server/utils/mongo_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReGenerateTestDatabase() error {
	conf := utils.GetConfig()
	m := utils.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass)
	// Drop database
	drops := []string{"users", "invites", "sequence"}
	for _, d := range drops {
		col := m.Database("accounts").Collection(d)
		err := col.Drop(context.Background())
		if err != nil {
			return err
		}
	}
	// Create account
	col := m.Database("accounts").Collection("users")
	user := mongo_models.MongoAccount{
		ID: primitive.NewObjectID(),
		AccountStruct: gen.AccountStruct{
			AccountID:   1,
			DisplayID:   "domao",
			ApiSeq:      0,
			Permission:  0,
			Password:    "DUMMY_PASSWORD",
			Mail:        "debug@example.com",
			TotpEnabled: false,
			Name:        "ドマオー",
			Description: "",
			Favorite:    0,
			Access: gen.AccountStructAccess{
				CanInvite:      true,
				CanLike:        true,
				CanComment:     true,
				CanCreatePost:  true,
				CanEditPost:    true,
				CanApprovePost: true,
			},
			Inviter: gen.LightAccountStruct{
				AccountID: 1,
			},
			Invite: gen.AccountStructInvite{
				Code:         "dev",
				InvitedCount: -1,
			},
			Notify: gen.AccountStructNotify{
				HasLineNotify: false,
				HasWebNotify:  false,
			},
			Ipfs: gen.AccountStructIpfs{
				GatewayUrl:     "https://cloudflare-ipfs.com",
				NodeUrl:        "",
				GatewayEnabled: false,
				NodeEnabled:    false,
				PinEnabled:     false,
			},
		},
	}
	if _, err := col.InsertOne(context.Background(), utils.ConvertStructToBson(user)); err != nil {
		return err
	}
	// Create invite
	col = m.Database("accounts").Collection("invites")
	invite := mongo_models.MongoInvite{
		ID:      primitive.NewObjectID(),
		Code:    "dev",
		Inviter: 1,
		Invitee: 0,
	}
	if _, err := col.InsertOne(context.Background(), utils.ConvertStructToBson(invite)); err != nil {
		return err
	}
	// Create sequence
	col = m.Database("accounts").Collection("sequence")
	seq := mongo_models.MongoSequence{
		ID:    primitive.NewObjectID(),
		Key:   "accountID",
		Value: 1,
	}
	if _, err := col.InsertOne(context.Background(), utils.ConvertStructToBson(seq)); err != nil {
		return err
	}
	return nil
}
