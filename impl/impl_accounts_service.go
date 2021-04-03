package impl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils"
	"github.com/UsagiBooru/accounts-server/utils/mongo_models"
	"github.com/elastic/go-elasticsearch/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AccountsApiImplService struct {
	gen.AccountsApiService
	es *elasticsearch.Client
	md *mongo.Client
}

func NewAccountsApiImplService() gen.AccountsApiServicer {
	conf := utils.GetConfig()
	return &AccountsApiImplService{
		AccountsApiService: gen.AccountsApiService{},
		es:                 utils.NewElasticSearchClient(conf.ElasticHost, conf.ElasticUser, conf.ElasticPass),
		md:                 utils.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass),
	}
}

// CreateAccount - Create account
func (s *AccountsApiImplService) CreateAccount(ctx context.Context, accountStruct gen.AccountStruct) (gen.ImplResponse, error) {
	// Timeout of this method is 3 seconds
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	v := ctx.Value("user-id")
	token, ok := v.(string)
	utils.Debug("User id is" + string(token))
	if !ok {
		return gen.Response(500, gen.AccountStruct{}), nil
	}

	// Validate request fields
	if resp := utils.ValidateRequiredFields(
		accountStruct,
		[]string{"name", "displayID", "password", "mail"},
	); resp.Code != http.StatusOK {
		return resp, nil
	}
	if resp := utils.ValidateRequiredFields(
		accountStruct.Invite,
		[]string{"code"},
	); resp.Code != http.StatusOK {
		return resp, nil
	}

	// Use transaction to prevent duplicate request
	err := s.md.UseSession(ctx, func(sc mongo.SessionContext) error {
		// Start transaction
		err := sc.StartTransaction()
		if err != nil {
			return err
		}
		// Find invite code
		col := s.md.Database("accounts").Collection("invites")
		filter := bson.M{
			"code":    accountStruct.Invite.Code,
			"invitee": 0,
		}
		var invite mongo_models.MongoInvite
		if err := col.FindOne(context.Background(), filter).Decode(&invite); err != nil {
			utils.Debug(err.Error())
			return errors.New("invite code was not found")
		}
		// Find inviter account
		col = s.md.Database("accounts").Collection("users")
		filter = bson.M{"accountID": invite.Inviter}
		var inviter mongo_models.MongoAccount
		if err := col.FindOne(context.Background(), filter).Decode(&inviter); err != nil {
			utils.Debug(err.Error())
			return errors.New("inviter account was not found")
		}
		// Get latest-1 accountID
		col = s.md.Database("accounts").Collection("sequence")
		filter = bson.M{"key": "accountID"}
		var seq mongo_models.MongoSequence
		if err := col.FindOne(context.Background(), filter).Decode(&seq); err != nil {
			utils.Debug(err.Error())
			return errors.New("accountID sequence was not found")
		}
		// Update invite invitee
		col = s.md.Database("accounts").Collection("invites")
		filter = bson.M{
			"_id": invite.ID,
		}
		set := bson.M{"$set": bson.M{"invitee": seq.Value + 1}}
		if _, err = col.UpdateOne(ctx, filter, set); err != nil {
			utils.Debug(err.Error())
			return errors.New("update invite invitee failed")
		}
		// Get password hash
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(accountStruct.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			utils.Debug(err.Error())
			return errors.New("password hash create failed")
		}
		// Create new invite for new account
		newInviteCodeForNew := utils.GetShortUUID(8)
		newInviteForNew := mongo_models.MongoInvite{
			ID:      primitive.NewObjectID(),
			Code:    newInviteCodeForNew,
			Inviter: seq.Value + 1,
			Invitee: 0,
		}
		if _, err = col.InsertOne(ctx, utils.ConvertStructToBson(newInviteForNew)); err != nil {
			utils.Debug(err.Error())
			return errors.New("insert new invite for new account failed")
		}
		// Create new invite for old account
		newInviteCodeForOld := utils.GetShortUUID(8)
		newInviteForOld := mongo_models.MongoInvite{
			ID:      primitive.NewObjectID(),
			Code:    newInviteCodeForOld,
			Inviter: invite.Inviter,
		}
		invite_str := fmt.Sprint(newInviteForOld)
		utils.Debug(invite_str)
		if _, err = col.InsertOne(ctx, utils.ConvertStructToBson(newInviteForOld)); err != nil {
			utils.Debug(err.Error())
			return errors.New("insert new invite for old account failed")
		}
		// Create new mongo user model
		col = s.md.Database("accounts").Collection("users")
		user := mongo_models.MongoAccount{
			ID: primitive.NewObjectID(),
			AccountStruct: gen.AccountStruct{
				AccountID:   seq.Value + 1,
				DisplayID:   accountStruct.DisplayID,
				ApiSeq:      0,
				Permission:  0,
				Password:    string(hashedPassword),
				Mail:        accountStruct.Mail,
				TotpEnabled: false,
				Name:        accountStruct.Name,
				Description: "",
				Favorite:    0,
				Access: gen.AccountStructAccess{
					CanInvite:      true,
					CanLike:        true,
					CanComment:     true,
					CanCreatePost:  true,
					CanEditPost:    false,
					CanApprovePost: false,
				},
				Inviter: gen.LightAccountStruct{
					AccountID: invite.Inviter,
				},
				Invite: gen.AccountStructInvite{
					InvitedCount: -1,
					Code:         newInviteCodeForNew,
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
		// Insert new user
		if _, err = col.InsertOne(ctx, utils.ConvertStructToBson(user)); err != nil {
			return errors.New("insert newuser failed")
		}
		// Update sequence
		col = s.md.Database("accounts").Collection("sequence")
		filter = bson.M{"key": "accountID"}
		set = bson.M{"$set": bson.M{"value": seq.Value + 1}}
		if _, err = col.UpdateOne(ctx, filter, set); err != nil {
			return errors.New("update accountID sequence failed")
		}
		// Update inviter's invite count
		col = s.md.Database("accounts").Collection("users")
		filter = bson.M{"accountID": invite.Inviter}
		if inviter.Invite.InvitedCount == -1 {
			inviter.Invite.InvitedCount = 0
		}
		set = bson.M{"$set": bson.M{
			"invite.invitedCount": inviter.Invite.InvitedCount + 1,
			"invite.code":         newInviteCodeForOld,
		}}
		resp, err := col.UpdateOne(ctx, filter, set)
		if err != nil {
			return errors.New("update inviter's invite count failed")
		}
		resp_str := fmt.Sprint(resp)
		utils.Debug(resp_str)
		// Commit insert user / update sequence / update invite code
		return sc.CommitTransaction(sc)

	})
	if err != nil {
		return gen.Response(500, gen.GeneralMessageResponse{Message: err.Error()}), nil
	}
	return gen.Response(200, gen.AccountStruct{}), nil
}

// GetAccount - Get account info
func (s *AccountsApiImplService) GetAccount(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": accountID}
	var account gen.AccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		utils.Debug(err.Error())
		return gen.Response(404, gen.GeneralMessageResponse{Message: "Specified account does not exist."}), nil
	}
	// Find inviter account
	col = s.md.Database("accounts").Collection("users")
	filter = bson.M{"accountID": account.Inviter.AccountID}
	var inviter gen.AccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&inviter); err != nil {
		utils.Debug(err.Error())
		return gen.Response(500, gen.GeneralMessageResponse{Message: "Internal server error."}), nil
	}
	// Create response
	accountResp := gen.AccountStruct{
		AccountID:   account.AccountID,
		DisplayID:   account.DisplayID,
		Permission:  account.Permission,
		Name:        account.Name,
		Description: account.Description,
		Favorite:    account.Favorite,
		Access: gen.AccountStructAccess{
			CanInvite:      account.Access.CanInvite,
			CanLike:        account.Access.CanLike,
			CanComment:     account.Access.CanComment,
			CanCreatePost:  account.Access.CanCreatePost,
			CanEditPost:    account.Access.CanEditPost,
			CanApprovePost: account.Access.CanApprovePost,
		},
		Inviter: gen.LightAccountStruct{
			AccountID: account.Inviter.AccountID,
			Name:      inviter.Name,
		},
	}
	return gen.Response(200, accountResp), nil
}
