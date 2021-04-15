package mongo_models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type MongoAccountHelper struct {
	col *mongo.Collection
}

func NewMongoAccountHelper(md *mongo.Client) MongoAccountHelper {
	return MongoAccountHelper{
		md.Database("accounts").Collection("users"),
	}
}

func (h *MongoAccountHelper) CreateAccount(
	accountID AccountID,
	displayID string,
	password string, mail string,
	name string,
	inviterID AccountID, inviteCode string,
) (*MongoAccountStruct, error) {
	// Get password hash
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.New("create password hash failed")
	}
	account := MongoAccountStruct{
		ID:            primitive.NewObjectID(),
		AccountStatus: 0,
		AccountID:     accountID,
		DisplayID:     displayID,
		ApiKey:        "",
		ApiSeq:        0,
		Permission:    0,
		Password:      string(hashedPassword),
		Mail:          mail,
		TotpCode:      "",
		TotpEnabled:   false,
		Name:          name,
		Description:   "",
		Favorite:      0,
		Access: MongoAccountStructAccess{
			CanInvite:      true,
			CanLike:        true,
			CanComment:     true,
			CanCreatePost:  true,
			CanEditPost:    false,
			CanApprovePost: false,
		},
		Inviter: LightMongoAccountStruct{
			AccountID: inviterID,
		},
		Invite: MongoAccountStructInvite{
			InvitedCount: 0,
			Code:         inviteCode,
		},
		Notify: MongoAccountStructNotify{
			HasLineNotify: false,
			HasWebNotify:  false,
		},
		Ipfs: MongoAccountStructIpfs{
			GatewayUrl:     "https://cloudflare-ipfs.com",
			NodeUrl:        "",
			GatewayEnabled: false,
			NodeEnabled:    false,
			PinEnabled:     false,
		},
	}
	// Insert new user
	if _, err = h.col.InsertOne(context.Background(), account); err != nil {
		return nil, errors.New("insert account failed")
	}
	return &account, nil
}

func (h *MongoAccountHelper) FindAccount(accountID AccountID) (*MongoAccountStruct, error) {
	filter := bson.M{"accountID": int32(accountID)}
	var account MongoAccountStruct
	if err := h.col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		return nil, errors.New("account was not found")
	}
	return &account, nil
}

func (h *MongoAccountHelper) DeleteAccount(accountID AccountID, deleteMethod int32) error {
	account, err := h.FindAccount(accountID)
	if err != nil {
		return err
	}
	account.AccountStatus = deleteMethod
	filter := bson.M{"accountID": int32(accountID)}
	set := bson.M{"$set": account}
	if _, err = h.col.UpdateOne(context.Background(), filter, set); err != nil {
		return errors.New("delete account failed")
	}
	return nil
}

func (h *MongoAccountHelper) UpdateInvite(accountID AccountID, code string, invitedCount int32) error {
	filter := bson.M{"accountID": int32(accountID)}
	set := bson.M{"$set": bson.M{
		"invite.invitedCount": invitedCount,
		"invite.code":         code,
	}}
	if _, err := h.col.UpdateOne(context.Background(), filter, set); err != nil {
		return errors.New("update inviter's invite count failed")
	}
	return nil
}