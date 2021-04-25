package tests

import (
	"context"
	"errors"

	"github.com/UsagiBooru/accounts-server/models/const_models/account_const"
	"github.com/UsagiBooru/accounts-server/models/mongo_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InitAccountDatabase(m *mongo.Client) error {
	// Get password hash
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(PASSWORD),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errors.New("password hash create failed")
	}
	// Create accounts
	col := m.Database("accounts").Collection("users")
	users := []interface{}{
		// Admin account
		mongo_models.MongoAccountStruct{
			ID:            primitive.NewObjectID(),
			TotpCode:      "Hogehoge",
			AccountStatus: account_const.STATUS_ACTIVE,
			AccountID:     1,
			DisplayID:     "domao",
			ApiSeq:        0,
			Permission:    account_const.PERMISSION_ADMIN,
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
		},
		// Modelator account
		mongo_models.MongoAccountStruct{
			ID:            primitive.NewObjectID(),
			TotpCode:      "Hogehoge",
			AccountStatus: account_const.STATUS_ACTIVE,
			AccountID:     2,
			DisplayID:     "kafuuchino",
			ApiSeq:        0,
			Permission:    account_const.PERMISSION_MOD,
			Password:      string(hashedPassword),
			Mail:          "debug2@example.com",
			TotpEnabled:   false,
			Name:          "香風智乃",
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
		},
		// User account
		mongo_models.MongoAccountStruct{
			ID:            primitive.NewObjectID(),
			TotpCode:      "Hogehoge",
			AccountStatus: account_const.STATUS_ACTIVE,
			AccountID:     3,
			DisplayID:     "hotococoa",
			ApiSeq:        0,
			Permission:    account_const.PERMISSION_USER,
			Password:      string(hashedPassword),
			Mail:          "debug3@example.com",
			TotpEnabled:   false,
			Name:          "保登心愛",
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
				AccountID: 2,
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
		},
		// Deleted account
		mongo_models.MongoAccountStruct{
			ID:            primitive.NewObjectID(),
			TotpCode:      "Hogehoge",
			AccountStatus: account_const.STATUS_DELETED_BY_MOD,
			AccountID:     4,
			DisplayID:     "deleted",
			ApiSeq:        0,
			Permission:    account_const.PERMISSION_USER,
			Password:      string(hashedPassword),
			Mail:          "debug4@example.com",
			TotpEnabled:   false,
			Name:          "削除済みアカウント",
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
				AccountID: 3,
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
		},
	}
	if _, err := col.InsertMany(context.Background(), users); err != nil {
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
		Value: 4,
	}
	if _, err := col.InsertOne(context.Background(), seq); err != nil {
		return err
	}
	return nil
}
